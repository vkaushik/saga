package subtx

import (
	"context"
	"reflect"

	"github.com/juju/errors"
)

func NewSubTxDefinitions() *Definitions {
	return &Definitions{}
}

// SubTxDefinitions contains the metadata for each SubTransaction with their identifiers i.e. SubTxID
type Definitions map[string]Definition

type Definition struct {
	subTxID    string
	action     reflect.Value
	compensate reflect.Value
}

func (d *Definition) GetAction() reflect.Value {
	return d.action
}

func (d *Definition) GetCompensate() reflect.Value {
	return d.compensate
}

func (d *Definitions) Get(subTxID string) (Definition, error) {
	if def, ok := (*d)[subTxID]; ok {
		return def, nil
	} else {
		return Definition{}, errors.New("could not find subTx definition for subTxID: " + subTxID)
	}
}

func (d *Definitions) Add(subTxID string, action interface{}, compensate interface{}) error {
	actionFunc, err := validateAndGetFuncValue(action)
	if err != nil {
		return errors.Annotatef(err, "invalid action provided for SubTxID: %s", subTxID)
	}

	compensateFunc, err := validateAndGetFuncValue(compensate)
	if err != nil {
		return errors.Annotatef(err, "invalid compensate provided for SubTxID: %s", subTxID)
	}

	(*d)[subTxID] = Definition{
		subTxID:    subTxID,
		action:     actionFunc,
		compensate: compensateFunc,
	}

	return nil
}

func validateAndGetFuncValue(obj interface{}) (reflect.Value, error) {
	funcValue := reflect.ValueOf(obj)
	if isNotAFunction(funcValue) {
		return reflect.Value{}, errors.Errorf("must be function")
	}
	if isNotFirstArgumentAContext(funcValue) {
		return reflect.Value{}, errors.Errorf("first argument must be of type context.Context")
	}
	if isNotFirstReturnValueOfTypeError(funcValue) {
		return reflect.Value{}, errors.Errorf("first return type must be of type error")
	}

	return funcValue, nil
}

func isNotAFunction(v reflect.Value) bool {
	return v.Kind() != reflect.Func
}

func isNotFirstArgumentAContext(v reflect.Value) bool {
	return v.Type().NumIn() < 1 || !v.Type().In(0).Implements(reflect.TypeOf((*context.Context)(nil)).Elem())
}

func isNotFirstReturnValueOfTypeError(v reflect.Value) bool {
	return v.Type().NumOut() == 0 || !v.Type().Out(0).Implements(reflect.TypeOf((*error)(nil)).Elem())
}
