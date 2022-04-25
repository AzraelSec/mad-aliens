package world

import (
	"fmt"
	"strings"

	"github.com/AzraelSec/mad-aliens/pkg/alien"
	"github.com/AzraelSec/mad-aliens/pkg/city"
	"github.com/AzraelSec/mad-aliens/pkg/utils"
)

type LinkMap map[string]map[Direction]*city.City
type CityMap map[string]*city.City
type AliensMap = map[int]*alien.Alien

type World struct {
	Cities CityMap
	Links  LinkMap

	Aliens      AliensMap
	StuckAliens AliensMap
}

func NewWorld(cities CityMap, links LinkMap, aliens, stuckAliens AliensMap) *World {
	return &World{
		Cities:      cities,
		Links:       links,
		Aliens:      aliens,
		StuckAliens: stuckAliens,
	}
}

func (w *World) findAlienPointer(id int) (*alien.Alien, error) {
	if alien, exists := w.Aliens[id]; exists {
		return alien, nil
	} else {
		return nil, fmt.Errorf("cannot find requested alien: %d", id)
	}
}

func (w *World) findCityPointer(name string) (*city.City, error) {
	if city, exists := w.Cities[name]; exists {
		return city, nil
	} else {
		return nil, fmt.Errorf("cannot find requested city: %s", name)
	}
}

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

func (w *World) RandomlyMove(id int) (string, error) {
	alien, err := w.findAlienPointer(id)
	if err != nil {
		return "", err
	}

	availableLinks, availableIds := w.Links[alien.City.Name], make([]string, 0)
	for link, arrival := range availableLinks {
		if arrival.Destroyed {
			delete(availableLinks, link)
		} else {
			availableIds = append(availableIds, arrival.Name)
		}
	}

	target := alien.City
	if len(availableLinks) == 0 {
		w.StuckAliens[alien.Id] = alien
	} else {
		idx := utils.RandomInt(len(availableIds))
		target = w.Cities[availableIds[idx]]
		alien.City = target
	}

	return target.Name, nil
}

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

func (w *World) GetAliensByCity(name string) []*alien.Alien {
	aliens := []*alien.Alien{}
	for _, alien := range w.Aliens {
		if alien.City.Name == name {
			aliens = append(aliens, alien)
		}
	}
	return aliens
}

func (w *World) DestroyAliens(ids []int) error {
	for _, id := range ids {
		if _, err := w.findAlienPointer(id); err != nil {
			return err
		} else {
			delete(w.Aliens, id)
			delete(w.StuckAliens, id)
		}
	}
	return nil
}

func (w *World) GetCity(name string) (city.City, error) {
	city, err := w.findCityPointer(name)
	return *city, err
}

func (w *World) GetAlien(id int) (alien.Alien, error) {
	alien, err := w.findAlienPointer(id)
	return *alien, err
}
