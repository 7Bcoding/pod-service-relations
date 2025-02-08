package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

const (
	jsonValidatorTag = "json-validator"
	jsonTag          = "json"
	tagValMandatory  = "mandatory"
	tagValOptional   = "optional"
)

type ValidatorFunc func(val interface{}) bool
type validator struct {
	functions map[string]ValidatorFunc
	mutex     *sync.RWMutex
}

func (v *validator) get(name string) (ValidatorFunc, bool) {
	v.mutex.RLock()
	defer v.mutex.RUnlock()
	fn, ok := v.functions[name]
	return fn, ok
}

func (v *validator) register(name string, fn ValidatorFunc) error {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	_, ok := v.functions[name]
	if ok {
		return fmt.Errorf("validator with name %s already exists", name)
	}
	v.functions[name] = fn
	return nil
}

var defaultValidator validator
var initOnce sync.Once

func init() {
	initOnce.Do(func() {
		defaultValidator = validator{
			functions: make(map[string]ValidatorFunc),
			mutex:     new(sync.RWMutex),
		}
		_ = defaultValidator.register("json", jsonValidator)
		_ = defaultValidator.register("base64", base64Validator)
	})
}

func Register(name string, fn ValidatorFunc) error {
	return defaultValidator.register(name, fn)
}

func BindJSONWithContext(container interface{}, ctx *gin.Context) (map[string]interface{}, error) {
	raw, err := ctx.GetRawData()
	if err != nil && err != io.EOF {
		return nil, err
	}
	return BindJSON(container, raw)
}

func BindJSON(container interface{}, rawJSON []byte) (map[string]interface{}, error) {
	var jsonMap = make(map[string]interface{})
	if err := json.Unmarshal(rawJSON, &jsonMap); err != nil {
		return nil, err
	}
	// panics if container is non-pointer since it indicates a bug in your code
	t := reflect.TypeOf(container).Elem()
	v := reflect.ValueOf(container).Elem()
	// panics for the same reason, container must be a pointer to a struct
	n := t.NumField()
	var res = make(map[string]interface{})
	for i := 0; i < n; i++ {
		f := t.Field(i)
		fv := v.Field(i)
		jsonTag := f.Tag.Get(jsonTag)
		if jsonTag == "" {
			continue
		}
		rawValidatorTags := f.Tag.Get(jsonValidatorTag)
		if rawValidatorTags == "" {
			continue
		}
		validatorTags := strings.Split(rawValidatorTags, ",")
		val, ok := jsonMap[jsonTag]
		switch validatorTags[0] {
		case tagValMandatory:
			if !ok {
				return nil, fmt.Errorf("mandatory field %s is missing", jsonTag)
			}
		case tagValOptional:
			if !ok {
				continue
			}
		default:
			panic("invalid json-validator tag")
		}
		for j := 1; j < len(validatorTags); j++ {
			validationFn, ok := defaultValidator.get(validatorTags[j])
			if !ok {
				panic("unregistered json-validator tag: " + validatorTags[j])
			}
			if !validationFn(val) {
				return nil, fmt.Errorf("%s validation on field %s failed with value: %v", validatorTags[j], jsonTag, val)
			}
		}
		if err := jsonTypeReflector(fv, reflect.ValueOf(val)); err != nil {
			return nil, fmt.Errorf("validation on field %s failed with error: %w", jsonTag, err)
		}

		res[jsonTag] = fv.Interface()
	}
	return res, nil
}

func jsonTypeReflector(v reflect.Value, jv reflect.Value) error {
	switch jvKind := jv.Kind(); jvKind {
	case reflect.String, reflect.Bool:
		if v.Kind() != jvKind {
			return fmt.Errorf("type mismatch")
		}
		v.Set(jv)
	case reflect.Float64:
		val := jv.Float()
		switch v.Kind() {
		case reflect.Float64, reflect.Float32:
			v.SetFloat(val)
		case reflect.Int, reflect.Int32, reflect.Int64:
			v.SetInt(int64(val))
		case reflect.Uint, reflect.Uint32, reflect.Uint64:
			if val < 0 {
				return fmt.Errorf("unsigned int overflow")
			}
			v.SetUint(uint64(val))
		}
	default:
		return fmt.Errorf("unsupported type %s", v.Kind())
	}
	return nil
}

var jsonValidator = func(val interface{}) bool {
	raw, ok := val.(string)
	return !ok || json.Valid([]byte(raw))
}

var base64Validator = func(val interface{}) bool {
	raw, ok := val.(string)
	if !ok {
		return true
	}
	_, err := base64.StdEncoding.DecodeString(raw)
	return err == nil
}
