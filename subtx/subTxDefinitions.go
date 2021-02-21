package subtx

import (
	"context"
	"reflect"

	"github.com/juju/errors"
)

// SubTxDefinitions contains the metadata for each SubTransaction with their identifiers i.e. SubTxID
type definitions map[string]definition

type definition struct {
	subTxID    string
	action     reflect.Value
	compensate reflect.Value
}

func (d *definitions) add(subTxID string, action interface{}, compensate interface{}) error {
	actionFunc, err := validateAndGetFuncValue(action)
	if err != nil {
		return errors.Annotatef(err, "invalid action provided for SubTxID: %s", subTxID)
	}

	compensateFunc, err := validateAndGetFuncValue(compensate)
	if err != nil {
		return errors.Annotatef(err, "invalid compensate provided for SubTxID: %s", subTxID)
	}

	(*d)[subTxID] = definition{
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
	if isFirstArgumentNotAContext(funcValue) {
		return reflect.Value{}, errors.Errorf("first argument must be of type context.Context")
	}

	return funcValue, nil
}

func isNotAFunction(v reflect.Value) bool {
	return v.Kind() != reflect.Func
}

func isFirstArgumentNotAContext(v reflect.Value) bool {
	return v.Type().NumIn() < 1 || v.Type().In(0) != reflect.TypeOf((context.Context)(nil))
}
