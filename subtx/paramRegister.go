package subtx

import (
	"reflect"

	"github.com/juju/errors"
)

func NewParamTypeRegister() *ParamTypeRegister {
	return &ParamTypeRegister{
		nameToType: map[string]reflect.Type{},
		typeToName: map[reflect.Type]string{},
	}
}

type ParamTypeRegister struct {
	nameToType map[string]reflect.Type
	typeToName map[reflect.Type]string
}

func (pr *ParamTypeRegister) GetRegisteredType(t reflect.Type) (typ string, err error) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if p, ok := pr.typeToName[t]; ok {
		return p, nil
	}
	return "", errors.New("could not find the type in param type register: " + t.String())
}

func (pr *ParamTypeRegister) Add(obj interface{}) error {
	funcValue, err := validateAndGetFuncValue(obj)
	if err != nil {
		return errors.Annotate(err, "invalid function")
	}

	funcType := funcValue.Type()
	pr.addInputParams(funcType)
	pr.addOutputParams(funcType)

	return nil
}

func (pr *ParamTypeRegister) addInputParams(funcType reflect.Type) {
	for i := 0; i < funcType.NumIn(); i++ {
		pr.addParam(funcType.In(i))
	}
}

func (pr *ParamTypeRegister) addOutputParams(funcType reflect.Type) {
	for i := 0; i < funcType.NumIn(); i++ {
		pr.addParam(funcType.Out(i))
	}
}

func (pr *ParamTypeRegister) addParam(paramType reflect.Type) {
	// TODO: Check if taking value out of pointer got any side-effects on abort
	if paramType.Kind() == reflect.Ptr {
		paramType = paramType.Elem()
	}
	paramName := paramType.PkgPath() + "/" + paramType.Name()
	pr.nameToType[paramName] = paramType
}
