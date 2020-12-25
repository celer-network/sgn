package types

// Params returns all of the global params
type Params struct {
}

func (gp Params) String() string {
	return ""
}

// NewParams creates a new global Params instance
func NewParams() Params {
	return Params{}
}

// DefaultParams default global params
func DefaultParams() Params {
	return NewParams()
}
