package log

import (
	"time"
)

// Type is the type of log being logged in the storage
type Type int

const (
	// StartTx denotes the start of a Transaction
	StartTx Type = iota + 1

	// StartSubTx denotes the start of a Sub-Transaction
	StartSubTx

	// EndSubTx denotes
	EndSubTx

	// StartCompensateSubTx denotes the start  of a Sub-Transaction
	StartCompensateSubTx

	// EndCompensateSubTx denotes the end of a Sub-Transaction
	EndCompensateSubTx

	// EndTx denotes the end of a Transaction
	EndTx

	// AbortTx denotes the abort of a Transaction
	AbortTx
)

// Log is used by Saga to compensate a SubTx. Logs persisted in storage are the source of truth to know the state of a Transaction.
type Log struct {
	Type    Type      `json:"type,omitempty"`
	SubTxID string    `json:"sub_tx_ID,omitempty"`
	Time    time.Time `json:"time,omitempty"`
	Args    []ArgData `json:"args,omitempty"`
}

// ArgData is used by Log to contain the arguments passed to SubTx. It's used to store and restore SubTx input args from logs.
type ArgData struct {
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}
