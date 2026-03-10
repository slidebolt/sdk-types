package types

import (
	"encoding/json"
	"time"
)

// decodeLabels unmarshals a raw label map in canonical array format:
// "room": ["kitchen"].
func decodeLabels(raw map[string]json.RawMessage) map[string][]string {
	if len(raw) == 0 {
		return nil
	}
	out := make(map[string][]string, len(raw))
	for k, v := range raw {
		var ss []string
		if json.Unmarshal(v, &ss) == nil {
			out[k] = ss
		}
	}
	return out
}

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
	ID      string             `json:"id"`
	Name    string             `json:"name"`
	Version string             `json:"version"`
	Schemas []DomainDescriptor `json:"schemas,omitempty"`
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

// Device represents a physical or virtual piece of hardware.
// Ownership Contract:
// - ID: Owned by the Plugin. Must be deterministic and stable.
// - SourceID: Owned by Hardware. The raw technical ID (e.g. MAC).
// - SourceName: Owned by Hardware. The name the device calls itself.
// - LocalName: Owned by User. Only modified via user API actions. Hardware discovery must never overwrite this.
// - Labels: Shared. User modifications take precedence over hardware defaults.
// Protocol-specific raw data lives in the plugin's RawStore, not here.
type Device struct {
	ID          string              `json:"id"`
	SourceID    string              `json:"source_id"`
	SourceName  string              `json:"source_name,omitempty"`
	LocalName   string              `json:"local_name"`
	Labels      map[string][]string `json:"labels,omitempty"`
	EntityQuery string              `json:"entity_query,omitempty"`
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

func (d *Device) UnmarshalJSON(data []byte) error {
	var w struct {
		ID          string                     `json:"id"`
		SourceID    string                     `json:"source_id"`
		SourceName  string                     `json:"source_name"`
		LocalName   string                     `json:"local_name"`
		Labels      map[string]json.RawMessage `json:"labels,omitempty"`
		EntityQuery string                     `json:"entity_query,omitempty"`
	}
	if err := json.Unmarshal(data, &w); err != nil {
		return err
	}
	d.ID, d.SourceID, d.SourceName, d.LocalName = w.ID, w.SourceID, w.SourceName, w.LocalName
	d.Labels = decodeLabels(w.Labels)
	d.EntityQuery = w.EntityQuery
	return nil
}

// EntitySnapshot captures the Effective state of an entity at a point in time.
// ID is the stable UUID used for data transmission.
// Name and Labels are for human-readable discovery.
type EntitySnapshot struct {
	ID        string              `json:"id"`
	Name      string              `json:"name"`
	State     json.RawMessage     `json:"state"`
	Labels    map[string][]string `json:"labels,omitempty"`
	CreatedAt time.Time           `json:"created_at"`
}

type Entity struct {
	ID         string                    `json:"id"`
	SourceID   string                    `json:"source_id"`
	SourceName string                    `json:"source_name,omitempty"`
	DeviceID   string                    `json:"device_id"`
	Domain     string                    `json:"domain"`
	LocalName  string                    `json:"local_name"`
	Actions    []string                  `json:"actions,omitempty"`
	Data       EntityData                `json:"data"`
	Labels     map[string][]string       `json:"labels,omitempty"`
	Snapshots  map[string]EntitySnapshot `json:"snapshots,omitempty"`
}

func (e *Entity) UnmarshalJSON(data []byte) error {
	var w struct {
		ID         string                     `json:"id"`
		SourceID   string                     `json:"source_id"`
		SourceName string                     `json:"source_name,omitempty"`
		DeviceID   string                     `json:"device_id"`
		Domain     string                     `json:"domain"`
		LocalName  string                     `json:"local_name"`
		Actions    []string                   `json:"actions,omitempty"`
		Data       EntityData                 `json:"data"`
		Labels     map[string]json.RawMessage `json:"labels,omitempty"`
		Snapshots  map[string]EntitySnapshot  `json:"snapshots,omitempty"`
	}
	if err := json.Unmarshal(data, &w); err != nil {
		return err
	}
	e.ID, e.SourceID, e.SourceName, e.DeviceID, e.Domain, e.LocalName = w.ID, w.SourceID, w.SourceName, w.DeviceID, w.Domain, w.LocalName
	e.Actions, e.Data = w.Actions, w.Data
	e.Labels = decodeLabels(w.Labels)
	e.Snapshots = w.Snapshots
	return nil
}

type EntityData struct {
	Desired       json.RawMessage `json:"desired,omitempty"`
	Reported      json.RawMessage `json:"reported,omitempty"`
	Effective     json.RawMessage `json:"effective,omitempty"`
	SyncStatus    SyncStatus      `json:"sync_status,omitempty"`
	LastCommandID string          `json:"last_command_id,omitempty"`
	LastEventID   string          `json:"last_event_id,omitempty"`
	UpdatedAt     time.Time       `json:"updated_at,omitempty"`
}

// SyncStatus is the canonical state sync enum persisted in EntityData.
type SyncStatus string

const (
	SyncStatusEmpty   SyncStatus = ""
	SyncStatusSynced  SyncStatus = "synced"
	SyncStatusPending SyncStatus = "pending"
	SyncStatusFailed  SyncStatus = "failed"
)

// NormalizeSyncStatus converts any non-canonical value to a canonical one.
func NormalizeSyncStatus(v SyncStatus) SyncStatus {
	switch v {
	case SyncStatusSynced, SyncStatusPending, SyncStatusFailed, SyncStatusEmpty:
		return v
	default:
		return SyncStatusEmpty
	}
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

// --- Batch ---

type BatchDeviceRef struct {
	PluginID string `json:"plugin_id"`
	DeviceID string `json:"device_id"`
}

type BatchEntityRef struct {
	PluginID string `json:"plugin_id"`
	DeviceID string `json:"device_id"`
	EntityID string `json:"entity_id"`
}

type BatchDeviceItem struct {
	PluginID string `json:"plugin_id"`
	Device   Device `json:"device"`
}

type BatchEntityItem struct {
	PluginID string `json:"plugin_id"`
	DeviceID string `json:"device_id"`
	Entity   Entity `json:"entity"`
}

type BatchCommandItem struct {
	PluginID string          `json:"plugin_id"`
	DeviceID string          `json:"device_id"`
	EntityID string          `json:"entity_id"`
	Payload  json.RawMessage `json:"payload"`
}

type BatchResult struct {
	PluginID string          `json:"plugin_id,omitempty"`
	DeviceID string          `json:"device_id,omitempty"`
	EntityID string          `json:"entity_id,omitempty"`
	OK       bool            `json:"ok"`
	Error    string          `json:"error,omitempty"`
	Data     json.RawMessage `json:"data,omitempty"`
}

type BatchCommandResult struct {
	PluginID  string       `json:"plugin_id,omitempty"`
	DeviceID  string       `json:"device_id,omitempty"`
	EntityID  string       `json:"entity_id,omitempty"`
	CommandID string       `json:"command_id,omitempty"`
	State     CommandState `json:"state,omitempty"`
	OK        bool         `json:"ok"`
	Error     string       `json:"error,omitempty"`
}

// --- Schema ---

type FieldDescriptor struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
	Min      *int   `json:"min,omitempty"`
	Max      *int   `json:"max,omitempty"`
}

type ActionDescriptor struct {
	Action string            `json:"action"`
	Fields []FieldDescriptor `json:"fields,omitempty"`
}

type DomainDescriptor struct {
	Domain   string             `json:"domain"`
	Commands []ActionDescriptor `json:"commands"`
	Events   []ActionDescriptor `json:"events"`
}

// --- Search ---

type SearchQuery struct {
	Pattern  string              `json:"pattern"`
	Labels   map[string][]string `json:"labels,omitempty"`
	PluginID string              `json:"plugin_id,omitempty"`
	DeviceID string              `json:"device_id,omitempty"`
	EntityID string              `json:"entity_id,omitempty"`
	Domain   string              `json:"domain,omitempty"`
	Limit    int                 `json:"limit,omitempty"`
}

type SearchPluginsResponse struct {
	PluginID string     `json:"plugin_id"`
	Matches  []Manifest `json:"matches"`
}

type SearchDevicesResponse struct {
	PluginID string   `json:"plugin_id"`
	Matches  []Device `json:"matches"`
}

type SearchEntitiesResponse struct {
	PluginID string   `json:"plugin_id"`
	Matches  []Entity `json:"matches"`
}

// --- Core Entities ---

// CoreDeviceID returns the management device ID for a plugin.
// By convention each plugin exposes a management device using its own pluginID as the device ID.
func CoreDeviceID(pluginID string) string { return pluginID }

// CoreEntities returns the standard health entity for the core management device.
func CoreEntities(pluginID string) []Entity {
	deviceID := CoreDeviceID(pluginID)
	return []Entity{
		{ID: "health", SourceID: "health", SourceName: "Health", DeviceID: deviceID, Domain: "plugin.health", LocalName: "Health"},
	}
}

// CoreDomains returns the DomainDescriptors for all core entity domains.
// These should be included in every plugin's Manifest Schemas.
func CoreDomains() []DomainDescriptor {
	return []DomainDescriptor{
		{
			Domain:   "stream",
			Commands: []ActionDescriptor{},
			Events: []ActionDescriptor{
				{
					Action: "updated",
					Fields: []FieldDescriptor{
						{Name: "url", Type: "string", Required: true},
						{Name: "format", Type: "string", Required: false},
						{Name: "kind", Type: "string", Required: false},
						{Name: "online", Type: "bool", Required: false},
					},
				},
			},
		},
		{
			Domain:   "image",
			Commands: []ActionDescriptor{},
			Events: []ActionDescriptor{
				{
					Action: "updated",
					Fields: []FieldDescriptor{
						{Name: "url", Type: "string", Required: true},
						{Name: "format", Type: "string", Required: false},
					},
				},
			},
		},
		{
			Domain:   "plugin.health",
			Commands: []ActionDescriptor{},
			Events: []ActionDescriptor{
				{
					Action: "status",
					Fields: []FieldDescriptor{
						{Name: "status", Type: "string", Required: true},
						{Name: "message", Type: "string", Required: false},
						{Name: "ts", Type: "string", Required: true},
					},
				},
			},
		},
	}
}
