package hass

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type EventListener struct {
	reader io.ReadCloser
	buffer *bufio.Reader
}

func (a *Access) ListenEvents() (*EventListener, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)

	return a.ListenEventsWithContext(ctx)
}

func (a *Access) ListenEventsWithContext(ctx context.Context) (*EventListener, error) {
	client := http.DefaultClient

	req, err := http.NewRequestWithContext(ctx, "GET", a.host+"/api/stream", nil)
	if err != nil {
		return nil, err
	}

	a.authorizeRequest(req)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return &EventListener{
		reader: resp.Body,
		buffer: bufio.NewReader(resp.Body),
	}, nil
}

type StateChangedEvent struct {
	Origin    string    `json:"origin"`
	EventType string    `json:"event_type"`
	TimeFired time.Time `json:"time_fired"`
	Data      struct {
		OldState struct {
			EntityID    string          `json:"entity_id"`
			State       string          `json:"state"`
			LastChanged time.Time       `json:"last_changed"`
			LastUpdated time.Time       `json:"last_updated"`
			Attributes  StateAttributes `json:"attributes"`
		} `json:"old_state"`
		EntityID string `json:"entity_id"`
		NewState struct {
			EntityID    string          `json:"entity_id"`
			State       string          `json:"state"`
			LastChanged time.Time       `json:"last_changed"`
			LastUpdated time.Time       `json:"last_updated"`
			Attributes  StateAttributes `json:"attributes"`
		} `json:"new_state"`
	} `json:"data"`
}

// NextStateChanged waits and returns for the next state_changed event.
func (e *EventListener) NextStateChanged() (StateChangedEvent, error) {
	for {
		line, err := e.buffer.ReadBytes('\n')
		if err != nil {
			return StateChangedEvent{}, err
		}

		if len(line) > 6 && string(line[:6]) == "data: " {
			jsonData := line[6:]

			if string(jsonData) == "ping\n" {
				continue
			}

			var eventTypeFinder struct {
				EventType string `json:"event_type"`
			}

			err := json.Unmarshal(jsonData, &eventTypeFinder)
			if err != nil {
				return StateChangedEvent{}, err
			}

			if eventTypeFinder.EventType != "state_changed" {
				continue
			}

			var stateChanged StateChangedEvent
			err = json.Unmarshal(jsonData, &stateChanged)
			if err != nil {
				return StateChangedEvent{}, err
			}

			return stateChanged, nil
		}
	}
}

// Close closes the event listener library
func (e *EventListener) Close() error {
	return e.reader.Close()
}
