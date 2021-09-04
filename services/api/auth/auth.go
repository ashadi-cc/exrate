package auth

import (
	"fmt"
	"xrate/common/store/driver"

	"github.com/google/uuid"
)

const storeKey = "access_key"

type accessKey map[string]string

//Auth base auth methods
type Auth interface {
	//Create returns new access key by given project name
	Create(projectName string) (string, error)
	//Validate access key
	Validate(autKey string) bool
}

type simpleAuth struct {
	store driver.Driver
}

//NewSimpleAuth returns a simple auth instance
func NewSimpleAuth(s driver.Driver) Auth {
	return &simpleAuth{store: s}
}

//Create implementing Auth.Create
func (a *simpleAuth) Create(projectName string) (string, error) {
	lists := accessKey{}
	i, err := a.store.Get(storeKey)
	if err == nil {
		l, ok := i.(accessKey)
		if ok {
			lists = l
		}
	}

	//check project name exits
	for _, p := range lists {
		if p == projectName {
			return "", fmt.Errorf("project: %s already used", projectName)
		}
	}

	//simply create uuid from google uuid lib
	token := uuid.New().String()
	lists[token] = projectName

	if err := a.store.Set(storeKey, lists); err != nil {
		return "", err
	}

	return token, nil
}

//Validate implementing Auth.Validate
func (a *simpleAuth) Validate(authKey string) bool {
	i, err := a.store.Get(storeKey)
	if err != nil {
		return false
	}

	lists, ok := i.(accessKey)
	if !ok {
		return false
	}

	_, ok = lists[authKey]
	return ok
}
