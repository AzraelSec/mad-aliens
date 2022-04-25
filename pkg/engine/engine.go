package engine

import (
	"io"
	"log"

	"github.com/AzraelSec/mad-aliens/pkg/world"
)

const (
	RUNNING = iota + 1
	MAX_ROUND_REACHED
	ALL_ALIENS_STUCK
	NO_ALIENS_LEFT
)

type ExecutionStatus int

type Engine struct {
	MaxRuns int
	Runs    int

	World *world.World
}

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

func (e *Engine) Run() (ExecutionStatus, error) {
	completed, err := ExecutionStatus(RUNNING), error(nil)
	for err == nil && completed == RUNNING {
		log.Printf("======= Run #%d =======\n", e.Runs)
		completed, err = e.tick()
	}
	return completed, err
}

func (e *Engine) tick() (ExecutionStatus, error) {
	if status := e.completed(); status != RUNNING {
		return status, nil
	}

	visited := make(map[string]int)

	for _, alien := range e.World.Aliens {
		arrival, err := e.World.RandomlyMove(alien.Id)
		if err != nil {
			log.Fatal(err)
		} else {
			if alien.City.Name == arrival {
				log.Printf("Alien %d is in a city %s with no available directions", alien.Id, arrival)
			} else {
				log.Printf("Alien %d is in city %s and moving to %s", alien.Id, alien.City.Name, arrival)
			}
		}

		if _, exists := visited[arrival]; exists {
			if err := e.handleFight(alien.Id, visited[arrival], arrival); err != nil {
				log.Fatal(err)
			} else {
				delete(visited, arrival)
			}
		} else {
			visited[arrival] = alien.Id
		}
	}

	e.Runs++
	return e.completed(), nil
}

func (e *Engine) completed() ExecutionStatus {
	aliveAliens, stuckAliens := e.World.Aliens, e.World.StuckAliens
	if len(aliveAliens) == 0 {
		return NO_ALIENS_LEFT
	}
	if len(stuckAliens) >= len(aliveAliens) {
		return ALL_ALIENS_STUCK
	}
	if e.Runs >= e.MaxRuns {
		return MAX_ROUND_REACHED
	}
	return RUNNING
}

func (e *Engine) handleFight(a1 int, a2 int, city string) error {
	ids := []int{}
	for _, alien := range e.World.GetAliensByCity(city) {
		ids = append(ids, alien.Id)
	}

	if err := e.World.DestroyAliens(ids); err != nil {
		return err
	}

	if err := e.World.DestroyCities([]string{city}); err != nil {
		return err
	}

	log.Printf("%s has been destroyed by alien %d and alien %d!", city, a1, a2)
	return nil
}
