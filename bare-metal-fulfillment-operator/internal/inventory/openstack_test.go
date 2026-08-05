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
	"testing"
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
