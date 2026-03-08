package types

import (
	"encoding/json"
	"fmt"
	"sort"
)

var domainRegistry = map[string]DomainDescriptor{}

func RegisterDomain(d DomainDescriptor) {
	domainRegistry[d.Domain] = d
}

func GetDomainDescriptor(domain string) (DomainDescriptor, bool) {
	d, ok := domainRegistry[domain]
	return d, ok
}

func AllDomainDescriptors() []DomainDescriptor {
	out := make([]DomainDescriptor, 0, len(domainRegistry))
	for _, d := range domainRegistry {
		out = append(out, d)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Domain < out[j].Domain })
	return out
}

// StateToCommandsFunc converts a domain-specific state blob into a slice of
// raw JSON command payloads. Registered by entity packages at init time.
type StateToCommandsFunc func(stateJSON json.RawMessage) ([]json.RawMessage, error)

var stateToCommandsRegistry = map[string]StateToCommandsFunc{}

// RegisterStateToCommands registers a state→commands converter for a domain.
// Called from entity package init() functions alongside RegisterDomain.
func RegisterStateToCommands(domain string, fn StateToCommandsFunc) {
	stateToCommandsRegistry[domain] = fn
}

// StateToCommands converts a state blob to command payloads using the
// registered converter for the given domain. Returns an error if no
// converter is registered for the domain.
func StateToCommands(domain string, stateJSON json.RawMessage) ([]json.RawMessage, error) {
	fn, ok := stateToCommandsRegistry[domain]
	if !ok {
		return nil, fmt.Errorf("no state-to-commands handler registered for domain %q", domain)
	}
	return fn(stateJSON)
}
