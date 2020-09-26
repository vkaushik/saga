package saga_test

import (
	"context"
	"log"

	"github.com/vkaushik/saga"
	"github.com/vkaushik/saga/storage"
)

func ExampleSaga() {
	// Create storage for persistence
	st := storage.NewKafka()

	// Create Saga
	s := saga.New(st)
	defer func() {
		if err := s.Close(); err != nil {
			log.Println("Failed to close saga.")
		}
	}()

	// Add SubTx Definitions to Saga. Saga uses these definitions for Rollback.
	// The action and compensate could be some member functions as well and must have same signature.
	// Also the arguments must be serializable/marshalable and first arg must be context.
	if err := s.AddSubTx("debit", debit, debitCompensate); err != nil {
		log.Println("Failed to add Debit SubTx to Saga")
	}

	if err := s.AddSubTx("credit", credit, creditCompensate); err != nil {
		log.Println("Failed to add Credit SubTx to Saga")
	}

	// Create New Transaction
	tx := s.NewTx(context.Background(), "transfer-100")
	defer func() {
		if err := tx.End(); err != nil {
			log.Println("Failed to End the Transaction")
		}
	}()

	// Execute Sub Transactions. Invoke transaction rollback if Sub-Transaction execution fails.
	// The Rollback call is made explicit because there may be some cases when Rollback is not required or
	// possibly somehing needs to be done before or after rollback.
	if err := tx.ExecSubTx("debit", 100, "sam"); err != nil {
		tryRollback(tx, err)
	}
	if err := tx.ExecSubTx("credit", 100, "pam"); err != nil {
		tryRollback(tx, err)
	}
}

// tryRollback to try Rollback infinitely until success
func tryRollback(tx *saga.Tx, err error) {
	log.Println("Initiating Rollback due to: ", err.Error())
	rollbackErr := tx.Rollback()

	for rollbackErr != nil {
		log.Println("Failed to Rollback with error", rollbackErr.Error())
		log.Println("Retrying Rollback.")
		rollbackErr = tx.Rollback()
	}
}

func debit(c context.Context, amount int, from string) {

}

func debitCompensate(c context.Context, amount int, from string) {

}

func credit(c context.Context, amount int, to string) {

}

func creditCompensate(c context.Context, amount int, to string) {

}
