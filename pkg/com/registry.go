package com

import (
	"fmt"
	"path"
	"reflect"
	"runtime"
	"strings"
	"sync"
)

var (
	defaultRegistry = newRegistry()
)

type registry struct {
	sync.Mutex
	components []component
	options    []option
	config     ConfigProvider
}

type component struct {
	name string
	obj  interface{}
}

type option struct {
	section     string
	name        string
	dfault      interface{}
	description string
}

// Option ...
func Option(name string, dfault interface{}, description string) interface{} {
	return option{
		name:        name,
		dfault:      dfault,
		description: description,
	}
}

func Register(name string, com interface{}, options ...interface{}) {
	defaultRegistry.Register(name, com, options...)
}
func (r *registry) Register(name string, com interface{}, options ...interface{}) {
	r.Lock()
	defer r.Unlock()
	if _, exists := findComponent(r.components, name); exists {
		panic("component '" + name + "' already registered")
	}
	r.components = append(r.components, component{name, com})
	for _, opt := range options {
		o := opt.(option)
		o.section = shortName(name) // TODO: shortest name
		r.options = append(r.options, o)
	}
}

func Select(name string, iface interface{}) interface{} {
	return defaultRegistry.Select(name, iface)
}
func (r *registry) Select(name string, iface interface{}) interface{} {
	r.Lock()
	defer r.Unlock()
	if !r.config.ComponentEnabled(name) {
		return nil
	}
	com, ok := findComponent(r.components, name)
	if !ok {
		return nil
	}
	if iface != nil {
		ifaceType := reflect.TypeOf(iface).Elem()
		typ := reflect.TypeOf(com.obj)
		if ifaceType.Kind() == reflect.Func && typ.AssignableTo(ifaceType) {
			return com.obj
		}
		if ifaceType.Kind() != reflect.Func && typ.Implements(ifaceType) {
			return com.obj
		}
		return nil
	}
	return com.obj
}

func Enabled(iface interface{}, ctx Context) []interface{} {
	return defaultRegistry.Enabled(iface, ctx)
}
func (r *registry) Enabled(iface interface{}, ctx Context) []interface{} {
	r.Lock()
	defer r.Unlock()
	if iface == nil {
		var none []interface{}
		return none
	}
	ifaceType := reflect.TypeOf(iface).Elem()
	var coms []interface{}
	for _, com := range r.components {
		if !r.config.ComponentEnabled(com.name) {
			continue
		}
		if ctx != nil && !ctx.ComponentEnabled(com.name) {
			continue
		}
		if iface == nil {
			coms = append(coms, com.obj)
		} else {
			typ := reflect.TypeOf(com.obj)
			if ifaceType.Kind() == reflect.Func && typ.AssignableTo(ifaceType) {
				coms = append(coms, com.obj)
			}
			if ifaceType.Kind() != reflect.Func && typ.Implements(ifaceType) {
				coms = append(coms, com.obj)
			}
		}
	}
	return coms
}

func SetConfig(c ConfigProvider) { defaultRegistry.SetConfig(c) }
func (r *registry) SetConfig(c ConfigProvider) {
	r.config = c
}

func GetString(name string) string {
	return defaultRegistry.GetString(pkgToSection(callerPkg()), name)
}
func (r *registry) GetString(section, name string) string {
	fqn := fmt.Sprintf("%s.%s", section, name)
	value, set := r.config.GetString(fqn)
	if set {
		return value
	}
	r.Lock()
	defer r.Unlock()
	for _, option := range r.options {
		if option.section == section && option.name == name {
			return option.dfault.(string)
		}
	}
	return ""
}

func GetInt(name string) int {
	return defaultRegistry.GetInt(pkgToSection(callerPkg()), name)
}
func (r *registry) GetInt(section, name string) int {
	fqn := fmt.Sprintf("%s.%s", section, name)
	value, set := r.config.GetInt(fqn)
	if set {
		return value
	}
	r.Lock()
	defer r.Unlock()
	for _, option := range r.options {
		if option.section == section && option.name == name {
			return option.dfault.(int)
		}
	}
	return 0
}

func GetBool(name string) bool {
	return defaultRegistry.GetBool(pkgToSection(callerPkg()), name)
}
func (r *registry) GetBool(section, name string) bool {
	fqn := fmt.Sprintf("%s.%s", section, name)
	value, set := r.config.GetBool(fqn)
	if set {
		return value
	}
	r.Lock()
	defer r.Unlock()
	for _, option := range r.options {
		if option.section == section && option.name == name {
			return option.dfault.(bool)
		}
	}
	return false
}

func findComponent(v []component, name string) (*component, bool) {
	for _, com := range v {
		if com.name == name {
			return &com, true
		}
	}
	return nil, false
}

func shortName(componentName string) string {
	parts := strings.Split(componentName, ".")
	return parts[len(parts)-1]
}

func callerPkg() string {
	pc := make([]uintptr, 10)
	runtime.Callers(3, pc) // this caller's caller
	f := runtime.FuncForPC(pc[0]).Name()
	base := path.Base(f)
	dir := path.Dir(f)
	parts := strings.Split(base, ".")
	return path.Join(dir, parts[0])
}

func pkgToSection(pkg string) string {
	parts := strings.Split(pkg, "/")
	return parts[len(parts)-1]
}

func newRegistry() registry {
	return registry{config: mapConfig{}}
}
