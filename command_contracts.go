package types

import "time"

// CommandRequestPayload marks a domain payload type that can be used in a typed
// command request.
type CommandRequestPayload interface {
	CommandRequestPayloadKind() string
}

// CommandResponsePayload marks a domain payload type that can be used in a typed
// command response.
type CommandResponsePayload interface {
	CommandResponsePayloadKind() string
}

// CommandRequest is a typed command contract used inside the framework/plugin
// boundary after payload decoding.
type CommandRequest[P CommandRequestPayload] struct {
	CommandID string `json:"command_id"`
	PluginID  string `json:"plugin_id"`
	Device    Device `json:"device"`
	Entity    Entity `json:"entity"`
	Payload   P      `json:"payload"`
}

// CommandResponse is a typed command result contract.
type CommandResponse[P CommandResponsePayload] struct {
	Device  Device `json:"device"`
	Entity  Entity `json:"entity"`
	Payload P      `json:"payload"`
}

// EventPayload marks a typed payload used in events/inbound events.
type EventPayload interface {
	EventPayloadKind() string
}

type InboundEventTyped[P EventPayload] struct {
	DeviceID      string `json:"device_id"`
	EntityID      string `json:"entity_id"`
	CorrelationID string `json:"correlation_id,omitempty"`
	Payload       P      `json:"payload"`
}

type EventTyped[P EventPayload] struct {
	ID            string    `json:"id"`
	PluginID      string    `json:"plugin_id"`
	DeviceID      string    `json:"device_id"`
	EntityID      string    `json:"entity_id"`
	EntityType    string    `json:"entity_type"`
	CorrelationID string    `json:"correlation_id,omitempty"`
	Payload       P         `json:"payload"`
	CreatedAt     time.Time `json:"created_at"`
}
