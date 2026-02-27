package types

import (
	"encoding/json"
	"time"
)

const JSONRPCVersion = "2.0"

// --- Transport ---

type Request struct {
	JSONRPC string           `json:"jsonrpc"`
	ID      *json.RawMessage `json:"id,omitempty"`
	Method  string           `json:"method"`
	Params  json.RawMessage  `json:"params,omitempty"`
}

type Response struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *RPCError       `json:"error,omitempty"`
}

type RPCError struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// --- Registry ---

type Manifest struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Registration struct {
	Manifest   Manifest `json:"manifest"`
	RPCSubject string   `json:"rpc_subject"`
}

// --- Domain ---

type Storage struct {
	Meta string          `json:"meta"`
	Data json.RawMessage `json:"data"`
}

type Device struct {
	ID         string  `json:"id"`
	SourceID   string  `json:"source_id"`
	SourceName string  `json:"source_name"`
	LocalName  string  `json:"local_name"`
	Config     Storage `json:"config"`
}

func (d Device) Name() string {
	if d.LocalName != "" {
		return d.LocalName
	}
	if d.SourceName != "" {
		return d.SourceName
	}
	if d.SourceID != "" {
		return d.SourceID
	}
	return d.ID
}

func (d Device) MarshalJSON() ([]byte, error) {
	type Alias Device
	return json.Marshal(&struct {
		Alias
		Name string `json:"name"`
	}{
		Alias: (Alias)(d),
		Name:  d.Name(),
	})
}

type Entity struct {
	ID        string     `json:"id"`
	DeviceID  string     `json:"device_id"`
	Domain    string     `json:"domain"`
	LocalName string     `json:"local_name"`
	Config    Storage    `json:"config"`
	Actions   []string   `json:"actions,omitempty"`
	Data      EntityData `json:"data"`
}

type EntityData struct {
	Desired       json.RawMessage `json:"desired,omitempty"`
	Reported      json.RawMessage `json:"reported,omitempty"`
	Effective     json.RawMessage `json:"effective,omitempty"`
	SyncStatus    string          `json:"sync_status,omitempty"`
	LastCommandID string          `json:"last_command_id,omitempty"`
	LastEventID   string          `json:"last_event_id,omitempty"`
	UpdatedAt     time.Time       `json:"updated_at,omitempty"`
}

// --- Commands ---

// Command carries intent from the gateway to a plugin.
// Payload shape is defined by the entity package for the given Domain.
type Command struct {
	ID         string          `json:"id"`
	PluginID   string          `json:"plugin_id"`
	DeviceID   string          `json:"device_id"`
	EntityID   string          `json:"entity_id"`
	EntityType string          `json:"entity_type"`
	Payload    json.RawMessage `json:"payload"`
	CreatedAt  time.Time       `json:"created_at"`
}

type CommandState string

const (
	CommandPending   CommandState = "pending"
	CommandSucceeded CommandState = "succeeded"
	CommandFailed    CommandState = "failed"
)

type CommandStatus struct {
	CommandID     string       `json:"command_id"`
	PluginID      string       `json:"plugin_id"`
	DeviceID      string       `json:"device_id"`
	EntityID      string       `json:"entity_id"`
	EntityType    string       `json:"entity_type"`
	State         CommandState `json:"state"`
	Error         string       `json:"error,omitempty"`
	CreatedAt     time.Time    `json:"created_at"`
	LastUpdatedAt time.Time    `json:"last_updated_at"`
}

// --- Events ---

// Event carries facts reported from a device or provider.
// Payload shape is defined by the entity package for the given Domain.
type Event struct {
	ID            string          `json:"id"`
	PluginID      string          `json:"plugin_id"`
	DeviceID      string          `json:"device_id"`
	EntityID      string          `json:"entity_id"`
	EntityType    string          `json:"entity_type"`
	CorrelationID string          `json:"correlation_id,omitempty"`
	Payload       json.RawMessage `json:"payload"`
	CreatedAt     time.Time       `json:"created_at"`
}

// InboundEvent is emitted by plugin code after provider-specific last-mile work.
// The runner derives EntityType from the stored entity's Domain.
type InboundEvent struct {
	DeviceID      string          `json:"device_id"`
	EntityID      string          `json:"entity_id"`
	CorrelationID string          `json:"correlation_id,omitempty"`
	Payload       json.RawMessage `json:"payload"`
}

// EntityEventEnvelope is the normalized event published on the NATS bus
// so the gateway can fan out updates (e.g. to virtual entities).
type EntityEventEnvelope struct {
	EventID       string          `json:"event_id"`
	PluginID      string          `json:"plugin_id"`
	DeviceID      string          `json:"device_id"`
	EntityID      string          `json:"entity_id"`
	EntityType    string          `json:"entity_type"`
	CorrelationID string          `json:"correlation_id,omitempty"`
	Payload       json.RawMessage `json:"payload"`
	CreatedAt     time.Time       `json:"created_at"`
}

// --- Search ---

type SearchQuery struct {
	Pattern string `json:"pattern"`
}
