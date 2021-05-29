package saga

import (
	"context"
	"errors"
	"fmt"
	"github.com/vkaushik/saga/storage/memory"
	"github.com/vkaushik/saga/trace"
	"github.com/vkaushik/saga/tx"
	"log"
	"testing"
)

// TODO: Add more tracing in saga.
func TestSaga(t *testing.T) {
	// create logger for tracing to help debug
	loggerForTx := trace.NewSimpleLogger()

	// create storage (
	storageForTx := memory.NewLogStorage()
	// kafka storage could look like this:
	//storageForTx, err := kafka.New("2.7.0", []string{"10.9.2.7"})
	//if err != nil {
	//	t.Fatal(err)
	//}

	// create saga and register subTx definitions with it
	sagaForTx := NewWithLogger(loggerForTx)
	if err := sagaForTx.AddSubTx("debit", debit, debitCompensate); err != nil {
		log.Println("Failed to add Debit SubTx to Saga")
		t.Fatal(err)
	}

	if err := sagaForTx.AddSubTx("credit", credit, creditCompensate); err != nil {
		log.Println("Failed to add Credit SubTx to Saga")
		t.Fatal(err)
	}

	// create a new transaction
	readyTx := tx.NewWithLogger(context.Background(), sagaForTx, storageForTx, "transfer-100-from-sam-to-pam", loggerForTx)
	defer func() {
		if err := readyTx.End(); err != nil {
			log.Println("Failed to close transaction.", err.Error())
			t.Fatal(err)
		}
	}()

	if err := readyTx.ExecSubTx("debit", 100, "sam"); err != nil {
		readyTx.RollbackWithInfiniteTries()
		t.Fatal(err)
	}
	if err := readyTx.ExecSubTx("credit", 100, "pam"); err != nil {
		readyTx.RollbackWithInfiniteTries()
		t.Fatal(err)
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

func creditError(c context.Context, amount int, to string) error {
	fmt.Printf("crediting amount: %v to account: %v\n", amount, to)
	return errors.New("error on credit")
}

// TODO: Add more tracing in saga.
func TestSagaError(t *testing.T) {
	// create logger for tracing to help debug
	loggerForTx := trace.NewSimpleLogger()

	// create storage (
	storageForTx := memory.NewLogStorage()
	// kafka storage could look like this:
	//storageForTx, err := kafka.New("2.7.0", []string{"10.9.2.7"})
	//if err != nil {
	//	t.Fatal(err)
	//}

	// create saga and register subTx definitions with it
	sagaForTx := NewWithLogger(loggerForTx)
	if err := sagaForTx.AddSubTx("debit", debit, debitCompensate); err != nil {
		log.Println("Failed to add Debit SubTx to Saga")
		t.Fatal(err)
	}

	if err := sagaForTx.AddSubTx("credit", creditError, creditCompensate); err != nil {
		log.Println("Failed to add Credit SubTx to Saga")
		t.Fatal(err)
	}

	// create a new transaction
	readyTx := tx.NewWithLogger(context.Background(), sagaForTx, storageForTx, "transfer-100-from-sam-to-pam", loggerForTx)
	defer func() {
		if err := readyTx.End(); err != nil {
			log.Println("Failed to close transaction.", err.Error())
			t.Fatal(err)
		}
	}()

	if err := readyTx.ExecSubTx("debit", 100, "sam"); err != nil {
		readyTx.RollbackWithInfiniteTries()
		t.Fatal(err)
	}
	if err := readyTx.ExecSubTx("credit", 100, "pam"); err != nil {
		readyTx.RollbackWithInfiniteTries()
		t.Log("This is expected to rollback " + err.Error())
	}
}
