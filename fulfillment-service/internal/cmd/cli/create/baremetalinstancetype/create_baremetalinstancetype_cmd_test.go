/*
Copyright (c) 2026 Red Hat Inc.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the
License. You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific
language governing permissions and limitations under the License.
*/

package baremetalinstancetype

import (
	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"
)

var _ = Describe("create baremetalinstancetype command", func() {

	Context("command structure", func() {
		It("should create command without error", func() {
			cmd := Cmd()
			Expect(cmd).NotTo(BeNil())
			Expect(cmd.Use).To(Equal("baremetalinstancetype"))
		})

		It("should have the correct aliases", func() {
			cmd := Cmd()
			Expect(cmd.Aliases).To(ContainElement("osac.private.v1.BareMetalInstanceType"))
		})

		It("should accept no arguments", func() {
			cmd := Cmd()
			Expect(cmd.Args).NotTo(BeNil())
		})
	})

	Context("flag registration", func() {
		It("should register --name flag", func() {
			cmd := Cmd()
			flag := cmd.Flags().Lookup("name")
			Expect(flag).NotTo(BeNil())
			Expect(flag.Usage).To(ContainSubstring("NAME"))
		})

		It("should register --cpu-cores flag", func() {
			cmd := Cmd()
			flag := cmd.Flags().Lookup("cpu-cores")
			Expect(flag).NotTo(BeNil())
			Expect(flag.Usage).To(ContainSubstring("CORES"))
		})

		It("should register --cpu-architecture flag", func() {
			cmd := Cmd()
			flag := cmd.Flags().Lookup("cpu-architecture")
			Expect(flag).NotTo(BeNil())
			Expect(flag.Usage).To(ContainSubstring("ARCHITECTURE"))
		})

		It("should register --memory-total-gb flag", func() {
			cmd := Cmd()
			flag := cmd.Flags().Lookup("memory-total-gb")
			Expect(flag).NotTo(BeNil())
			Expect(flag.Usage).To(ContainSubstring("MEMORY"))
		})

		It("should register optional --description flag", func() {
			cmd := Cmd()
			flag := cmd.Flags().Lookup("description")
			Expect(flag).NotTo(BeNil())
			Expect(flag.Usage).To(ContainSubstring("DESCRIPTION"))
		})

		It("should register optional --cpu-model flag", func() {
			cmd := Cmd()
			flag := cmd.Flags().Lookup("cpu-model")
			Expect(flag).NotTo(BeNil())
			Expect(flag.Usage).To(ContainSubstring("MODEL"))
		})

		It("should register optional --cpu-threads-per-core flag", func() {
			cmd := Cmd()
			flag := cmd.Flags().Lookup("cpu-threads-per-core")
			Expect(flag).NotTo(BeNil())
			Expect(flag.Usage).To(ContainSubstring("THREADS"))
		})

		It("should register optional --memory-type flag", func() {
			cmd := Cmd()
			flag := cmd.Flags().Lookup("memory-type")
			Expect(flag).NotTo(BeNil())
			Expect(flag.Usage).To(ContainSubstring("TYPE"))
		})
	})

	Context("help text", func() {
		It("should have proper short help", func() {
			cmd := Cmd()
			Expect(cmd.Short).To(Equal("Create a bare metal instance type."))
		})

		It("should have proper long help with examples", func() {
			cmd := Cmd()
			Expect(cmd.Long).To(ContainSubstring("Create a bare metal instance type."))
			Expect(cmd.Long).To(ContainSubstring("create baremetalinstancetype"))
			Expect(cmd.Long).To(ContainSubstring("--name gpu-large"))
			Expect(cmd.Long).To(ContainSubstring("--cpu-cores"))
			Expect(cmd.Long).To(ContainSubstring("--cpu-architecture"))
			Expect(cmd.Long).To(ContainSubstring("--memory-total-gb"))
		})
	})

})
