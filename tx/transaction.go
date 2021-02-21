package tx

import (
	"context"
	"reflect"
	"time"

	"github.com/juju/errors"
	"github.com/vkaushik/saga"
	"github.com/vkaushik/saga/log"
)

// ID is the unique identifier for each transaction.
type ID string

// Tx is the Transaction object to perfrom sub-transactions in the Saga
type Tx struct {
	ctx     context.Context
	txID    ID
	saga    Saga
	storage Storage
	log     saga.Logger
}

// Saga is the depencency for Transaction that keeps Sub-Transaction definitions.
type Saga interface {
}

// Storage is the dependency for Transaction that provides persistence to Saga
type Storage interface {
	TxIDAlreadyExists(string) (bool, error)
	AppendLog(string, string) error
}

// ReadyTx is the Ready-Transaction that exposes the executable actions for the Saga Transaction
type ReadyTx interface {
	Start() error
	ExecSubTx(subTxID string, args ...interface{}) error
	ExecSubTxAndGetResult(subTxID string, args ...interface{}) ([]reflect.Value, error)
	End() error
	SetLogger(l saga.Logger)
	SetContext(ctx context.Context)
	RollbackWithInfiniteTries()
	Rollback(tryCount int) error
}

// New returns an instance of type ReadyTx
func New(ctx context.Context, sg Saga, st Storage, txID ID) ReadyTx {

	return &Tx{ctx: ctx, saga: sg, storage: st, txID: txID, log: saga.NewDummyLogger()}
}

// Start starts the transaction
// If it returns any error like "could not rollback TxID", you can retry start, because the start tries rollback just once.
func (tx *Tx) Start() error {
	logMsg := &log.Log{
		Type: log.StartTx,
		Time: time.Now(),
	}

	if txIDAlreadyExists, err := tx.storage.TxIDAlreadyExists(string(tx.txID)); err != nil {
		return errors.Annotate(err, "could not find if the TxID is already in use")
	} else if txIDAlreadyExists {
		tx.log.Info("TxID is already in use, calling rollback on this TxID: %s to avoid any inconsistencies", tx.txID)
		if err = tx.Rollback(1); err != nil {
			return errors.Annotatef(err, "could not rollback TxID: ", string(tx.txID))
		}
	}

	logData, err := log.Marshal(logMsg)
	if err != nil {
		return errors.Annotatef(err, "could not start Tx: %s, because the log: %v is not serializable", tx.txID, logMsg)
	}

	if err = tx.storage.AppendLog(string(tx.txID), logData); err != nil {
		return errors.Annotatef(err, "could not append logs to storage for TxID: ", tx.txID)
	}

	return nil
}

// ExecSubTx executes the sub-transaction that's already defined in saga and identified by the identifier
func (tx *Tx) ExecSubTx(subTxID string, args ...interface{}) error {
	return nil
}

// ExecSubTxAndGetResult executes and returns the results of the sub-transaction that's already defined in saga and identified by the identifier
func (tx *Tx) ExecSubTxAndGetResult(subTxID string, args ...interface{}) ([]reflect.Value, error) {
	return nil, nil
}

// End ends the Transaction
func (tx *Tx) End() error {
	return nil
}

// RollbackWithInfiniteTries tries rolling back the transaction. It'll keep retrying the rollback until it's successful.
func (tx *Tx) RollbackWithInfiniteTries() {

}

// Rollback tries rolling back the transaction. If rollback is failing it'll try rolling back only tryCount times.
func (tx *Tx) Rollback(tryCount int) error {
	return nil
}

// SetLogger to change the Transaction logger.
func (tx *Tx) SetLogger(l saga.Logger) {
	tx.log = l
}

// SetContext to change the Transaction context.
func (tx *Tx) SetContext(ctx context.Context) {

}
