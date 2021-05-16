package tx

import (
	"context"
	"github.com/vkaushik/saga/marshal"
	"github.com/vkaushik/saga/subtx"
	"github.com/vkaushik/saga/trace"
	"reflect"
	"time"

	"github.com/juju/errors"
	"github.com/vkaushik/saga/log"
)

// Tx is the Transaction object to perfrom sub-transactions in the Saga
type Tx struct {
	ctx     context.Context
	txID    string // ID is the unique identifier for each transaction.
	saga    Saga
	storage Storage
	log     trace.Logger
}

// Saga is the dependency for Transaction that keeps Sub-Transaction definitions.
type Saga interface {
	GetSubTxDef(subTxID string) (subtx.Definition, error)
	MarshallArgs(args []interface{}) ([]log.ArgData, error)
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
	SetLogger(l trace.Logger)
	SetContext(ctx context.Context)
	RollbackWithInfiniteTries()
	Rollback(tryCount int) error
	IsTxIDAlreadyInUse() (bool, error)
}

// New returns an instance of type ReadyTx
func New(ctx context.Context, sg Saga, st Storage, txID string) ReadyTx {

	return NewWithLogger(ctx, sg, st, txID, trace.NewDummyLogger())
}

// NewWithLogger returns an instance of type ReadyTx with given logger
func NewWithLogger(ctx context.Context, sg Saga, st Storage, txID string, logger trace.Logger) ReadyTx {

	return &Tx{ctx: ctx, saga: sg, storage: st, txID: txID, log: logger}
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
			return errors.Annotatef(err, "could not rollback TxID: %s", tx.txID)
		}
	}

	logData, err := marshal.Marshal(logMsg)
	if err != nil {
		return errors.Annotatef(err, "could not start Tx: %s, because the log: %v is not serializable", tx.txID, logMsg)
	}

	if err = tx.storage.AppendLog(string(tx.txID), logData); err != nil {
		return errors.Annotatef(err, "could not append logs to storage for TxID: %s", tx.txID)
	}

	return nil
}

// ExecSubTx executes the sub-transaction that's already defined in saga and identified by the identifier
func (tx *Tx) ExecSubTx(subTxID string, args ...interface{}) error {
	_, err := tx.ExecSubTxAndGetResult(subTxID, args)
	return err
}

func getErrorFrom(result []reflect.Value) error {
	if result[0].IsNil() {
		return nil
	}
	return result[0].Interface().(error)
}

// ExecSubTxAndGetResult executes and returns the results of the sub-transaction that's already defined in saga and identified by the identifier
func (tx *Tx) ExecSubTxAndGetResult(subTxID string, args ...interface{}) ([]reflect.Value, error) {
	var res []reflect.Value

	// validate SubTxID and get the definition from saga
	subTxDef, err := tx.saga.GetSubTxDef(subTxID)
	if err != nil {
		return res, errors.Annotatef(err, "could not get SubTx definition for subTxID: %s", subTxID)
	}

	// log the starting of subTx
	marshalledArgs, err := tx.saga.MarshallArgs(args)
	if err != nil {
		return res, errors.Annotatef(err, "could not marshal params: %v", args)
	}

	logMsg := &log.Log{
		Type:    log.StartSubTx,
		SubTxID: subTxID,
		Time:    time.Now(),
		Args:    marshalledArgs,
	}

	l, err := marshal.Marshal(logMsg)
	if err != nil {
		return res, errors.Annotate(err, "could not marshal log message for start of SubTx")
	}

	if err := tx.storage.AppendLog(subTxID, l); err != nil {
		return res, errors.Annotate(err, "could not append start SubTx log for subTxID: "+subTxID)
	}

	// prepare actual arguments to execute SubTx
	actualArgs := make([]reflect.Value, 0, len(args)+1) // +1 for context which is first arg
	actualArgs = append(actualArgs, reflect.ValueOf(tx.ctx))
	for _, arg := range args {
		actualArgs = append(actualArgs, reflect.ValueOf(arg))
	}

	// execute subTx
	tx.log.Info("calling action for SubTxID: %s", subTxID)
	res = subTxDef.GetAction().Call(actualArgs)
	err = getErrorFrom(res)
	if err != nil {
		return res, errors.Annotatef(err, "subTx action execution returned error for subTxID: %s", subTxID)
	}

	// log the end of subTx action
	logMsg = &log.Log{
		Type:    log.EndSubTx,
		SubTxID: subTxID,
		Time:    time.Now(),
	}
	l, err = marshal.Marshal(logMsg)
	if err != nil {
		return res, errors.Annotate(err, "could not marshal log message for end of SubTx")
	}

	err = tx.storage.AppendLog(subTxID, l)
	if err != nil {
		return res, errors.Annotate(err, "could not append end SubTx log for subTxID: "+subTxID)
	}

	return res, nil
}

// End ends the Transaction
func (tx *Tx) End() error {
	logMsg := &log.Log{
		Type: log.EndTx,
		Time: time.Now(),
	}
	logData, err := marshal.Marshal(logMsg)
	if err != nil {
		return errors.Annotatef(err, "could not end Tx: %s, because the log: %v is not serializable", tx.txID, logMsg)
	}

	err = tx.storage.AppendLog(tx.txID, logData)
	if err != nil {
		return errors.Annotatef(err, "could not append end SubTx log for TxID: %s", tx.txID)
	}

	// Cleanup

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
func (tx *Tx) SetLogger(l trace.Logger) {
	tx.log = l
}

// SetContext to change the Transaction context.
func (tx *Tx) SetContext(ctx context.Context) {
	tx.ctx = ctx
}

// IsTxIDAlreadyInUse checks if TxID is already in use i.e. another Tx was already initiated.
func (tx *Tx) IsTxIDAlreadyInUse() (bool, error) {
	return tx.storage.TxIDAlreadyExists(string(tx.txID))
}
