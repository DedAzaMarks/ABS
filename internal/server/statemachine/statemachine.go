package statemachine

import (
	"errors"
	"fmt"
)

// State represents a statemachine in the statemachine machine
type State string

const (
	StateInit             State = "start"
	StateSearch           State = "search"
	StateFilmSelection    State = "film_selection"
	StateVersionSelection State = "version_selection"
	StateDeviceSelection  State = "device_selection"
	StateDone             State = "done"
)

// Event represents an event that triggers statemachine transitions
type Event string

const (
	EventCancel        Event = "cancel"
	EventNewSearch     Event = "new"
	EventSelectFilm    Event = "select_film"
	EventSelectVersion Event = "select_version"
	EventSelectDevice  Event = "select_device"
	EventFinish        Event = "finish"
)

var transitions = map[State]map[Event]State{
	StateInit: {
		EventNewSearch: StateSearch,
		EventCancel:    StateInit,
	},
	StateSearch: {
		EventNewSearch:  StateSearch,
		EventSelectFilm: StateFilmSelection,
		EventCancel:     StateInit,
	},
	StateFilmSelection: {
		EventNewSearch:     StateSearch,
		EventSelectVersion: StateVersionSelection,
		EventCancel:        StateInit,
	},
	StateVersionSelection: {
		EventNewSearch:    StateSearch,
		EventSelectDevice: StateDeviceSelection,
		EventFinish:       StateInit,
		EventCancel:       StateInit,
	},
	StateDeviceSelection: {
		EventNewSearch: StateSearch,
		EventFinish:    StateInit,
		EventCancel:    StateInit,
	},
}

// StateMachine represents the finite statemachine machine
type StateMachine struct {
	currentState State
	transitions  map[State]map[Event]State
}

// NewStateMachine initializes a new statemachine machine with the given initial statemachine
func NewStateMachine() *StateMachine {
	return &StateMachine{
		currentState: StateInit,
		transitions:  transitions,
	}
}

// TriggerEvent triggers an event and changes the state if the transition is valid
func (sm *StateMachine) TriggerEvent(event Event) error {
	if nextState, ok := sm.transitions[sm.currentState][event]; ok {
		fmt.Printf("Transition: %s -> %s on event %s\n", sm.currentState, nextState, event)
		sm.currentState = nextState
		return nil
	}
	return errors.New("invalid event for current state")
}

// CurrentState returns the current state of the state machine
func (sm *StateMachine) CurrentState() State {
	return sm.currentState
}

// Reset sets StateMachine to StateInit
func (sm *StateMachine) Reset() {
	sm.currentState = StateInit
}

// SetState should be used ONLY for setting state after deserialization
func SetState(state State) *StateMachine {
	return &StateMachine{
		currentState: state,
		transitions:  transitions,
	}
}
