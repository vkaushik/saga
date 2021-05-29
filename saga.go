package saga

import (
	"github.com/juju/errors"
	"github.com/vkaushik/saga/log"
	"github.com/vkaushik/saga/marshal"
	"github.com/vkaushik/saga/subtx"
	"github.com/vkaushik/saga/trace"
	"reflect"
)

// New creates and returns a new Saga instance
func New() *Saga {
	return NewWithLogger(trace.NewDummyLogger())
}

// NewWithLogger creates and returns a new Saga instance with given Logger
func NewWithLogger(l trace.Logger) *Saga {
	return &Saga{
		log:      l,
		subTxDef: subtx.NewSubTxDefinitions(),
		params:   subtx.NewParamTypeRegister(),
	}
}

// Saga helps SubTransaction execution and rollback
type Saga struct {
	log      trace.Logger
	subTxDef SubTxDefinitions
	params   ParamRegister
}

// SubTxDefinitions contains methods to add sub-transaction definitions
type SubTxDefinitions interface {
	Add(subTxID string, action interface{}, compensate interface{}) error
	Get(subTxID string) (subtx.Definition, error)
}

// ParamRegister contains methods to add sub-transaction parameters metadata
type ParamRegister interface {
	Add(funcObj interface{}) error
	GetRegisteredTypeName(t reflect.Type) (typ string, err error)
	GetRegisteredType(typ string) (t reflect.Type, err error)
}

// AddSubTx registers the action and compensate methods for a SubTx that'll be identified with the SubTxID.
// While Transaction execution, the SubTxID is used to identify the SubTx and execute it's action in success flow
// or compensate if Tx is being rollback.
func (s *Saga) AddSubTx(ID string, action interface{}, compensate interface{}) error {
	if err := s.params.Add(action); err != nil {
		return errors.Annotatef(err, "could not parse action parameters for SubTxID: %s", ID)
	}

	if err := s.params.Add(compensate); err != nil {
		return errors.Annotatef(err, "could not parse compensate parameters for SubTxID: %s", ID)
	}

	if err := s.subTxDef.Add(string(ID), action, compensate); err != nil {
		return errors.Annotate(err, "could not add sub-transaction definitions")
	}

	return nil
}

func (s *Saga) GetSubTxDef(subTxID string) (subtx.Definition, error) {
	return s.subTxDef.Get(subTxID)
}

func (s *Saga) MarshallArgs(args []interface{}) ([]log.ArgData, error) {
	res := make([]log.ArgData, 0, len(args))

	for _, arg := range args {
		t, err := s.params.GetRegisteredTypeName(reflect.ValueOf(arg).Type())
		if err != nil {
			return res, errors.Annotate(err, "could not find argument type registered in saga")
		}

		m, err := marshal.Marshal(arg)
		if err != nil {
			return res, errors.Annotate(err, "could not marshal arg")
		}

		ad := log.ArgData{
			Type:  t,
			Value: m,
		}

		res = append(res, ad)
	}

	return res, nil
}

func (s *Saga) UnmarshallArgs(argData []log.ArgData) ([]reflect.Value, error) {
	var actualArgs []reflect.Value
	for _, arg := range argData {
		typ, err := s.params.GetRegisteredType(arg.Type)
		if err != nil {
			return actualArgs, errors.Annotate(err, "could not find argument registered")
		}
		obj := reflect.New(typ).Interface()
		err = marshal.Unmarshal([]byte(arg.Value), obj)
		if err != nil {
			return actualArgs, errors.Annotatef(err, "could not unmarshal type: %s", arg.Type)
		}
		objValue := reflect.ValueOf(obj)
		// if pointer, get its value
		if objValue.Type().Kind() == reflect.Ptr && objValue.Type() != typ {
			objValue = objValue.Elem()
		}
		actualArgs = append(actualArgs, objValue)
	}

	return actualArgs, nil
}
