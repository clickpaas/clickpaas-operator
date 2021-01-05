package operator



type IOperator interface {
	Sync(key string)error
	Healthy()error
}


