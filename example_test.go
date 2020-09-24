package saga

import (
	"context"
	"log"
	"testing"

	"github.com/vkaushik/saga/storage"
)

func TestSaga(t *testing.T) {
	st := storage.NewKafka()
	s := New(st)
	defer func() {
		if err := s.Close(); err != nil {
			log.Println("Failed to close saga.")
		}
	}()

	if err := s.AddSubTx("debit", debit, debitCompensate); err != nil {
		log.Println("Failed to add Debit SubTx to Saga")
	}

	if err := s.AddSubTx("credit", credit, creditCompensate); err != nil {
		log.Println("Failed to add Credit SubTx to Saga")
	}

	tx := s.NewTx(context.Background(), "transfer-100")
	defer func() {
		if err := tx.End(); err != nil {
			log.Println("Failed to End the Transaction")
		}
	}()

	if err := tx.ExecSubTx("debit", 100, "sam"); err != nil {
		tryRollback(tx, err)
	}
	if err := tx.ExecSubTx("credit", 100, "pam"); err != nil {
		tryRollback(tx, err)
	}
}

// tryRollback to try Rollback infinitely until success
func tryRollback(tx *Tx, err error) {
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
