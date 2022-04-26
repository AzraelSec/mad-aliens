package world

import (
	"fmt"
	"strings"

	"github.com/AzraelSec/mad-aliens/pkg/alien"
	"github.com/AzraelSec/mad-aliens/pkg/city"
	"github.com/AzraelSec/mad-aliens/pkg/utils"
)

func NewWorld(cities CityMap, links LinkMap, aliens AliensMap) *World {
	return &World{
		Cities:          cities,
		Links:           links,
		Aliens:          aliens,
		StuckAliens:     0,
		DestroyedAliens: 0,
	}
}

// Utility method to retrieve an alien pointer
func (w *World) findAlienPointer(id int) (*alien.Alien, error) {
	if alien, exists := w.Aliens[id]; exists {
		return alien, nil
	} else {
		return nil, fmt.Errorf("cannot find requested alien: %d", id)
	}
}

// Utility method to retrieve a city pointer
func (w *World) findCityPointer(name string) (*city.City, error) {
	if city, exists := w.Cities[name]; exists {
		return city, nil
	} else {
		return nil, fmt.Errorf("cannot find requested city: %s", name)
	}
}

// Method that counts the alive aliens.
// Since destroyed aliens are soft-deleted an additional function is needed to count the alive aliens.
func (w *World) CountAliveAliens() int {
	return len(w.Aliens) - w.DestroyedAliens
}

// Method that serialize the world with the same format used to define it
func (w *World) String() (o string) {
	lines := make([]string, 0)

	for name, city := range w.Cities {
		if city.Destroyed {
			continue
		}

		if links, exists := w.Links[name]; exists {
			if len(links) == 0 {
				continue
			}

			line := make([]string, 0)
			anyValid := false

			for direction, arrival := range links {
				directionStr, err := directionString(direction)
				if err == nil && !arrival.Destroyed {
					line = append(line, fmt.Sprintf("%s=%s", directionStr, arrival.Name))
					anyValid = true
				}
			}

			// Only cities having at least one link are printed out
			if anyValid {
				lines = append(lines, fmt.Sprintf("%s %s", name, strings.Join(line, " ")))
			}
		}
	}

	if len(lines) == 0 {
		return "No entry-point cities left"
	} else {
		return strings.Join(lines, "\n")
	}
}

// Method that moves an alien into a new random city following the available links.
// If no moves are available, the alien is stuck.
func (w *World) RandomlyMove(id int) (bool, error) {
	alien, err := w.findAlienPointer(id)
	if err != nil {
		return false, err
	}

	// Iterate over linked cities filtering destroyed ones to get the available next moves
	availableLinks, availableIds := w.Links[alien.City.Name], make([]string, 0)
	for _, arrival := range availableLinks {
		if !arrival.Destroyed {
			availableIds = append(availableIds, arrival.Name)
		}
	}

	if len(availableIds) == 0 {
		// If no moves are available, the alien is stuck
		w.StuckAliens++
		alien.Stuck = true
		return false, nil
	} else {
		// Else the alien moves into a new city
		idx := utils.RandomInt(len(availableIds))
		alien.City = w.Cities[availableIds[idx]]
		return true, nil
	}
}

// Method that soft-delete a group of cities
func (w *World) DestroyCities(names []string) error {
	for _, name := range names {
		if city, err := w.findCityPointer(name); err != nil {
			return err
		} else {
			city.Destroyed = true
		}
	}
	return nil
}

// Method that soft-delete a group of aliens
func (w *World) DestroyAliens(ids []int) error {
	for _, id := range ids {
		if al, err := w.findAlienPointer(id); err != nil {
			return err
		} else {
			w.DestroyedAliens++
			al.Destroyed = true
		}
	}
	return nil
}

// Utility method to retrieve a city from the cities map
func (w *World) GetCity(name string) (city.City, error) {
	city, err := w.findCityPointer(name)
	return *city, err
}
