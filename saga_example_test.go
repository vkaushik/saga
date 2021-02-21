package saga

func Example_sagaTransaction() {
	// storage := kafka.NewKafka(version, brokerAddresses, options...)
	// saga := NewSaga()
	// saga.AddSubTx("debit", DebitAccount, CompensateDebitAccount)
	// saga.AddSubTx("credit, "CreditAccount, CompensateCreditAccount)
	// tx := NewTx(ctx, saga, storage, logger, txID)
	// tx.Start()
	// tx.ExecuteAndGetResults("initialize transfer", from, to)
	// tx.Execute("debit", account, amount) -> if err do -> tx.Abort
	// tx.Execute("credit", credit, amount) -> if err do -> tx.Abort
	// tx.End()
}

// storage will provide - inmemory implementation as well
// repo will have golint, test, coverage, formatting etc. enabled
// a Dockerfile -> to construct image to run all go commands in local - container
// makefile -> hooks etc.
// I guess kafka is only using the configs, otherwise everything is taking args.
//

// Outside you'll interact with 3 objects -> Saga, Storage, Transaction
// Add SubTx definitions to Saga - action and compensate
// Storage is nothing but the persistence for Saga. Create Kafka instance with custom configs and use as storage.
// Create Transaction using the defined Saga and Storage.
// Use Transaction to execute SubTransactions and perform Abort if something goes unexpected.
