package city

type City struct {
	Name      string
	Destroyed bool
}

func NewCity(name string) *City {
	return &City{
		Name:      name,
		Destroyed: false,
	}
}
