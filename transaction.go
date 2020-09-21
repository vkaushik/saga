package saga

// Tx is the Transaction composed of multiple Sub-Transactions.
type Tx struct {
}

func newTx() *Tx {
	return &Tx{}
}

// ExecSubTx to execute the Sub-Transaction identified by SubTxID.
// The SubTxID and the corresponding Action & Compensate definitions must already be registered with Saga.
// The arguments to ExecSubTx are passed to the Sub-Transaction Action.
func (tx *Tx) ExecSubTx(ID SubTxID, args ...interface{}) error {
	return nil
}

// Rollback to stop transaction execution and bring back the stystem to consistent state.
// Rollback will call Compensate for all Sub-Transactions which were started.
func (tx *Tx) Rollback() error {
	return nil
}

// End to delete the Storage data for Transaction.
func (tx *Tx) End() error {
	return nil
}