/*
Copyright (c) 2026 Red Hat Inc.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the
License. You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific
language governing permissions and limitations under the License.
*/

package auth_test

import (
	"context"
	"log/slog"

	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"

	"github.com/osac-project/fulfillment-service/internal/auth"
	"github.com/osac-project/fulfillment-service/internal/collections"
)

var _ = Describe("SharedTenancyLogic", func() {
	var (
		ctx    context.Context
		logger *slog.Logger
		logic  *auth.SharedTenancyLogic
		err    error
	)

	BeforeEach(func() {
		ctx = context.Background()
		logger = slog.Default()

		logic, err = auth.NewSharedTenancyLogic().
			SetLogger(logger).
			Build()
		Expect(err).ToNot(HaveOccurred())
		Expect(logic).ToNot(BeNil())
	})

	Describe("DetermineAssignableTenants", func() {
		It("Always returns shared tenants regardless of user context", func() {
			// Test with no subject in context
			result, err := logic.DetermineAssignableTenants(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(auth.SharedTenants))

			// Test with regular user subject
			userSubject := &auth.Subject{
				User:    "test-user",
				Tenants: collections.NewSet("tenant-a", "tenant-b"),
			}
			userCtx := auth.ContextWithSubject(ctx, userSubject)

			result, err = logic.DetermineAssignableTenants(userCtx)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(auth.SharedTenants))
			Expect(result.Contains("shared")).To(BeTrue())
			Expect(result.Contains("tenant-a")).To(BeFalse())
			Expect(result.Contains("tenant-b")).To(BeFalse())

			// Test with admin user (universal tenant set)
			adminSubject := &auth.Subject{
				User:    "admin-user",
				Tenants: auth.AllTenants,
			}
			adminCtx := auth.ContextWithSubject(ctx, adminSubject)

			result, err = logic.DetermineAssignableTenants(adminCtx)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(auth.SharedTenants))
		})
	})

	Describe("DetermineDefaultTenant", func() {
		It("Always returns shared tenant regardless of user context", func() {
			// Test with no subject in context
			result, err := logic.DetermineDefaultTenant(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal("shared"))

			// Test with regular user subject
			userSubject := &auth.Subject{
				User:    "test-user",
				Tenants: collections.NewSet("tenant-a", "tenant-b"),
			}
			userCtx := auth.ContextWithSubject(ctx, userSubject)

			result, err = logic.DetermineDefaultTenant(userCtx)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal("shared"))

			// Test with admin user (universal tenant set)
			adminSubject := &auth.Subject{
				User:    "admin-user",
				Tenants: auth.AllTenants,
			}
			adminCtx := auth.ContextWithSubject(ctx, adminSubject)

			result, err = logic.DetermineDefaultTenant(adminCtx)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal("shared"))
		})
	})

	Describe("DetermineVisibleTenants", func() {
		It("Returns shared tenant plus user's tenants for regular users", func() {
			// Test with regular user subject
			userSubject := &auth.Subject{
				User:    "test-user",
				Tenants: collections.NewSet("tenant-a", "tenant-b"),
			}
			userCtx := auth.ContextWithSubject(ctx, userSubject)

			result, err := logic.DetermineVisibleTenants(userCtx)
			Expect(err).ToNot(HaveOccurred())

			expected := auth.SharedTenants.Union(collections.NewSet("tenant-a", "tenant-b"))
			Expect(result).To(Equal(expected))
			Expect(result.Contains("shared")).To(BeTrue())
			Expect(result.Contains("tenant-a")).To(BeTrue())
			Expect(result.Contains("tenant-b")).To(BeTrue())
		})

		It("Returns all tenants for admin users", func() {
			// Test with admin user (universal tenant set)
			adminSubject := &auth.Subject{
				User:    "admin-user",
				Tenants: auth.AllTenants,
			}
			adminCtx := auth.ContextWithSubject(ctx, adminSubject)

			result, err := logic.DetermineVisibleTenants(adminCtx)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(auth.AllTenants))
		})

		It("Returns only shared tenant for users with empty tenant set", func() {
			// Test with user that has empty tenant set
			emptyTenantSubject := &auth.Subject{
				User:    "empty-tenant-user",
				Tenants: collections.NewSet[string](),
			}
			emptyCtx := auth.ContextWithSubject(ctx, emptyTenantSubject)

			result, err := logic.DetermineVisibleTenants(emptyCtx)
			Expect(err).ToNot(HaveOccurred())

			expected := auth.SharedTenants.Union(collections.NewSet[string]())
			Expect(result).To(Equal(expected))
			Expect(result.Contains("shared")).To(BeTrue())
		})
	})
})
