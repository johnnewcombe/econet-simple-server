package piconet

type TxResult struct {
	Result string
	Ok     bool
}

func NewTxResult(result string) TxResult {

	txResult := TxResult{
		Result: result,
	}
	if result == "OK" {
		txResult.Ok = true
	}
	return txResult
}
