package alien

import "github.com/AzraelSec/mad-aliens/pkg/city"

type Alien struct {
	Id        int        // Aliens identification number (assigned during the world initialization)
	City      *city.City // Pointer to the city the alien is currently located in
	Stuck     bool       // A boolean indicating if alien is stuck or not
	Destroyed bool       // A boolean indicating if alien is destroyed or not
}

// Function to instanciate a new Alien
func NewAlien(id int, city *city.City) *Alien {
	return &Alien{
		Id:        id,
		City:      city,
		Stuck:     false,
		Destroyed: false,
	}
}
