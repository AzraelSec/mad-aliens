package world

import (
	"github.com/AzraelSec/mad-aliens/pkg/alien"
	"github.com/AzraelSec/mad-aliens/pkg/city"
)

type AliensMap = map[int]*alien.Alien
type CityMap map[string]*city.City
type LinkMap map[string]map[Direction]*city.City

type World struct {
	Cities          CityMap   // A map that associates each city name to its object
	Links           LinkMap   // A map that assiciates each city name to its available links
	Aliens          AliensMap // A map that associates each alien name to its alien
	StuckAliens     int       // A stuck aliens counter
	DestroyedAliens int       // A destroyed aliens counter
}
