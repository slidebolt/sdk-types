package types

import (
	"reflect"
	"testing"
)

// TestDTOPurity enforces that transport and registry types carry no behaviour.
// DTOs must remain plain data — logic belongs in the runner or entity packages.
func TestDTOPurity(t *testing.T) {
	allowedMethods := map[string]int{
		"Device": 1, // Name()
	}

	typesToTest := []interface{}{
		Request{},
		Response{},
		RPCError{},
		Manifest{},
		Device{},
	}

	for _, typ := range typesToTest {
		rt := reflect.TypeOf(typ)
		limit := allowedMethods[rt.Name()]
		if rt.NumMethod() > limit {
			t.Errorf("Type %s has %d methods (limit %d). DTOs must remain pure.", rt.Name(), rt.NumMethod(), limit)
		}
	}
}

// Custom Payload Types for Testing
type MyCustomRequestPayload struct {
	Value int `json:"value"`
}

func (MyCustomRequestPayload) CommandRequestPayloadKind() string { return "my_custom_request" }

type MyCustomResponsePayload struct {
	Result string `json:"result"`
}

func (MyCustomResponsePayload) CommandResponsePayloadKind() string { return "my_custom_response" }

type MyCustomEventPayload struct {
	EventData string `json:"event_data"`
}

func (MyCustomEventPayload) EventPayloadKind() string { return "my_custom_event" }

// TestCustomPayloadSupport proves that consumers can define their own concrete
// payload types and use them with the generic transport boundaries, removing
// the need for a built-in GenericPayload type.
func TestCustomPayloadSupport(t *testing.T) {
	// 1. Verify custom command requests
	req := CommandRequest[MyCustomRequestPayload]{
		CommandID: "cmd-1",
		Payload: MyCustomRequestPayload{
			Value: 42,
		},
	}
	if req.Payload.Value != 42 {
		t.Errorf("Expected payload value 42, got %d", req.Payload.Value)
	}
	if req.Payload.CommandRequestPayloadKind() != "my_custom_request" {
		t.Errorf("Expected kind 'my_custom_request', got %s", req.Payload.CommandRequestPayloadKind())
	}

	// 2. Verify custom command responses
	res := CommandResponse[MyCustomResponsePayload]{
		Payload: MyCustomResponsePayload{
			Result: "success",
		},
	}
	if res.Payload.Result != "success" {
		t.Errorf("Expected payload result 'success', got %s", res.Payload.Result)
	}
	if res.Payload.CommandResponsePayloadKind() != "my_custom_response" {
		t.Errorf("Expected kind 'my_custom_response', got %s", res.Payload.CommandResponsePayloadKind())
	}

	// 3. Verify custom inbound events
	inboundEvt := InboundEventTyped[MyCustomEventPayload]{
		Payload: MyCustomEventPayload{
			EventData: "sensor_triggered",
		},
	}
	if inboundEvt.Payload.EventData != "sensor_triggered" {
		t.Errorf("Expected event data 'sensor_triggered', got %s", inboundEvt.Payload.EventData)
	}
	if inboundEvt.Payload.EventPayloadKind() != "my_custom_event" {
		t.Errorf("Expected kind 'my_custom_event', got %s", inboundEvt.Payload.EventPayloadKind())
	}

	// 4. Verify custom internal events
	evt := EventTyped[MyCustomEventPayload]{
		Payload: MyCustomEventPayload{
			EventData: "state_changed",
		},
	}
	if evt.Payload.EventData != "state_changed" {
		t.Errorf("Expected event data 'state_changed', got %s", evt.Payload.EventData)
	}
	if evt.Payload.EventPayloadKind() != "my_custom_event" {
		t.Errorf("Expected kind 'my_custom_event', got %s", evt.Payload.EventPayloadKind())
	}
}
