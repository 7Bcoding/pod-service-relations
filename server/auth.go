package server

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin"
)

type BasicAuth struct {
	User  string
	Token string
}

var (
	basicAuthPattern    = regexp.MustCompile("Basic (.*)")
	splitter            = []byte(":")
	errInvalidBasicAuth = fmt.Errorf("invalid basic auth format")
)

func extractBasicAuth(ctx *gin.Context) (BasicAuth, error) {
	authStr := ctx.Request.Header.Get("Authorization")
	groups := basicAuthPattern.FindStringSubmatch(authStr)
	if len(groups) != 2 {
		return BasicAuth{}, errInvalidBasicAuth
	}
	decoded, err := base64.StdEncoding.DecodeString(groups[1])
	if err != nil {
		return BasicAuth{}, err
	}
	pair := bytes.Split(decoded, splitter)
	if len(pair) != 2 {
		return BasicAuth{}, errInvalidBasicAuth
	}
	return BasicAuth{User: string(pair[0]), Token: string(pair[1])}, nil
}
