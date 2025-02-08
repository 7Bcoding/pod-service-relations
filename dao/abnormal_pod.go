package dao

import (
	"fmt"
	"go/types"
	"time"

	"gorm.io/gorm"

	"pod-service-relations/model"
)

type abnormalPodDao struct {
	db *gorm.DB
}

type AbnormalPodDao interface {
	Create(AbnormalPod *model.AbnormalPod) error
	Get(PodName string) (*model.AbnormalPod, error)
	ALL(filters map[string]interface{}) ([]*model.AbnormalPod, error)
	Delete(ID uint, PodName string) error
	Update(Pod *model.AbnormalPod, fields map[string]interface{}) error
}

func NewAbnormalPodDao(db *gorm.DB) AbnormalPodDao {
	return &abnormalPodDao{db: db}
}

// Create
/*
* insert Pod
* PARAMS:
*	- Pod: model Pod
* RETURN:
* 	- error: if success, nil, else gorm error and other error
 */
func (ap abnormalPodDao) Create(AbnormalPod *model.AbnormalPod) error {
	return ap.db.Create(AbnormalPod).Error
}

// Get
/*
* select single Pod
* params:
*	- id: pk
*	- name: unique key
* RETURN:
*	- model.Pod
*	- error: if success, nil, else gorm error and other error
 */
func (ap abnormalPodDao) Get(PodName string) (*model.AbnormalPod, error) {
	if PodName == "" {
		return nil, fmt.Errorf(
			"Pod_name must set one: %s", PodName,
		)
	}
	Pod := new(model.AbnormalPod)
	db := ap.db
	if PodName != "" {
		db = db.Where("pod_name = ?", PodName)
	}
	err := db.First(Pod).Error
	if err != nil {
		return nil, err
	}
	return Pod, nil
}

// ALL
/*
* select batch Pods
* params:
*	- filters: map[string]interface{}
*		key: db column
*		value: interface{}, maybe uint32, string, slice
*		like {"id": 1, "name": ["test1", "test2"]}
* RETURN:
*	- model.Pod slice
*	- error: if success, nil, else gorm error and other error
 */
func (ap abnormalPodDao) ALL(filters map[string]interface{}) ([]*model.AbnormalPod, error) {
	var AbnormalPods []*model.AbnormalPod
	db := ap.db
	for k, v := range filters {
		switch value := v.(type) {
		case string, uint32, time.Time:
			db = db.Where(k+"= ?", value)
		case types.Slice:
			db = db.Where(k+"IN ?", value)
		default:
			return nil, fmt.Errorf("unsupport value type: %s", v)
		}
	}
	err := db.Find(&AbnormalPods).Error
	return AbnormalPods, err
}

// Delete
/*
* delete single Pod
* PARAMS:
* 	- id: pk
*	- name: unique key
* RETURN:
*	- error: if success, nil, else gorm error and other error
 */
func (ap abnormalPodDao) Delete(id uint, name string) error {
	db := ap.db
	// id or name must set, else will delete all tables
	if id > 0 || name != "" {
		if id > 0 {
			db = db.Where("id = ?", id)
		}
		if name != "" {
			db = db.Where("pod_name = ?", name)
		}
	} else {
		return fmt.Errorf("id or name must set one")
	}

	return db.Delete(model.AbnormalPod{}).Error
}

// Update
/*
* update single Pod
* PARAMS:
* 	- action: model.Pod
*	- filters: map[string]interface{}
*		key: db column
*		value: interface{}, maybe uint32, string
*		like {"bin_repo": "baidu@123"}
* RETURN:
* 	- error: if success nil, else gorm error or other error
 */
func (ap abnormalPodDao) Update(Pod *model.AbnormalPod, fields map[string]interface{}) error {
	return ap.db.Model(Pod).Updates(fields).Error
}
