package saga

import (
	"github.com/juju/errors"
	"github.com/vkaushik/saga/log"
	"github.com/vkaushik/saga/marshal"
	"github.com/vkaushik/saga/subtx"
	"reflect"
)

// New creates and returns a new Saga instance
func New() *Saga {
	return NewWithLogger(NewDummyLogger())
}

// NewWithLogger creates and returns a new Saga instance with given Logger
func NewWithLogger(l Logger) *Saga {
	return &Saga{
		log:      l,
		subTxDef: subtx.NewSubTxDefinitions(),
		params:   subtx.NewParamTypeRegister(),
	}
}

// Saga helps SubTransaction execution and rollback
type Saga struct {
	log      Logger
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
	GetRegisteredType(t reflect.Type) (typ string, err error)
}

// AddSubTx registers the action and compensate methods for a SubTx that'll be identified with the SubTxID.
// While Transaction execution, the SubTxID is used to identify the SubTx and execute it's action in success flow
// or compensate if Tx is being rollback.
func (s *Saga) AddSubTx(ID string, action interface{}, compensate interface{}) error {
	if err := s.params.Add(action); err != nil {
		return errors.Annotatef(err, "could not parse action parameters for SubTxID: ", ID)
	}

	if err := s.params.Add(compensate); err != nil {
		return errors.Annotatef(err, "could not parse compensate parameters for SubTxID: ", ID)
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
		t, err := s.params.GetRegisteredType(reflect.ValueOf(arg).Type())
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
