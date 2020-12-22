package capabilities

//SimpleCapability represent
type SimpleCapability struct {
	Callable func() error `json:"-"`
}
