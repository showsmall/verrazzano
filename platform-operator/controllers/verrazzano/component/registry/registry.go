// Copyright (c) 2021, Oracle and/or its affiliates.
// Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl.

package registry

import (
	"fmt"
	"path/filepath"

	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/coherence"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/externaldns"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/mysql"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/nginx"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/oam"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/rancher"

	"github.com/verrazzano/verrazzano/platform-operator/constants"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/appoper"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/helm"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/istio"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/keycloak"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/spi"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/verrazzano"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/weblogic"
	"github.com/verrazzano/verrazzano/platform-operator/internal/config"
)

type GetCompoentsFnType func() []spi.Component

var getComponentsFn = getComponents

// OverrideGetComponentsFn Allows overriding the set of registry components for testing purposes
func OverrideGetComponentsFn(fnType GetCompoentsFnType) {
	getComponentsFn = fnType
}

// ResetGetComponentsFn Restores the GetComponents implementation to the default if it's been overridden for testing
func ResetGetComponentsFn() {
	getComponentsFn = getComponents
}

// GetComponents returns the list of components that are installable and upgradeable.
// The components will be processed in the order items in the array
// The components will be processed in the order items in the array
func GetComponents() []spi.Component {
	return getComponentsFn()
}

const defaultImagePullSecretKeyName = "imagePullSecrets[0].name"

// getComponents is the internal impl function for GetComponents, to allow overriding it for testing purposes
func getComponents() []spi.Component {
	overridesDir := config.GetHelmOverridesDir()
	helmChartsDir := config.GetHelmChartsDir()
	thirdPartyChartsDir := config.GetThirdPartyDir()
	injectedSystemNamespaces := config.GetInjectedSystemNamespaces()

	return []spi.Component{
		helm.HelmComponent{
			ReleaseName:             nginx.ComponentName,
			ChartDir:                filepath.Join(thirdPartyChartsDir, "ingress-nginx"), // Note name is different than release name
			ChartNamespace:          nginx.ComponentNamespace,
			IgnoreNamespaceOverride: true,
			SupportsOperatorInstall: true,
			ImagePullSecretKeyname:  defaultImagePullSecretKeyName,
			ValuesFile:              filepath.Join(overridesDir, nginx.ValuesFileOverride),
			PreInstallFunc:          nginx.PreInstall,
			AppendOverridesFunc:     nginx.AppendOverrides,
			PostInstallFunc:         nginx.PostInstall,
			Dependencies:            []string{istio.ComponentName},
			ReadyStatusFunc:         nginx.IsReady,
		},
		helm.HelmComponent{
			ReleaseName:             "cert-manager",
			ChartDir:                filepath.Join(thirdPartyChartsDir, "cert-manager"),
			ChartNamespace:          "cert-manager",
			IgnoreNamespaceOverride: true,
			ValuesFile:              filepath.Join(overridesDir, "cert-manager-values.yaml"),
		},
		helm.HelmComponent{
			ReleaseName:             externaldns.ComponentName,
			ChartDir:                filepath.Join(thirdPartyChartsDir, externaldns.ComponentName),
			ChartNamespace:          "cert-manager",
			IgnoreNamespaceOverride: true,
			ValuesFile:              filepath.Join(overridesDir, "external-dns-values.yaml"),
		},
		helm.HelmComponent{
			ReleaseName:             rancher.ComponentName,
			ChartDir:                filepath.Join(thirdPartyChartsDir, rancher.ComponentName),
			ChartNamespace:          "cattle-system",
			IgnoreNamespaceOverride: true,
			ValuesFile:              filepath.Join(overridesDir, "rancher-values.yaml"),
		},
		helm.HelmComponent{
			ReleaseName:             verrazzano.ComponentName,
			ChartDir:                filepath.Join(helmChartsDir, verrazzano.ComponentName),
			ChartNamespace:          constants.VerrazzanoSystemNamespace,
			IgnoreNamespaceOverride: true,
			ResolveNamespaceFunc:    verrazzano.ResolveVerrazzanoNamespace,
			PreUpgradeFunc:          verrazzano.VerrazzanoPreUpgrade,
		},
		helm.HelmComponent{
			ReleaseName:             coherence.ComponentName,
			ChartDir:                filepath.Join(thirdPartyChartsDir, coherence.ComponentName),
			ChartNamespace:          constants.VerrazzanoSystemNamespace,
			IgnoreNamespaceOverride: true,
			SupportsOperatorInstall: true,
			ImagePullSecretKeyname:  defaultImagePullSecretKeyName,
			ValuesFile:              filepath.Join(overridesDir, "coherence-values.yaml"),
			ReadyStatusFunc:         coherence.IsCoherenceOperatorReady,
		},
		helm.HelmComponent{
			ReleaseName:             weblogic.ComponentName,
			ChartDir:                filepath.Join(thirdPartyChartsDir, weblogic.ComponentName),
			ChartNamespace:          constants.VerrazzanoSystemNamespace,
			IgnoreNamespaceOverride: true,
			SupportsOperatorInstall: true,
			ImagePullSecretKeyname:  defaultImagePullSecretKeyName,
			ValuesFile:              filepath.Join(overridesDir, "weblogic-values.yaml"),
			PreInstallFunc:          weblogic.WeblogicOperatorPreInstall,
			AppendOverridesFunc:     weblogic.AppendWeblogicOperatorOverrides,
			Dependencies:            []string{istio.ComponentName},
			ReadyStatusFunc:         weblogic.IsWeblogicOperatorReady,
		},
		helm.HelmComponent{
			ReleaseName:             oam.ComponentName,
			ChartDir:                filepath.Join(thirdPartyChartsDir, oam.ComponentName),
			ChartNamespace:          constants.VerrazzanoSystemNamespace,
			IgnoreNamespaceOverride: true,
			SupportsOperatorInstall: true,
			ValuesFile:              filepath.Join(overridesDir, "oam-kubernetes-runtime-values.yaml"),
			ImagePullSecretKeyname:  defaultImagePullSecretKeyName,
			ReadyStatusFunc:         oam.IsOAMReady,
		},
		helm.HelmComponent{
			ReleaseName:             appoper.ComponentName,
			ChartDir:                filepath.Join(helmChartsDir, appoper.ComponentName),
			ChartNamespace:          constants.VerrazzanoSystemNamespace,
			IgnoreNamespaceOverride: true,
			SupportsOperatorInstall: true,
			ValuesFile:              filepath.Join(overridesDir, "verrazzano-application-operator-values.yaml"),
			AppendOverridesFunc:     appoper.AppendApplicationOperatorOverrides,
			ImagePullSecretKeyname:  "global.imagePullSecrets[0]",
			ReadyStatusFunc:         appoper.IsApplicationOperatorReady,
			Dependencies:            []string{"oam-kubernetes-runtime"},
			PreUpgradeFunc:          appoper.ApplyCRDYaml,
		},
		mysql.NewComponent(),
		keycloak.NewComponent(),
		istio.IstioComponent{
			ValuesFile:               filepath.Join(overridesDir, "istio-cr.yaml"),
			InjectedSystemNamespaces: injectedSystemNamespaces,
		},
	}
}

