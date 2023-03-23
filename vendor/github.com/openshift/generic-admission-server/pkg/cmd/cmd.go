package cmd

import (
	"os"
	"runtime"

	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/apiserver/pkg/server/mux"
	"k8s.io/component-base/cli"
	"k8s.io/component-base/metrics/legacyregistry"
	"k8s.io/klog/v2"

	"github.com/openshift/generic-admission-server/pkg/apiserver"
	"github.com/openshift/generic-admission-server/pkg/cmd/server"
	"github.com/openshift/run-once-duration-override/pkg/runoncedurationoverride"
	apiServer "k8s.io/apiserver/pkg/server"
	apiserveroptions "k8s.io/apiserver/pkg/server/options"
	restclient "k8s.io/client-go/rest"
)

const (
	DefaultRODOMetricsPort = 9448
)

// AdmissionHook is what callers provide, in the mutating, the validating variant or implementing even both interfaces.
// We define it here to limit how much of the import tree callers have to deal with for this plugin. This means that
// callers need to match levels of apimachinery, api, client-go, and apiserver.
type AdmissionHook apiserver.AdmissionHook
type ValidatingAdmissionHook apiserver.ValidatingAdmissionHook
type MutatingAdmissionHook apiserver.MutatingAdmissionHook

func RunAdmissionServer(admissionHooks ...AdmissionHook) {
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	stopCh := genericapiserver.SetupSignalHandler()

	var LoopbackClientConfig *restclient.Config
	var SecureServing *apiServer.SecureServingInfo

	secureServing := apiserveroptions.NewSecureServingOptions().WithLoopback()
	secureServing.BindPort = DefaultRODOMetricsPort

	if err := secureServing.ApplyTo(&SecureServing, &LoopbackClientConfig); err != nil {
		klog.V(1).Infof("name=%s failed to initialize webhook mectrics - %s", runoncedurationoverride.Name, err.Error())
	}
	pathRecorderMux := mux.NewPathRecorderMux("rodo")
	pathRecorderMux.Handle("/metrics", legacyregistry.HandlerWithReset())
	_, _, err := SecureServing.Serve(pathRecorderMux, 0, stopCh)
	if err != nil {
		klog.Fatalf("failed to start secure server: %v", err)
	}

	klog.V(1).Infof("name=%s admission webhook metrics up", runoncedurationoverride.Name)

	// done to avoid cannot use admissionHooks (type []AdmissionHook) as type []apiserver.AdmissionHook in argument to "github.com/openshift/kubernetes-namespace-reservation/pkg/genericadmissionserver/cmd/server".NewCommandStartAdmissionServer
	var castSlice []apiserver.AdmissionHook
	for i := range admissionHooks {
		castSlice = append(castSlice, admissionHooks[i])
	}

	code := cli.Run(server.NewCommandStartAdmissionServer(os.Stdout, os.Stderr, stopCh, castSlice...))
	os.Exit(code)
}
