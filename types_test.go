package types

import (
	"reflect"
	"testing"
)

// TestDTOPurity enforces that transport and registry types carry no behaviour.
// DTOs must remain plain data — logic belongs in the runner or entity packages.
func TestDTOPurity(t *testing.T) {
	allowedMethods := map[string]int{
		"Device": 2, // Name() and MarshalJSON()
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
