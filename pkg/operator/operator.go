package operator



type IOperator interface {
	Reconcile(key string)error
	Healthy()error
}


