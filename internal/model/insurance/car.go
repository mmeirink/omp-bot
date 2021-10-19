package insurance

type Car struct {
	ID uint64
	Title string
}

func (c Car) String() string {
	return c.Title
}
