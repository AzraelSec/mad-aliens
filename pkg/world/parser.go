package world

import (
	"bufio"
	"errors"
	"io"
	"log"
	"strings"

	"github.com/AzraelSec/mad-aliens/pkg/alien"
	"github.com/AzraelSec/mad-aliens/pkg/city"
	"github.com/AzraelSec/mad-aliens/pkg/utils"
)

// Method that parses a world definition and randomly
// position a given number of aliens inside of them
func Parse(in io.Reader, nAliens int) (*World, error) {
	var (
		scanner = bufio.NewScanner(in)
		cities  = make(CityMap)
		links   = make(LinkMap)

		// This slice keeps track of the created city names.
		// In this way, no additional loops are lately performed to deploy the aliens.
		ids = make([]string, 0)
	)

	for scanner.Scan() {
		line := scanner.Text()
		items := strings.Split(strings.TrimSpace(line), " ")
		sourceName, directions := items[0], items[1:]

		// Create a city and add it to the cities map
		attachCity(city.NewCity(sourceName), cities)
		ids = append(ids, sourceName)

		// The links map uses the directions as a key since a single link
		// can exist per each direction per each city.
		// Because of this, if multiple links with the same directions are defined
		// for the same source, the last one is used.
		for _, directionConfig := range directions {
			sep := strings.Index(directionConfig, "=")
			if sep == -1 {
				return nil, errors.New("invalid direction: " + directionConfig)
			}

			directionName := directionConfig[:sep]
			direction, ok := parseDirection(directionName)
			if !ok {
				return nil, errors.New("invalid direction: " + directionName)
			}

			targetName := directionConfig[sep+1:]
			ids = append(ids, targetName)

			var target *city.City
			if found, exists := cities[targetName]; exists {
				target = found
			} else {
				// If the linked city does not exist yet, it gets created
				target = city.NewCity(targetName)
			}

			// A new link between source and target city is added to the world links map
			attachCity(target, cities)
			addLink(sourceName, target, direction, links)
		}
	}

	return NewWorld(
		cities,
		links,
		deployAliens(
			nAliens,
			cities,
			ids,
		),
	), nil
}

// Method that adds a city to the cities map
func attachCity(ct *city.City, cts CityMap) {
	if _, exists := cts[ct.Name]; !exists {
		cts[ct.Name] = ct
	}
}

// Method that adds a link between two cities to the links map
func addLink(from string, to *city.City, direction Direction, links LinkMap) {
	if links[from] == nil {
		links[from] = map[Direction]*city.City{}
	}
	links[from][direction] = to
}

// Method that randomly define an aliens map in order to deploy alines in random locations
func deployAliens(nAliens int, cp CityMap, cityNames []string) AliensMap {
	aliens := make(map[int]*alien.Alien, nAliens)
	for i := 0; i < nAliens; i++ {
		idx := utils.RandomInt(len(cityNames))
		aliens[i] = alien.NewAlien(i, cp[cityNames[idx]])
		log.Printf("Alien %d is located at %s", aliens[i].Id, aliens[i].City.Name)
	}
	return aliens
}
