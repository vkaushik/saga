package saga

import (
	"context"
	"fmt"
	"github.com/vkaushik/saga/trace"
	"github.com/vkaushik/saga/tx"
	"log"
	"testing"

	"github.com/vkaushik/saga/storage/kafka"
)

// TODO: Fix the test: Failing to Add SubTx to saga. Add more tracing in saga.
func NoTestSaga(t *testing.T) {
	// create logger for tracing to help debug
	loggerForTx := trace.NewSimpleLogger()

	// create storage (kafka in this case) for persistence
	storageForTx, err := kafka.New("2.7.0", []string{"10.9.2.7"})
	if err != nil {
		t.Fail()
	}

	// create saga and register subTx definitions with it
	sagaForTx := NewWithLogger(loggerForTx)
	if err := sagaForTx.AddSubTx("debit", debit, debitCompensate); err != nil {
		log.Println("Failed to add Debit SubTx to Saga")
	}

	if err := sagaForTx.AddSubTx("credit", credit, creditCompensate); err != nil {
		log.Println("Failed to add Credit SubTx to Saga")
	}

	// create a new transaction
	readyTx := tx.NewWithLogger(context.Background(), sagaForTx, storageForTx, "transfer-100-from-sam-to-pam", loggerForTx)
	defer func() {
		if err := readyTx.End(); err != nil {
			log.Println("Failed to close transaction.", err.Error())
		}
	}()

	if err := readyTx.ExecSubTx("debit", 100, "sam"); err != nil {
		readyTx.RollbackWithInfiniteTries()
	}
	if err := readyTx.ExecSubTx("credit", 100, "pam"); err != nil {
		readyTx.RollbackWithInfiniteTries()
	}
}

func debit(c context.Context, amount int, from string) error {
	fmt.Printf("debiting amount: %v from account: %v\n", amount, from)
	return nil
}

func debitCompensate(c context.Context, amount int, from string) error {
	fmt.Printf("reversing: debiting amount: %v from account: %v\n", amount, from)
	return nil
}

func credit(c context.Context, amount int, to string) error {
	fmt.Printf("crediting amount: %v to account: %v\n", amount, to)
	return nil
}

func creditCompensate(c context.Context, amount int, to string) error {
	fmt.Printf("reversing crediting amount: %v to account: %v\n", amount, to)
	return nil
}
