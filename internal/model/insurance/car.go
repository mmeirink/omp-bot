package insurance

type Car struct {
	Title string
}

func (c Car) String() string {
	return c.Title
}
