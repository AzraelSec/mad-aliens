package world

import (
	"testing"

	"github.com/AzraelSec/mad-aliens/pkg/alien"
	"github.com/AzraelSec/mad-aliens/pkg/city"
)

var (
	cityA = city.NewCity("A")
	cityB = city.NewCity("B")
	cityC = city.NewCity("C")

	destroyed1 = city.NewCity("D")
	destroyed2 = city.NewCity("E")
	destroyed3 = city.NewCity("F")
)

func init() {
	destroyed1.Destroyed = true
	destroyed2.Destroyed = true
	destroyed3.Destroyed = true
}

func TestString(t *testing.T) {
	var tests = []struct {
		cities CityMap
		links  LinkMap
		output string
	}{
		// Empty world => empty string
		{
			cities: CityMap{},
			links:  LinkMap{},
			output: "No entry-point cities left",
		},
		// Prints one city per line
		// note: cities B and C are not placed on separate lines since their links are empty
		{
			cities: CityMap{
				"A": cityA,
				"B": cityB,
				"C": cityC,
			},
			links:  LinkMap{"A": map[Direction]*city.City{North: cityB}},
			output: "A north=B",
		},
		// Prints only cities that has not been destroyed
		{
			cities: CityMap{
				"A": cityA,
				"B": cityB,
				"D": destroyed1,
				"E": destroyed2,
			},
			links: LinkMap{
				"A": map[Direction]*city.City{South: cityB},
				"B": map[Direction]*city.City{North: destroyed1},
				"D": map[Direction]*city.City{North: destroyed2},
			},
			output: "A south=B",
		},
	}

	for _, test := range tests {
		w := World{
			Cities: test.cities,
			Links:  test.links,
		}

		if w.String() != test.output {
			t.Errorf("Expected %s, got %s", test.output, w.String())
		}
	}
}

func TestRandomlyMove(t *testing.T) {
	var tests = []struct {
		cities           CityMap
		links            LinkMap
		startCity        *city.City
		expectedCity     *city.City
		expectedStuckNum int
	}{
		// Aliens successfully move to B
		{
			cities: CityMap{
				"A": cityA,
				"B": cityB,
			},
			links: LinkMap{
				"A": map[Direction]*city.City{North: cityB},
				"B": map[Direction]*city.City{South: cityA},
			},
			startCity:        cityA,
			expectedCity:     cityB,
			expectedStuckNum: 0,
		},
		// Alien cannot move anywhere since its current location has no links, so it stays in A
		{
			cities:           CityMap{"A": cityA},
			links:            LinkMap{},
			startCity:        cityA,
			expectedCity:     cityA,
			expectedStuckNum: 1,
		},
		// Alien cannot move onto destroyed city
		{
			cities: CityMap{
				"A": cityA,
				"D": destroyed1,
			},
			links:            LinkMap{"A": map[Direction]*city.City{North: destroyed1}},
			startCity:        cityA,
			expectedCity:     cityA,
			expectedStuckNum: 1,
		},
	}

	for _, test := range tests {
		a := alien.NewAlien(0, test.startCity)
		w := NewWorld(
			test.cities,
			test.links,
			AliensMap{a.Id: a},
		)

		_, err := w.RandomlyMove(a.Id)
		if a.City.Name != test.expectedCity.Name {
			t.Errorf("Expected %s, got %s", test.expectedCity.Name, a.City.Name)
		}

		if err != nil {
			t.Error(err)
		}

		if w.StuckAliens != test.expectedStuckNum {
			t.Errorf("Expected %v stuck aliens, got %v", test.expectedStuckNum, w.StuckAliens)
		}
	}
}
