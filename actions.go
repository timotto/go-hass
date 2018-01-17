package hass

import (
	"errors"
	"strings"
	"time"
)

// CheckAPI checks whether or not the API is running. It returns an error
// if it is not running.
func (a *Access) CheckAPI() error {
	response := struct {
		Message string `json:"message"`
	}{}
	err := a.httpGet("/api/", &response)
	if err != nil {
		return err
	}

	if response.Message == "" {
		return errors.New("hass: API is not running")
	}

	return nil
}

// State is the struct for an object state
type State struct {
	Attributes struct {
		Auto         bool   `json:"auto"`
		FriendlyName string `json:"friendly_name"`
		Hidden       bool   `json:"hidden"`
		Order        int    `json:"order"`
		AssumedState bool   `json:"assumed_state"`
	} `json:"attributes"`
	EntityID    string    `json:"entity_id"`
	LastChanged time.Time `json:"last_changed"`
	LastUpdated time.Time `json:"last_updated"`
	State       string    `json:"state"`
}

// States is an array of State objects
type States []State

// StateChange is used for changing state on an entity
type StateChange struct {
	EntityID string `json:"entityid"`
	State    string `json:"state"`
}

// FireEvent fires an event.
func (a *Access) FireEvent(eventType string, eventData interface{}) error {
	return a.httpPost("/api/events/"+eventType, eventData)
}

// CallService calls a service with a domain, service, and entity id.
func (a *Access) CallService(domain, service string, entityID string) error {
	serviceData := struct {
		EntityID string `json:"entity_id"`
	}{entityID}

	return a.httpPost("/api/services/"+domain+"/"+service, serviceData)
}

// ListStates gets an array of state objects
func (a *Access) ListStates() (s States, err error) {
	var list States
	err = a.httpGet("/api/states", &list)
	if err != nil {
		return States{}, err
	}
	return list, nil
}

// GetState retrieves one stateobject for the entity id
func (a *Access) GetState(id string) (s State, err error) {
	var state State
	err = a.httpGet("/api/states/"+id, &state)
	if err != nil {
		return State{}, err
	}
	return state, nil
}

// FilterStates returns a list of states filtered by the list of domains
func (a *Access) FilterStates(domains ...string) (s States, err error) {
	list, err := a.ListStates()
	if err != nil {
		return States{}, err
	}
	for d := range list {
		dom := strings.TrimSuffix(strings.SplitAfter(list[d].EntityID, ".")[0], ".")
		for _, fdom := range domains {
			if fdom == dom {
				s = append(s, list[d])
			}
		}
		if err != nil {
			panic(err)
		}
	}

	return s, err
}

// ChangeState changes the state of a device
func (a *Access) ChangeState(id, state string) (s State, err error) {
	s.EntityID = id
	s.State = state
	err = a.httpPost("/api/states/"+id, s)
	return State{}, err
}
