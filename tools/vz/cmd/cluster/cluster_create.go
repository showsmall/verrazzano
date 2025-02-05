// Copyright (c) 2022, Oracle and/or its affiliates.
// Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl.

package cluster

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/verrazzano/verrazzano/pkg/capi"
	cmdhelpers "github.com/verrazzano/verrazzano/tools/vz/cmd/helpers"
	"github.com/verrazzano/verrazzano/tools/vz/pkg/constants"
	"github.com/verrazzano/verrazzano/tools/vz/pkg/helpers"
)

const (
	createSubCommandName = "create"
	createHelpShort      = "Verrazzano cluster create"
	createHelpLong       = `Creates a new local cluster`
	createHelpExample    = `vz cluster create --name mycluster`
)

func newSubcmdCreate(vzHelper helpers.VZHelper) *cobra.Command {
	cmd := cmdhelpers.NewCommand(vzHelper, createSubCommandName, createHelpShort, createHelpLong)
	cmd.Example = createHelpExample
	cmd.PersistentFlags().String(constants.ClusterNameFlagName, constants.ClusterNameFlagDefault, constants.ClusterNameFlagHelp)
	cmd.PersistentFlags().String(constants.ClusterTypeFlagName, constants.ClusterTypeFlagDefault, constants.ClusterTypeFlagHelp)
	cmd.PersistentFlags().String(constants.ClusterImageFlagName, constants.ClusterImageFlagDefault, constants.ClusterImageFlagHelp)
	// the image and type flags should be hidden since they are not intended for general use
	cmd.PersistentFlags().MarkHidden(constants.ClusterTypeFlagName)
	cmd.PersistentFlags().MarkHidden(constants.ClusterImageFlagName)
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return runCmdClusterCreate(cmd, args)
	}

	return cmd
}

func runCmdClusterCreate(cmd *cobra.Command, args []string) error {
	clusterName, err := cmd.PersistentFlags().GetString(constants.ClusterNameFlagName)
	if err != nil {
		return fmt.Errorf("Failed to get the %s flag: %v", constants.ClusterNameFlagName, err)
	}

	clusterType, err := cmd.PersistentFlags().GetString(constants.ClusterTypeFlagName)
	if err != nil {
		return fmt.Errorf("Failed to get the %s flag: %v", constants.ClusterTypeFlagName, err)
	}

	clusterImg, err := cmd.PersistentFlags().GetString(constants.ClusterImageFlagName)
	if err != nil {
		return fmt.Errorf("Failed to get the %s flag: %v", constants.ClusterImageFlagName, err)
	}

	cluster, err := capi.NewBoostrapCluster(capi.ClusterConfig{
		ClusterName:    clusterName,
		Type:           clusterType,
		ContainerImage: clusterImg,
	})
	if err != nil {
		return err
	}
	if err := cluster.Create(); err != nil {
		return err
	}
	fmt.Printf("Cluster %s created successfully, initializing...\n", clusterName)
	if err := cluster.Init(); err != nil {
		return err
	}
	fmt.Println("Cluster initialization complete")
	fmt.Printf("To get the kubeconfig for this cluster, run: vz cluster get-kubeconfig --name %s (for more details, run vz cluster get-kubeconfig -h)\n", clusterName)
	return nil
}
