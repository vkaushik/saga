package subtx

import (
	"reflect"

	"github.com/juju/errors"
)

type paramTypeRegister struct {
	nameToType map[string]reflect.Type
	typeToName map[reflect.Value]string
}

func (pr *paramTypeRegister) add(obj interface{}) error {
	funcValue, err := validateAndGetFuncValue(obj)
	if err != nil {
		return errors.Annotate(err, "invalid function")
	}

	funcType := funcValue.Type()
	pr.addInputParams(funcType)
	pr.addOutputParams(funcType)

	return nil
}

func (pr *paramTypeRegister) addInputParams(funcType reflect.Type) {
	for i := 0; i < funcType.NumIn(); i++ {
		pr.addParam(funcType.In(i))
	}
}

func (pr *paramTypeRegister) addOutputParams(funcType reflect.Type) {
	for i := 0; i < funcType.NumIn(); i++ {
		pr.addParam(funcType.Out(i))
	}
}

func (pr *paramTypeRegister) addParam(paramType reflect.Type) {
	// TODO: Check if taking value out of pointer got any side-effects on abort
	if paramType.Kind() == reflect.Ptr {
		paramType = paramType.Elem()
	}
	paramName := paramType.PkgPath() + "/" + paramType.Name()
	pr.nameToType[paramName] = paramType
}
