package engine

import (
	"io"
	"log"

	"github.com/AzraelSec/mad-aliens/pkg/world"
)

// Method to instance a new Engine
func NewEngine(nAliens int, mRounds int, in io.Reader) (*Engine, error) {
	world, err := world.Parse(in, nAliens)
	if err != nil {
		return nil, err
	}

	return &Engine{
		MaxRuns: mRounds,
		Runs:    0,
		World:   world,
	}, nil
}

// Method that starts the evaluation loop.
// The execution ends when all the aliens are dead, all the aliens are stuck or the max number of runs is reached.
func (e *Engine) Run() (ExecutionStatus, error) {
	completed, err := ExecutionStatus(RUNNING), error(nil)
	for err == nil && completed == RUNNING {
		log.Printf("======= Run #%d =======\n", e.Runs)
		completed, err = e.tick()
	}
	return completed, err
}

// Method that performs a single execution step
func (e *Engine) tick() (ExecutionStatus, error) {
	// If already completed, return the previous result
	if status := e.completed(); status != RUNNING {
		return status, nil
	}

	// Array that keeps track of the cities that gets visited and the id of the visitor alien
	visited := make(map[string]int)

	for _, alien := range e.World.Aliens {
		// If alien has been destroyed, skip
		if alien.Destroyed {
			continue
		}

		// Identify a random move that each alien will perform
		currentCityName := alien.City.Name
		moved, err := e.World.RandomlyMove(alien.Id)
		if err != nil {
			log.Fatal(err)
		}

		if !moved {
			log.Printf("Alien %d is in a city %s with no available directions", alien.Id, alien.City.Name)

			/*
			* Assume that a city has no links (aliens get stuck) and two aliens, A and B, land on it and destroy it.
			* If a third alien C already was on that same city it results to be on an already destroyed city.
			 */
			if alien.City.Destroyed {
				if err := e.World.DestroyAliens([]int{alien.Id}); err != nil {
					log.Fatal(err)
				}

			}
		} else {
			log.Printf("Alien %d is in city %s and moving to %s", alien.Id, currentCityName, alien.City.Name)
		}

		/*
		* If arrival city has already been visited before during this turn, check if a fight is required
		*
		* Note: this implementation assumes that two aliens fight only if they LAND on the same city.
		* If an alien moves to a city that already has another alien whose move has not been evaluated yet,
		* no fight is performed.
		 */
		if _, exists := visited[alien.City.Name]; exists {
			if err := e.handleFight(alien.Id, visited[alien.City.Name], alien.City.Name); err != nil {
				log.Fatal(err)
			} else {
				delete(visited, alien.City.Name)
			}
		} else {
			// Else register that visit
			visited[alien.City.Name] = alien.Id
		}
	}

	// Increment the runs counter and return the ending condition
	e.Runs++
	return e.completed(), nil
}

// Method that checks if at least an ending condition is met
func (e *Engine) completed() ExecutionStatus {
	aliveAliens, stuckAliens := e.World.CountAliveAliens(), e.World.StuckAliens
	if aliveAliens == 0 {
		return NO_ALIENS_LEFT // There are no alive aliens left
	}
	if stuckAliens >= aliveAliens {
		return ALL_ALIENS_STUCK // All alive aliens are stuck
	}
	if e.Runs >= e.MaxRuns {
		return MAX_ROUND_REACHED // Max number of runs reached
	}
	return RUNNING
}

// Method to handle a fight between two aliens that land on the same city.
func (e *Engine) handleFight(a1 int, a2 int, city string) error {
	// Destroy both aliens that fought
	if err := e.World.DestroyAliens([]int{a1, a2}); err != nil {
		return err
	}

	// Destroy the involved city
	if err := e.World.DestroyCities([]string{city}); err != nil {
		return err
	}

	log.Printf("%s has been destroyed by alien %d and alien %d!", city, a1, a2)
	return nil
}
