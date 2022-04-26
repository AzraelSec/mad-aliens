package engine

import (
	"testing"

	"github.com/AzraelSec/mad-aliens/pkg/alien"
	"github.com/AzraelSec/mad-aliens/pkg/city"
	"github.com/AzraelSec/mad-aliens/pkg/utils"
	"github.com/AzraelSec/mad-aliens/pkg/world"
)

// This test is intended to stress the engine in order to eventually find pitfalls
func TestMultipleRun(t *testing.T) {
	for i := 0; i < 100; i++ {
		e := Engine{
			MaxRuns: 100,
			Runs:    0,
			World:   randomWorld(),
		}

		_, err := e.Run()
		if err != nil {
			t.Errorf("No error expected, got %v", err)
		}
	}
}

func randomWorld() *world.World {
	var (
		nAliens = utils.RandomInt(4) + 1
		nCities = utils.RandomInt(10) + 1
		cities  = make(world.CityMap)
		cityIds = make([]string, 0)
		links   = make(world.LinkMap)
		aliens  = make(world.AliensMap)
	)

	for i := 1; i <= nCities; i++ {
		city := city.NewCity(utils.RandomString(utils.RandomInt(5) + 1))
		cities[city.Name] = city
		cityIds = append(cityIds, city.Name)
	}

	for _, source := range cities {
		if utils.RandomBool() {
			targetName := cityIds[utils.RandomInt(len(cityIds))]
			target := cities[targetName]
			if links[source.Name] != nil {
				links[source.Name][world.North] = target
			} else {
				links[source.Name] = map[world.Direction]*city.City{world.North: target}
			}
		}
	}

	for i := 0; i < nAliens; i++ {
		targetName := cityIds[utils.RandomInt(len(cityIds))]
		target := cities[targetName]
		aliens[i] = alien.NewAlien(i, target)
	}

	return world.NewWorld(
		cities,
		links,
		aliens,
	)
}

func TestCompleted(t *testing.T) {
	var (
		cityA = city.NewCity("A")
		cityB = city.NewCity("B")
		cityC = city.NewCity("C")
		a1    = alien.NewAlien(0, cityA)
	)

	var tests = []struct {
		engineItem  *Engine
		expectError bool
		expectValue ExecutionStatus
	}{
		// All aliens are stuck
		{
			engineItem: &Engine{
				World: world.NewWorld(
					world.CityMap{"A": cityA},
					world.LinkMap{},
					world.AliensMap{0: a1},
				),
				MaxRuns: 1,
				Runs:    0,
			},
			expectError: false,
			expectValue: ALL_ALIENS_STUCK,
		},
		// No alive aliens left
		{
			engineItem: &Engine{
				World: world.NewWorld(
					world.CityMap{"A": cityA},
					world.LinkMap{},
					world.AliensMap{},
				),
				MaxRuns: 1,
				Runs:    0,
			},
			expectError: false,
			expectValue: NO_ALIENS_LEFT,
		},
		// Max round reached
		{
			engineItem: &Engine{
				World: world.NewWorld(
					world.CityMap{"A": cityA, "B": cityB},
					world.LinkMap{"A": map[world.Direction]*city.City{world.North: cityB}},
					world.AliensMap{0: a1},
				),
				MaxRuns: 1,
				Runs:    0,
			},
			expectError: false,
			expectValue: MAX_ROUND_REACHED,
		},
		// Still running
		{
			engineItem: &Engine{
				World: world.NewWorld(
					world.CityMap{"A": cityA, "B": cityB, "C": cityC},
					world.LinkMap{
						"A": map[world.Direction]*city.City{world.North: cityB},
						"B": map[world.Direction]*city.City{world.South: cityC},
					},
					world.AliensMap{0: a1},
				),
				MaxRuns: 2,
				Runs:    0,
			},
			expectError: false,
			expectValue: RUNNING,
		},
	}

	for _, test := range tests {
		status, err := test.engineItem.tick()
		if test.expectError && err == nil {
			t.Error("Expected error, got nil")
		}

		if test.expectValue != status {
			t.Errorf("Expected %d, got %d", test.expectValue, status)
		}
	}
}

func TestHandleFight(t *testing.T) {
	var (
		cityA = city.NewCity("A")
		a1    = alien.NewAlien(0, cityA)
		a2    = alien.NewAlien(1, cityA)
		a3    = alien.NewAlien(2, cityA)

		cities = world.CityMap{"A": cityA}
		Links  = world.LinkMap{}
		Aliens = world.AliensMap{a1.Id: a1, a2.Id: a2, a3.Id: a3}
	)

	cityA.Destroyed = false

	w := world.NewWorld(cities, Links, Aliens)
	e := &Engine{
		World:   w,
		MaxRuns: 1,
		Runs:    0,
	}

	var tests = []struct {
		alien1, alien2 int
		city           string
		expectError    bool
		aliveAliens    int
	}{
		// A1 and A2 destroy the city and also B dies during the fight
		{a1.Id, a2.Id, cityA.Name, false, 1},
		// B exists even if the city has been destroyed since it could still be about to leave
		{a1.Id, a2.Id, "B", true, 1},
	}

	for _, test := range tests {
		if err := e.handleFight(test.alien1, test.alien2, test.city); err != nil && !test.expectError {
			t.Errorf("Expected no error, got %v", err)
		}

		if test.expectError {
			continue
		}

		if nAliens := w.CountAliveAliens(); nAliens != test.aliveAliens {
			t.Errorf("Expected %d aliens, got %d", test.aliveAliens, nAliens)
		}

		if c, err := e.World.GetCity("A"); err != nil {
			t.Errorf("An error occurred during city A retrieving: %v", c)
		} else if !c.Destroyed {
			t.Errorf("Expected city A to be destroyed, got %t", c.Destroyed)
		}
	}
}
