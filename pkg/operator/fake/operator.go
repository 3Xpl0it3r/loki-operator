package fake

type fakeOperator struct {
}

func (f fakeOperator) Reconcile(obj interface{}) error {
	return nil
}

func NewOperator() *fakeOperator {
	return &fakeOperator{}
}
