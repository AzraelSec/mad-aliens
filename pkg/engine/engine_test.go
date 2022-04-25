package engine

import (
	"testing"

	"github.com/AzraelSec/mad-aliens/pkg/alien"
	"github.com/AzraelSec/mad-aliens/pkg/city"
	"github.com/AzraelSec/mad-aliens/pkg/utils"
	"github.com/AzraelSec/mad-aliens/pkg/world"
)

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

	return world.NewWorld(
		cities,
		links,
		aliens,
		world.AliensMap{},
	)
}

func TestCompleted(t *testing.T) {
	var (
		cityA = city.NewCity("A")
		cityB = city.NewCity("B")
		cityC = city.NewCity("C")

		a1 = &alien.Alien{
			Id:   0,
			City: cityA,
		}
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
					world.AliensMap{},
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
					world.AliensMap{},
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
					world.AliensMap{},
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

		a1 = &alien.Alien{
			Id:   0,
			City: cityA,
		}
		a2 = &alien.Alien{
			Id:   1,
			City: cityA,
		}
		a3 = &alien.Alien{
			Id:   2,
			City: cityA,
		}

		cities        = world.CityMap{"A": cityA}
		Links         = world.LinkMap{}
		Aliens        = world.AliensMap{0: a1, 1: a2, 2: a3}
		StuckedAliens = world.AliensMap{}
	)

	cityA.Destroyed = false

	w := world.NewWorld(cities, Links, Aliens, StuckedAliens)
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
		{a1.Id, a2.Id, cityA.Name, false, 0},
		{a1.Id, a2.Id, "B", true, 0},
	}

	for _, test := range tests {
		err := e.handleFight(test.alien1, test.alien2, test.city)

		if err != nil && !test.expectError {
			t.Errorf("Expected no error, got %v", err)
		}

		if nAliens := len(w.Aliens); nAliens != test.aliveAliens {
			t.Errorf("Expected %d aliens, got %d", test.aliveAliens, nAliens)
		}

		if c, err := e.World.GetCity("A"); err != nil {
			t.Errorf("An error occurred during city A retrieving: %v", c)
		} else if !c.Destroyed {
			t.Errorf("Expected city A to be destroyed, got %t", c.Destroyed)
		}
	}
}
