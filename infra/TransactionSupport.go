package infra

type TransactionSupport interface {
	TxBegin()
	TxCommit()
	TxRollback()
}
