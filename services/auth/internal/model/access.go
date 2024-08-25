package model

// EndpointAccessRule contains access permissions for endpoint
type EndpointAccessRule struct {
	Endpoint string
	Roles    map[int]struct{}
}

// HasRole checks if role has acces to endpoint
func (r *EndpointAccessRule) HasRole(role int) bool {
	_, ok := r.Roles[role]
	return ok
}
