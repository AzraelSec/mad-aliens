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

func Parse(in io.Reader, nAliens int) (*World, error) {
	var (
		scanner = bufio.NewScanner(in)
		cities  = make(CityMap)
		links   = make(LinkMap)
		ids     = make([]string, 0)
	)

	for scanner.Scan() {
		line := scanner.Text()
		items := strings.Split(strings.TrimSpace(line), " ")
		sourceName, directions := items[0], items[1:]

		attachCity(city.NewCity(sourceName), cities)
		ids = append(ids, sourceName)

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
				target = city.NewCity(targetName)
			}

			attachCity(target, cities)
			addLink(sourceName, target, direction, links)
		}
	}

	return &World{
		Cities: cities,
		Links:  links,
		Aliens: deployAliens(
			nAliens,
			cities,
			ids,
		),
		StuckAliens: make(map[int]*alien.Alien),
	}, nil
}

func attachCity(ct *city.City, cts CityMap) {
	if _, exists := cts[ct.Name]; !exists {
		cts[ct.Name] = ct
	}
}

func addLink(from string, to *city.City, direction Direction, links LinkMap) {
	if links[from] == nil {
		links[from] = map[Direction]*city.City{}
	}
	links[from][direction] = to
}

func deployAliens(nAliens int, cp CityMap, cityNames []string) AliensMap {
	aliens := make(map[int]*alien.Alien, nAliens)
	for i := 0; i < nAliens; i++ {
		idx := utils.RandomInt(len(cityNames))
		aliens[i] = &alien.Alien{
			Id:   i,
			City: cp[cityNames[idx]],
		}
		log.Printf("Alien %d is located at %s", aliens[i].Id, aliens[i].City.Name)
	}
	return aliens
}
