package kubevip

// NeedsServiceElection returns true if per-service leader election should be
// used for a service with the given annotations. This is the case when:
//   - EnableServicesElection is globally enabled, OR
//   - The service has egress enabled and there is no global leader election
//     to serialize service processing
func NeedsServiceElection(config *Config, annotations map[string]string) bool {
	if config.EnableServicesElection {
		return true
	}
	if annotations[Egress] == "true" && !config.EnableLeaderElection {
		return true
	}
	return false
}

// IsElectionEnabled returns true if any form of leader election (global or
// per-service) applies for a service with the given annotations.
func IsElectionEnabled(config *Config, annotations map[string]string) bool {
	return config.EnableLeaderElection || NeedsServiceElection(config, annotations)
}