func FindComponent(releaseName string) (bool, spi.Component) {
	for _, comp := range GetComponents() {
		if comp.Name() == releaseName {
			return true, comp
		}
	}
	return false, &helm.HelmComponent{}
}

// ComponentDependenciesMet Checks if the declared dependencies for the component are ready and available
func ComponentDependenciesMet(c spi.Component, context spi.ComponentContext) bool {
	log := context.Log()
	trace, err := checkDependencies(c, context, nil)
	if err != nil {
		log.Error(err.Error())
		return false
	}
	if len(trace) == 0 {
		log.Infof("No dependencies declared for %s", c.Name())
		return true
	}
	log.Infof("Trace results for %s: %v", c.Name(), trace)
	for _, value := range trace {
		if !value {
			return false
		}
	}
	return true
}

// checkDependencies Check the ready state of any dependencies and check for cycles
func checkDependencies(c spi.Component, context spi.ComponentContext, trace map[string]bool) (map[string]bool, error) {
	for _, dependencyName := range c.GetDependencies() {
		if trace == nil {
			trace = make(map[string]bool)
		}
		if _, ok := trace[dependencyName]; ok {
			return trace, fmt.Errorf("Illegal state, dependency cycle found for %s: %s", c.Name(), dependencyName)
		}
		found, dependency := FindComponent(dependencyName)
		if !found {
			return trace, fmt.Errorf("Illegal state, declared dependency not found for %s: %s", c.Name(), dependencyName)
		}
		if trace, err := checkDependencies(dependency, context, trace); err != nil {
			return trace, err
		}
		if !dependency.IsReady(context) {
			trace[dependencyName] = false // dependency is not ready
			continue
		}
		trace[dependencyName] = true // dependency is ready
	}
	return trace, nil
}
