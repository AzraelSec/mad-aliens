package city

type City struct {
	Name      string // City identification number (assigned during the world initialization)
	Destroyed bool   // When a city is destroyed, a soft-deleted is performed
}

// Function to instanciate a new City
func NewCity(name string) *City {
	return &City{
		Name:      name,
		Destroyed: false,
	}
}
