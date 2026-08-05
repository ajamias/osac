/*
Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package inventory

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/v1/nodes"
)

func TestValidateMatchExpressions(t *testing.T) {
	tests := []struct {
		name             string
		matchExpressions map[string]string
		wantError        bool
		errorContains    string
	}{
		{
			name:             "valid simple labels",
			matchExpressions: map[string]string{"key1": "value1", "key2": "value2"},
			wantError:        false,
		},
		{
			name:             "empty map is valid",
			matchExpressions: map[string]string{},
			wantError:        false,
		},
		{
			name:             "nil map is valid",
			matchExpressions: nil,
			wantError:        false,
		},
		{
			name:             "empty key is invalid",
			matchExpressions: map[string]string{"": "value1"},
			wantError:        true,
			errorContains:    "empty label key",
		},
		{
			name:             "empty value is valid",
			matchExpressions: map[string]string{"key1": ""},
			wantError:        false,
		},
		{
			name:             "key with special characters",
			matchExpressions: map[string]string{"osac.openshift.io/hardware-profile": "gpu-large"},
			wantError:        false,
		},
		{
			name:             "key with spaces is invalid",
			matchExpressions: map[string]string{"key with spaces": "value1"},
			wantError:        true,
			errorContains:    "invalid label key",
		},
		{
			name:             "reserved bareMetalInstanceId key is invalid",
			matchExpressions: map[string]string{"bareMetalInstanceId": "some-value"},
			wantError:        true,
			errorContains:    "reserved label key",
		},
		{
			name:             "reserved managedBy key is invalid",
			matchExpressions: map[string]string{"managedBy": "some-value"},
			wantError:        true,
			errorContains:    "reserved label key",
		},
		{
			name:             "multiple valid labels",
			matchExpressions: map[string]string{"hostType": "gpu-large", "environment": "production"},
			wantError:        false,
		},
		{
			name:             "multiple labels with one invalid key",
			matchExpressions: map[string]string{"validKey": "validValue", "": "invalidEmptyKey"},
			wantError:        true,
			errorContains:    "empty label key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMatchExpressions(tt.matchExpressions)
			if tt.wantError {
				if err == nil {
					t.Errorf("validateMatchExpressions() expected error but got nil")
				} else if tt.errorContains != "" && !containsString(err.Error(), tt.errorContains) {
					t.Errorf("validateMatchExpressions() error = %v, want error containing %q", err, tt.errorContains)
				}
			} else {
				if err != nil {
					t.Errorf("validateMatchExpressions() unexpected error = %v", err)
				}
			}
		})
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 &&
			stringContainsSubstring(s, substr)))
}

func stringContainsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestMatchesLabels(t *testing.T) {
	tests := []struct {
		name             string
		node             *nodes.Node
		matchExpressions map[string]string
		wantMatch        bool
		wantError        bool
		errorContains    string
	}{
		{
			name: "node with matching labels",
			node: &nodes.Node{
				UUID: "test-node-1",
				Extra: map[string]interface{}{
					"osac_labels": map[string]interface{}{
						"environment": "production",
						"hardware":    "gpu-large",
					},
				},
			},
			matchExpressions: map[string]string{"environment": "production", "hardware": "gpu-large"},
			wantMatch:        true,
			wantError:        false,
		},
		{
			name: "node with partial matching labels",
			node: &nodes.Node{
				UUID: "test-node-2",
				Extra: map[string]interface{}{
					"osac_labels": map[string]interface{}{
						"environment": "production",
						"hardware":    "gpu-small",
					},
				},
			},
			matchExpressions: map[string]string{"environment": "production", "hardware": "gpu-large"},
			wantMatch:        false,
			wantError:        false,
		},
		{
			name: "node missing required labels",
			node: &nodes.Node{
				UUID: "test-node-3",
				Extra: map[string]interface{}{
					"osac_labels": map[string]interface{}{
						"environment": "staging",
					},
				},
			},
			matchExpressions: map[string]string{"environment": "production", "hardware": "gpu-large"},
			wantMatch:        false,
			wantError:        false,
		},
		{
			name: "node with no osac_labels",
			node: &nodes.Node{
				UUID:  "test-node-4",
				Extra: map[string]interface{}{},
			},
			matchExpressions: map[string]string{"environment": "production"},
			wantMatch:        false,
			wantError:        false,
		},
		{
			name: "node with nil Extra",
			node: &nodes.Node{
				UUID:  "test-node-5",
				Extra: nil,
			},
			matchExpressions: map[string]string{"environment": "production"},
			wantMatch:        false,
			wantError:        false,
		},
		{
			name: "empty match expressions matches any node",
			node: &nodes.Node{
				UUID: "test-node-6",
				Extra: map[string]interface{}{
					"osac_labels": map[string]interface{}{
						"environment": "production",
					},
				},
			},
			matchExpressions: map[string]string{},
			wantMatch:        true,
			wantError:        false,
		},
		{
			name: "nil match expressions matches any node",
			node: &nodes.Node{
				UUID: "test-node-7",
				Extra: map[string]interface{}{
					"osac_labels": map[string]interface{}{
						"environment": "production",
					},
				},
			},
			matchExpressions: nil,
			wantMatch:        true,
			wantError:        false,
		},
		{
			name: "node with malformed osac_labels type",
			node: &nodes.Node{
				UUID: "test-node-8",
				Extra: map[string]interface{}{
					"osac_labels": "not-a-map",
				},
			},
			matchExpressions: map[string]string{"environment": "production"},
			wantMatch:        false,
			wantError:        false,
		},
		{
			name: "exact string matching",
			node: &nodes.Node{
				UUID: "test-node-9",
				Extra: map[string]interface{}{
					"osac_labels": map[string]interface{}{
						"version": "1.2.3",
					},
				},
			},
			matchExpressions: map[string]string{"version": "1.2"},
			wantMatch:        false,
			wantError:        false,
		},
		{
			name: "empty value matching",
			node: &nodes.Node{
				UUID: "test-node-10",
				Extra: map[string]interface{}{
					"osac_labels": map[string]interface{}{
						"tag": "",
					},
				},
			},
			matchExpressions: map[string]string{"tag": ""},
			wantMatch:        true,
			wantError:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := matchesLabels(tt.node, tt.matchExpressions)

			if tt.wantError {
				if err == nil {
					t.Errorf("matchesLabels() expected error but got nil")
				} else if tt.errorContains != "" && !containsString(err.Error(), tt.errorContains) {
					t.Errorf("matchesLabels() error = %v, want error containing %q", err, tt.errorContains)
				}
				return
			}

			if err != nil {
				t.Errorf("matchesLabels() unexpected error = %v", err)
				return
			}

			if match != tt.wantMatch {
				t.Errorf("matchesLabels() = %v, want %v", match, tt.wantMatch)
			}
		})
	}
}

func TestOpenStackClientFindFreeHostValidation(t *testing.T) {
	tests := []struct {
		name             string
		matchExpressions map[string]string
		wantError        bool
		errorContains    string
	}{
		{
			name:             "invalid empty key should error before querying Ironic",
			matchExpressions: map[string]string{"": "value1"},
			wantError:        true,
			errorContains:    "empty label key",
		},
		{
			name:             "reserved bareMetalInstanceId key should error before querying Ironic",
			matchExpressions: map[string]string{"bareMetalInstanceId": "some-value"},
			wantError:        true,
			errorContains:    "reserved label key",
		},
		{
			name:             "reserved managedBy key should error before querying Ironic",
			matchExpressions: map[string]string{"managedBy": "some-value"},
			wantError:        true,
			errorContains:    "reserved label key",
		},
		{
			name:             "key with spaces should error before querying Ironic",
			matchExpressions: map[string]string{"key with spaces": "value1"},
			wantError:        true,
			errorContains:    "invalid label key",
		},
	}

	// Create a mock client to test validation without hitting actual API calls
	client := &OpenStackClient{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.findFreeHost(context.Background(), tt.matchExpressions)

			if tt.wantError {
				if err == nil {
					t.Errorf("findFreeHost() expected validation error but got nil")
				} else if tt.errorContains != "" && !containsString(err.Error(), tt.errorContains) {
					t.Errorf("findFreeHost() error = %v, want error containing %q", err, tt.errorContains)
				}
			}
		})
	}
}
