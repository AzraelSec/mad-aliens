package engine

import "github.com/AzraelSec/mad-aliens/pkg/world"

// Constants that define the engine execution status
const (
	RUNNING           = iota + 1 // Still running
	MAX_ROUND_REACHED            // Completed because of max round reached
	ALL_ALIENS_STUCK             // All the aliens are stuck and next executions would not change
	NO_ALIENS_LEFT               // All the aliens are dead fighting
)

type ExecutionStatus int

// Engine that handles the world's events and define the way aliens and city should behave
type Engine struct {
	MaxRuns int          // Max number of execution rounds
	Runs    int          // Number of execution rounds already performed
	World   *world.World // Pointer to the world the engine should manage
}
