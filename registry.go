package types

import "sort"

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
