package saga

import (
	"github.com/juju/errors"
)

// New creates and returns a new Saga instance
func New(l Logger) *Saga {
	return &Saga{
		log: l,
	}
}

// Saga helps SubTranscation execution and rollback
type Saga struct {
	log      Logger
	subTxDef SubTxDefinitions
	params   ParamRegister
}

// SubTxDefinitions contains methods to add sub-transaction definitions
type SubTxDefinitions interface {
	add(string, interface{}, interface{}) error
}

// ParamRegister contains methods to add sub-transaction parameters metadata
type ParamRegister interface {
	add(funcObj interface{}) error
}

// AddSubTx registers the action and compensate methods for a SubTx that'll be identified with the SubTxID.
// While Transaction execution, the SubTxID is used to identify the SubTx and execute it's action in success flow
// or compensate if Tx is being rollbacked.
func (s *Saga) AddSubTx(ID string, action interface{}, compensate interface{}) error {
	if err := s.params.add(action); err != nil {
		return errors.Annotatef(err, "could not parse action parameters for SubTxID: ", ID)
	}

	if err := s.params.add(compensate); err != nil {
		return errors.Annotatef(err, "could not parse compensate parameters for SubTxID: ", ID)
	}

	if err := s.subTxDef.add(string(ID), action, compensate); err != nil {
		return errors.Annotate(err, "could not add sub-transaction definitions")
	}

	return nil
}
