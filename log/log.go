package log

import (
	"encoding/json"
	"time"

	"github.com/juju/errors"
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

// Marshal converts the value to json string
func Marshal(log *Log) (string, error) {
	if b, err := json.Marshal(*log); err != nil {
		return "", errors.Annotatef(err, "could not marshal the log %v", log)
	} else {
		return string(b), nil
	}
}

// Unmarshal unmarshals the data bytes and fills in the Log object
func Unmarshal(data []byte, log *Log) error {
	if err := json.Unmarshal(data, log); err != nil {
		return errors.Annotatef(err, "could not unmarshal the log data: %s", data)
	}

	return nil
}
