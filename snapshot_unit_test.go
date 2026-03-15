package types

import (
	"encoding/json"
	"testing"
	"time"
)

func TestEntitySnapshotRoundTrip(t *testing.T) {
	snap := EntitySnapshot{
		ID:        "550e8400-e29b-41d4-a716-446655440000",
		Name:      "MovieTime",
		State:     json.RawMessage(`{"power":true,"brightness":40,"rgb":[255,100,0]}`),
		Labels:    map[string][]string{"room": {"living-room"}},
		CreatedAt: time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC),
	}

	b, err := json.Marshal(snap)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var got EntitySnapshot
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if got.ID != snap.ID {
		t.Errorf("ID: got %q want %q", got.ID, snap.ID)
	}
	if got.Name != snap.Name {
		t.Errorf("Name: got %q want %q", got.Name, snap.Name)
	}
	if string(got.State) != string(snap.State) {
		t.Errorf("State: got %s want %s", got.State, snap.State)
	}
	if len(got.Labels["room"]) != 1 || got.Labels["room"][0] != "living-room" {
		t.Errorf("Labels: got %v", got.Labels)
	}
	if !got.CreatedAt.Equal(snap.CreatedAt) {
		t.Errorf("CreatedAt: got %v want %v", got.CreatedAt, snap.CreatedAt)
	}
}

func TestEntityWithSnapshotsRoundTrip(t *testing.T) {
	snap := EntitySnapshot{
		ID:        "snap-uuid-001",
		Name:      "MovieTime",
		State:     json.RawMessage(`{"power":true,"brightness":40}`),
		CreatedAt: time.Now().UTC().Truncate(time.Second),
	}

	ent := Entity{
		ID:       "ent-001",
		DeviceID: "dev-001",
		Domain:   "light",
		Snapshots: map[string]EntitySnapshot{
			snap.ID: snap,
		},
	}

	b, err := json.Marshal(ent)
	if err != nil {
		t.Fatalf("marshal entity failed: %v", err)
	}

	var got Entity
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal entity failed: %v", err)
	}

	if len(got.Snapshots) != 1 {
		t.Fatalf("expected 1 snapshot, got %d", len(got.Snapshots))
	}
	gotSnap, ok := got.Snapshots[snap.ID]
	if !ok {
		t.Fatalf("snapshot %q missing from decoded entity", snap.ID)
	}
	if gotSnap.Name != "MovieTime" {
		t.Errorf("snapshot Name: got %q want %q", gotSnap.Name, "MovieTime")
	}
	if string(gotSnap.State) != string(snap.State) {
		t.Errorf("snapshot State: got %s want %s", gotSnap.State, snap.State)
	}
}

func TestEntityWithNoSnapshotsOmitted(t *testing.T) {
	ent := Entity{ID: "ent-001", DeviceID: "dev-001", Domain: "switch"}

	b, err := json.Marshal(ent)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		t.Fatalf("unmarshal to raw map failed: %v", err)
	}
	if _, ok := raw["snapshots"]; ok {
		t.Error("snapshots key should be omitted when nil")
	}
}

func TestStateToCommandsRegistryRoundTrip(t *testing.T) {
	const domain = "test.domain"
	registered := false

	RegisterStateToCommands(domain, func(stateJSON json.RawMessage) ([]json.RawMessage, error) {
		registered = true
		return []json.RawMessage{stateJSON}, nil
	})

	payloads, err := StateToCommands(domain, json.RawMessage(`{"power":true}`))
	if err != nil {
		t.Fatalf("StateToCommands failed: %v", err)
	}
	if !registered {
		t.Error("registered function was not called")
	}
	if len(payloads) != 1 {
		t.Fatalf("expected 1 payload, got %d", len(payloads))
	}
}

func TestStateToCommandsUnknownDomainReturnsError(t *testing.T) {
	_, err := StateToCommands("domain.that.does.not.exist", json.RawMessage(`{}`))
	if err == nil {
		t.Fatal("expected error for unregistered domain, got nil")
	}
}
