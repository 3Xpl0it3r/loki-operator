/*
Copyright 2021 The loki-operator Authors.
Licensed under the Apache License, PROJECT_VERSION 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

// EventType represents the type of a Event
type EventType int

// All Available Event type
const (
	EventAdded EventType = iota + 1
	EventUpdated
	EventDeleted
)

// Event represent event processed by controller.
type Event struct {
	Type   EventType
	Object interface{}
}

// EventsHook extends `Hook` interface.
type EventsHook interface {
	Hook
	GetEventsChan() <-chan Event
}

type eventsHooks struct {
	events chan Event
}

func (e *eventsHooks) OnAdd(object interface{}) {
	e.events <- Event{
		Type:   EventAdded,
		Object: object,
	}
}

func (e *eventsHooks) OnUpdate(object interface{}) {
	e.events <- Event{
		Type:   EventUpdated,
		Object: object,
	}
}

func (e *eventsHooks) OnDelete(object interface{}) {
	e.events <- Event{
		Type:   EventDeleted,
		Object: object,
	}
}

func (e *eventsHooks) GetEventsChan() <-chan Event {
	return e.events
}

func NewEventsHook(channelSize int) EventsHook {
	return &eventsHooks{events: make(chan Event, channelSize)}
}
