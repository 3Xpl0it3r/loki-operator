package operator

type Operator interface {
	Reconcile(obj interface{}) error
}
