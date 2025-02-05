// Copyright (c) 2022, Oracle and/or its affiliates.
// Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl.

package fluentd

import (
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	"github.com/verrazzano/verrazzano/pkg/test/framework"
	vzapi "github.com/verrazzano/verrazzano/platform-operator/apis/verrazzano/v1alpha1"
	pcons "github.com/verrazzano/verrazzano/platform-operator/constants"
	"github.com/verrazzano/verrazzano/tests/e2e/pkg"
)

const (
	labelValidation = "f:platform-lcm.fluentd-update-validation"
	opensearchURL   = "https://opensearch.example.com:9200"
)

var (
	t        = framework.NewTestFramework("update fluentd")
	tempuuid = uuid.NewString()[:7]
	extEsSec = "my-extsec-" + tempuuid
	wrongSec = "wrong-sec-" + tempuuid
	ociLgSec = "my-ocilog-" + tempuuid
	sysLogID = "my-sysLog-" + tempuuid
	defLogID = "my-defLog-" + tempuuid
)

var _ = t.AfterSuite(func() {
	pkg.DeleteSecret(pcons.VerrazzanoInstallNamespace, extEsSec)
	pkg.DeleteSecret(pcons.VerrazzanoInstallNamespace, wrongSec)
	ValidateDaemonset(pkg.VmiESURL, pkg.VmiESInternalSecret, "")
})

var _ = t.Describe("Update Fluentd", Label("f:platform-lcm.update"), func() {
	t.Describe("fluentd verify", Label("f:platform-lcm.fluentd-verify"), func() {
		t.It("fluentd default config", func() {
			ValidateDaemonset(pkg.VmiESURL, pkg.VmiESInternalSecret, "")
		})
	})

	t.Describe("Validate external Opensearch config", Label(labelValidation), func() {
		t.It("secret validation", func() {
			m := &FluentdModifier{Component: vzapi.FluentdComponent{
				ElasticsearchSecret: extEsSec + "missing",
				ElasticsearchURL:    opensearchURL,
			}}
			ValidateUpdate(m, "must be created")
		})
	})

	t.Describe("Update external Opensearch", Label("f:platform-lcm.fluentd-external-opensearch"), func() {
		t.It("external Opensearch", func() {
			pkg.CreateCredentialsSecret(pcons.VerrazzanoInstallNamespace, extEsSec, "user", "pw", map[string]string{})
			m := &FluentdModifier{Component: vzapi.FluentdComponent{
				ElasticsearchSecret: extEsSec,
				ElasticsearchURL:    opensearchURL,
			}}
			ValidateUpdate(m, "")
			ValidateDaemonset(opensearchURL, extEsSec, "")
		})
	})

	t.Describe("Validate OCI logging config", Label(labelValidation), func() {
		t.It("secret validation", func() {
			m := &FluentdModifier{Component: vzapi.FluentdComponent{
				OCI: &vzapi.OciLoggingConfiguration{APISecret: wrongSec},
			}}
			ValidateUpdate(m, "must be created")
			pkg.CreateCredentialsSecret(pcons.VerrazzanoInstallNamespace, wrongSec, "api", "pw", map[string]string{})
			ValidateUpdate(m, "Did not find OCI configuration")
		})
	})

	t.Describe("Update OCI logging", Label("f:platform-lcm.fluentd-oci-logging"), func() {
		t.It(" OCI logging", func() {
			createOciLoggingSecret(ociLgSec)
			m := &FluentdModifier{Component: vzapi.FluentdComponent{OCI: &vzapi.OciLoggingConfiguration{
				APISecret:       ociLgSec,
				SystemLogID:     sysLogID,
				DefaultAppLogID: defLogID,
			}}}
			ValidateUpdate(m, "")
			ValidateDaemonset("", "", ociLgSec)
			ValidateConfigMap(sysLogID, defLogID)
		})
	})

	t.Describe("Validate extra Volume Mounts", Label(labelValidation), func() {
		t.It("extraVolumeMounts validation", func() {
			m := &FluentdModifier{Component: vzapi.FluentdComponent{
				ExtraVolumeMounts: []vzapi.VolumeMount{{Source: "/var/log"}},
			}}
			ValidateUpdate(m, "duplicate mount path found")
		})
	})

	t.Describe("Update extraVolumeMounts", Label("f:platform-lcm.fluentd-extra-volume-mounts"), func() {
		t.It("extraVolumeMounts", func() {
			vm := vzapi.VolumeMount{Source: "/var/log", Destination: "/home/varLog"}
			m := &FluentdModifier{Component: vzapi.FluentdComponent{
				ExtraVolumeMounts: []vzapi.VolumeMount{vm},
			}}
			ValidateUpdate(m, "")
			ValidateDaemonset(pkg.VmiESURL, pkg.VmiESInternalSecret, "", vm)
		})
	})
})
