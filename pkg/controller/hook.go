package controller


type IHook interface {
	OnAdd(obj interface{})
	OnUpdate(obj interface{})
	OnDelete(obj interface{})
}
