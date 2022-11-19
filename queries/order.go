package queries

type Order string

const (
	Desc Order = "desc"
	Asc  Order = "asc"
)

func (o Order) String() string {
	return string(o)
}

func (o Order) IsValid() bool {
	switch o {
	case Desc:
	case Asc:
	default:
		return false
	}
	return true
}
