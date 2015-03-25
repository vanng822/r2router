package r2router

import (

)

type Params interface {
	Get(key string) string
	AppSet(key string, val interface{})
	AppGet(key string) interface{}
}

// Holding value for named parameters
type params_ struct {
	requestParams map[string]string
	appData       map[string]interface{}
}

func (p *params_) Get(key string) string {
	return p.requestParams[key]
}

func (p *params_) AppSet(key string, val interface{}) {
	p.appData[key] = val
}

func (p *params_) AppGet(key string) interface{} {
	return p.appData[key]
}

