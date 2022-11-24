package event

import (
	"errors"
	"fmt"
	"time"
)

type Event struct {
	ID          int       `json:"id"`
	Date        time.Time `event:"date" json:"date"`
	Name        string    `event:"name" json:"name"`
	Description string    `event:"description" json:"description"`
}

type Result struct {
	result []Event `json:"result"`
}

type eventList struct {
	events map[int]Event
	lastID int
}

type EventManager map[int]*eventList

func NewEventManager() *EventManager {
	ev := make(EventManager)
	return &ev
}

func (ue *EventManager) SetEvent(userID int, ev Event) {
	if _, found := (*ue)[userID]; !found {
		(*ue)[userID] = &eventList{
			events: make(map[int]Event),
			lastID: 0,
		}
	}

	ev.ID = (*ue)[userID].lastID

	(*ue)[userID].events[(*ue)[userID].lastID] = ev
	(*ue)[userID].lastID++
}

func (ue *EventManager) UpdateEvent(userID int, id int, ev Event) error {
	if _, found := (*ue)[userID]; !found {
		return errors.New("User has not events")
	}
	if _, found := (*ue)[userID].events[id]; !found {
		return errors.New(fmt.Sprintf("User has not event ID=%d", id))
	}

	(*ue)[userID].events[id] = ev
	return nil
}

func (ue *EventManager) DeleteEvent(userID int, id int) error {
	if _, found := (*ue)[userID]; !found {
		return errors.New("User has not events")
	}
	if _, found := (*ue)[userID].events[id]; !found {
		return errors.New(fmt.Sprintf("User has not event ID=%d", id))
	}

	delete((*ue)[userID].events, id)
	return nil
}

// from, to include
func (ue *EventManager) GetEvents(userID int, from time.Time, to time.Time) []Event {
	if _, found := (*ue)[userID]; !found {
		return []Event{}
	}

	result := []Event{}

	for _, event := range (*ue)[userID].events {
		if from.Unix() <= event.Date.Unix() && event.Date.Unix() < to.Unix() {
			result = append(result, event)
		}
	}

	return result
}
