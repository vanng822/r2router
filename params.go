package r2router

import (

)

// Params is for parameters that are matched from URL.
// It is also brigde to forward data from middleware.
// An Example could be that a middleware to identify
// the user API key, verify and get user data
type Params interface {
	// Get returns param value for the given key
	Get(key string) string
	// Has is for checking if parameter value exists
	// for given key
	Has(key string) bool
	// AppSet is for application to set own data
	AppSet(key interface{}, val interface{})
	// AppGet returns the value from AppSet
	AppGet(key interface{}) interface{}
	// AppHas is for checking if the application
	// has set data for given key
	AppHas(key interface{}) bool
}

// Holding value for named parameters
type params_ struct {
	requestParams map[string]string
	appData       map[interface{}]interface{}
}

func (p *params_) Get(key string) string {
	return p.requestParams[key]
}

func (p *params_) Has(key string) bool {
	_, exists := p.requestParams[key]
	return exists
}

func (p *params_) AppSet(key interface{}, val interface{}) {
	p.appData[key] = val
}

func (p *params_) AppGet(key interface{}) interface{} {
	return p.appData[key]
}

func (p *params_) AppHas(key interface{}) bool {
	_, exists := p.appData[key]
	return exists
}