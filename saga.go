package saga

import (
	"context"
)

// New to create Saga Instance.
func New(s Storage) *Saga {
	return &Saga{}
}

// Saga to co-ordinate and control transaction.
type Saga struct {
	// Storage
	//
}

// Storage to provide persistence for Saga.
type Storage interface {
	// if topic exists
	// create topic
	// close/cleanup
}

// SubTxID is an unique identifier for each SubTx.
type SubTxID string

// Action to execute the steps/tasks for a SubTx.
type Action interface{}

// Compensate to rollback the SubTx to help system bring back to consistent state.
// Compensate has to be Idempotent and can be called more than once while rolling-back the Tx.
type Compensate interface{}

// AddSubTx to add Action and Compensate definitions for given SubTxID.
// Both Action and Compensate must have same signature, arguments must be serializable/marshalable and
// first argument must be the cotnext. (context is not marshalable and is kept just to widen the usability)
func (s *Saga) AddSubTx(id SubTxID, a Action, c Compensate) error {
	// Validate action and compensate - same signature, first argument context
	// Add definitions to maps: subTxId to compensate. paramType to paramName. paramName to paramType.

	return nil
}

// TxID is an unique identifier for the Transaction.
type TxID string

// NewTx to create a new Transaction
func (s *Saga) NewTx(c context.Context, txID TxID) *Tx {
	// If TxID storage topic already exists then rollback the Tx
	// else Create storage topic with new TxID

	return newTx()
}

// Close to perform cleanup tasks and close connections i.e. kafka-client.
func (s *Saga) Close() error {
	// storage close
	return nil
}
