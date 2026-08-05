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
	"fmt"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"

	privatev1 "github.com/osac-project/osac/fulfillment-service/internal/api/osac/private/v1"
	"github.com/osac-project/osac/fulfillment-service/internal/config"
	"github.com/osac-project/osac/fulfillment-service/internal/terminal"
)

// Cmd creates the command to create a bare metal instance type.
func Cmd() *cobra.Command {
	runner := &runnerContext{}
	result := &cobra.Command{
		Use:                   "baremetalinstancetype",
		Aliases:               []string{string(proto.MessageName((*privatev1.BareMetalInstanceType)(nil)))},
		Short:                 shortHelp,
		Long:                  longHelp,
		DisableFlagsInUseLine: true,
		Args:                  cobra.NoArgs,
		RunE:                  runner.run,
	}
	flags := result.Flags()
	flags.StringVar(
		&runner.name,
		"name",
		"",
		nameFlagHelp,
	)
	flags.StringVar(
		&runner.description,
		"description",
		"",
		descriptionFlagHelp,
	)
	// CPU flags
	flags.Int32Var(
		&runner.cpuCores,
		"cpu-cores",
		0,
		cpuCoresFlagHelp,
	)
	flags.StringVar(
		&runner.cpuArchitecture,
		"cpu-architecture",
		"",
		cpuArchitectureFlagHelp,
	)
	flags.StringVar(
		&runner.cpuModel,
		"cpu-model",
		"",
		cpuModelFlagHelp,
	)
	flags.Int32Var(
		&runner.cpuThreadsPerCore,
		"cpu-threads-per-core",
		0,
		cpuThreadsPerCoreFlagHelp,
	)
	// Memory flags
	flags.Int64Var(
		&runner.memoryTotalGb,
		"memory-total-gb",
		0,
		memoryTotalGbFlagHelp,
	)
	flags.StringVar(
		&runner.memoryType,
		"memory-type",
		"",
		memoryTypeFlagHelp,
	)
	return result
}

type runnerContext struct {
	console           *terminal.Console
	name              string
	description       string
	cpuCores          int32
	cpuArchitecture   string
	cpuModel          string
	cpuThreadsPerCore int32
	memoryTotalGb     int64
	memoryType        string
}

func (c *runnerContext) run(cmd *cobra.Command, args []string) error {
	// Get the context:
	ctx := cmd.Context()

	// Get the console:
	c.console = terminal.ConsoleFromContext(ctx)

	// Get the configuration:
	cfg := config.SettingsFromContext(ctx)
	if !cfg.Armed() {
		return fmt.Errorf("there is no configuration, run the 'login' command")
	}

	// Check the required parameters:
	if c.name == "" {
		return fmt.Errorf("name is required")
	}
	if c.cpuCores <= 0 {
		return fmt.Errorf("cpu-cores must be greater than zero")
	}
	if c.cpuThreadsPerCore < 0 {
		return fmt.Errorf("cpu-threads-per-core must not be negative")
	}
	if c.cpuArchitecture == "" {
		return fmt.Errorf("cpu-architecture is required")
	}
	if c.memoryTotalGb <= 0 {
		return fmt.Errorf("memory-total-gb must be greater than zero")
	}

	// Create the gRPC connection from the configuration:
	conn, err := cfg.Connect(ctx, cmd.Flags())
	if err != nil {
		return fmt.Errorf("failed to create gRPC connection: %w", err)
	}
	defer conn.Close()

	// Create the client:
	client := privatev1.NewBareMetalInstanceTypesClient(conn)

	// Build CPU specification:
	cpuSpec := privatev1.BareMetalCPUSpec_builder{
		Cores:        c.cpuCores,
		Architecture: c.cpuArchitecture,
	}
	if c.cpuModel != "" {
		cpuSpec.Model = c.cpuModel
	}
	if c.cpuThreadsPerCore > 0 {
		cpuSpec.ThreadsPerCore = c.cpuThreadsPerCore
	}

	// Build memory specification:
	memorySpec := privatev1.BareMetalMemorySpec_builder{
		TotalGb: c.memoryTotalGb,
	}
	if c.memoryType != "" {
		memorySpec.Type = c.memoryType
	}

	// Build hardware specification:
	hardwareSpec := privatev1.BareMetalHardwareSpec_builder{
		Cpu:    cpuSpec.Build(),
		Memory: memorySpec.Build(),
	}.Build()

	// Prepare the bare metal instance type:
	bareMetalInstanceType := privatev1.BareMetalInstanceType_builder{
		Id: c.name,
		Metadata: privatev1.Metadata_builder{
			Name: c.name,
		}.Build(),
		Spec: privatev1.BareMetalInstanceTypeSpec_builder{
			Hardware:    hardwareSpec,
			Description: c.description,
		}.Build(),
	}.Build()

	// Create the bare metal instance type:
	response, err := client.Create(ctx, privatev1.BareMetalInstanceTypesCreateRequest_builder{
		Object: bareMetalInstanceType,
	}.Build())
	if err != nil {
		return fmt.Errorf("failed to create bare metal instance type: %w", err)
	}

	// Display the result:
	c.console.Infof(ctx, "Created bare metal instance type '%s'.\n", response.GetObject().GetId())

	return nil
}

const shortHelp = `Create a bare metal instance type.`

const longHelp = `
Create a bare metal instance type.

A bare metal instance type defines a pre-configured hardware bundle (CPU, memory, and other
specifications) that can be referenced by name when provisioning bare metal instances.
Bare metal instance types are managed by Cloud Provider Admins.

To create a bare metal instance type:

{{ bt 3 }}shell
{{ binary }} create baremetalinstancetype \
  --name gpu-large \
  --description 'Large GPU node for ML workloads' \
  --cpu-cores 32 \
  --cpu-architecture x86_64 \
  --cpu-model 'Intel Xeon Gold 6338' \
  --cpu-threads-per-core 2 \
  --memory-total-gb 512 \
  --memory-type DDR4
{{ bt 3 }}

Required fields: name, cpu-cores, cpu-architecture, memory-total-gb
`

const nameFlagHelp = `
_NAME_ - Name of the bare metal instance type. Must be a unique, human-readable identifier
(e.g., {{ bt }}gpu-large{{ bt }}).
`

const descriptionFlagHelp = `
_DESCRIPTION_ - Human friendly description of the bare metal instance type.
`

const cpuCoresFlagHelp = `
_CORES_ - Number of CPU cores. Must be greater than zero.
`

const cpuArchitectureFlagHelp = `
_ARCHITECTURE_ - CPU architecture (e.g., {{ bt }}x86_64{{ bt }}, {{ bt }}aarch64{{ bt }}). Required.
`

const cpuModelFlagHelp = `
_MODEL_ - CPU model name (e.g., {{ bt }}Intel Xeon Gold 6338{{ bt }}). Optional.
`

const cpuThreadsPerCoreFlagHelp = `
_THREADS_ - Number of threads per CPU core. Optional.
`

const memoryTotalGbFlagHelp = `
_MEMORY_ - Total memory in gigabytes. Must be greater than zero.
`

const memoryTypeFlagHelp = `
_TYPE_ - Memory type (e.g., {{ bt }}DDR4{{ bt }}, {{ bt }}DDR5{{ bt }}). Optional.
`
