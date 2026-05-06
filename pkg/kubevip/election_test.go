package kubevip

import "testing"

func TestNeedsServiceElection(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		annotations map[string]string
		want        bool
	}{
		{
			name:   "global services election enabled, no annotations",
			config: &Config{EnableServicesElection: true},
			want:   true,
		},
		{
			name:        "global services election enabled, egress true",
			config:      &Config{EnableServicesElection: true},
			annotations: map[string]string{Egress: "true"},
			want:        true,
		},
		{
			name:   "no election flags, no annotations",
			config: &Config{},
			want:   false,
		},
		{
			name:        "no election flags, non-egress service",
			config:      &Config{},
			annotations: map[string]string{"some-key": "some-value"},
			want:        false,
		},
		{
			name:        "no election flags, egress enabled",
			config:      &Config{},
			annotations: map[string]string{Egress: "true"},
			want:        true,
		},
		{
			name:        "global leader election, egress enabled — no per-service election needed",
			config:      &Config{KubernetesLeaderElection: KubernetesLeaderElection{EnableLeaderElection: true}},
			annotations: map[string]string{Egress: "true"},
			want:        false,
		},
		{
			name:        "both elections enabled, egress enabled",
			config:      &Config{EnableServicesElection: true, KubernetesLeaderElection: KubernetesLeaderElection{EnableLeaderElection: true}},
			annotations: map[string]string{Egress: "true"},
			want:        true,
		},
		{
			name:        "egress annotation is false",
			config:      &Config{},
			annotations: map[string]string{Egress: "false"},
			want:        false,
		},
		{
			name:        "nil annotations",
			config:      &Config{},
			annotations: nil,
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NeedsServiceElection(tt.config, tt.annotations)
			if got != tt.want {
				t.Errorf("NeedsServiceElection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsElectionEnabled(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		annotations map[string]string
		want        bool
	}{
		{
			name:   "no elections, no egress",
			config: &Config{},
			want:   false,
		},
		{
			name:   "global leader election only",
			config: &Config{KubernetesLeaderElection: KubernetesLeaderElection{EnableLeaderElection: true}},
			want:   true,
		},
		{
			name:   "services election only",
			config: &Config{EnableServicesElection: true},
			want:   true,
		},
		{
			name:        "no global flags, egress service",
			config:      &Config{},
			annotations: map[string]string{Egress: "true"},
			want:        true,
		},
		{
			name:        "global leader election, egress service",
			config:      &Config{KubernetesLeaderElection: KubernetesLeaderElection{EnableLeaderElection: true}},
			annotations: map[string]string{Egress: "true"},
			want:        true,
		},
		{
			name:        "no flags, non-egress service",
			config:      &Config{},
			annotations: map[string]string{"other": "annotation"},
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsElectionEnabled(tt.config, tt.annotations)
			if got != tt.want {
				t.Errorf("IsElectionEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}
