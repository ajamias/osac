/*
Copyright (c) 2026 Red Hat Inc.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the
License. You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific
language governing permissions and limitations under the License.
*/

package auth

import (
	"context"
	"log/slog"

	"github.com/osac-project/fulfillment-service/internal/collections"
)

// SharedTenancyLogicBuilder contains the data and logic needed to create shared tenancy logic.
type SharedTenancyLogicBuilder struct {
	logger *slog.Logger
}

// SharedTenancyLogic is a tenancy logic implementation that assigns the shared tenant to all objects
// and makes the shared tenant visible to all users. This is used for platform-scoped resources like
// BareMetalInstanceType that should be globally visible across tenants.
type SharedTenancyLogic struct {
	logger *slog.Logger
}

// NewSharedTenancyLogic creates a new builder for shared tenancy logic.
func NewSharedTenancyLogic() *SharedTenancyLogicBuilder {
	return &SharedTenancyLogicBuilder{}
}

// SetLogger sets the logger that will be used by the tenancy logic.
func (b *SharedTenancyLogicBuilder) SetLogger(value *slog.Logger) *SharedTenancyLogicBuilder {
	b.logger = value
	return b
}

// Build creates the shared tenancy logic.
func (b *SharedTenancyLogicBuilder) Build() (result *SharedTenancyLogic, err error) {
	result = &SharedTenancyLogic{
		logger: b.logger,
	}
	return
}

// DetermineAssignableTenants returns a set containing only the shared tenant, regardless of the user's identity.
// Platform-scoped resources like BareMetalInstanceType can only be assigned to the shared tenant.
func (p *SharedTenancyLogic) DetermineAssignableTenants(_ context.Context) (result collections.Set[string], err error) {
	result = SharedTenants
	return
}

// DetermineDefaultTenant returns the shared tenant, regardless of the user's identity.
// All platform-scoped resources default to the shared tenant for global visibility.
func (p *SharedTenancyLogic) DetermineDefaultTenant(_ context.Context) (result string, err error) {
	result = SharedTenant
	return
}

// DetermineVisibleTenants returns a set containing the shared tenant plus any tenants the user can access.
// This ensures shared resources are visible to all authenticated users while preserving tenant isolation.
func (p *SharedTenancyLogic) DetermineVisibleTenants(ctx context.Context) (result collections.Set[string], err error) {
	subject := SubjectFromContext(ctx)
	result = subject.Tenants
	if result.Finite() {
		result = SharedTenants.Union(result)
	}
	return
}
