package runoncedurationoverride

import (
	"encoding/json"
	"fmt"
	"os"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog"

	admissionresponse "github.com/openshift/run-once-duration-override/pkg/response"
)

const (
	Resource = "runoncedurationoverrides"
	Singular = "runoncedurationoverride"
	Name     = "runoncedurationoverride"
)

const (
	configurationEnvName = "CONFIGURATION_PATH"
)

// ConfigLoaderFunc loads a Config object from appropriate source and returns it.
type ConfigLoaderFunc func() (config *Config, err error)

// Admission interface encapsulates the admission logic for ClusterResourceOverride plugin.
type Admission interface {
	// GetConfiguration returns the configuration in use by the admission logic.
	GetConfiguration() *Config

	// IsApplicable returns true if the given resource inside the request is
	// applicable to this admission controller. Otherwise it returns false.
	IsApplicable(request *admissionv1.AdmissionRequest) bool

	// IsExempt returns true if the given resource is exempt from being admitted.
	// Otherwise it returns false. On any error, response is set with appropriate
	// status and error message.
	// If response is not nil, the caller should not proceed with the admission.
	IsExempt(request *admissionv1.AdmissionRequest) (exempt bool, response *admissionv1.AdmissionResponse)

	// Admit makes an attempt to admit the specified resource in the request.
	// It returns an AdmissionResponse that is set appropriately. On success,
	// the response should contain the patch for update.
	Admit(admissionSpec *admissionv1.AdmissionRequest) *admissionv1.AdmissionResponse
}

// NewInClusterAdmission returns a new instance of Admission that is appropriate
// to be consumed in cluster.
func NewInClusterAdmission() (admission Admission, err error) {
	configLoader := func() (config *Config, err error) {
		configPath := os.Getenv(configurationEnvName)
		if configPath == "" {
			err = fmt.Errorf("name=%s no configuration file specified, env var %s is not set", Name, configurationEnvName)
			return
		}

		externalConfig, decodeErr := DecodeWithFile(configPath)
		if decodeErr != nil {
			err = fmt.Errorf("name=%s file=%s failed to decode configuration - %s", Name, configPath, decodeErr.Error())
			return
		}

		config = ConvertExternalConfig(externalConfig)
		return
	}

	return NewAdmission(configLoader)
}

// NewInClusterAdmission returns a new instance of Admission that is appropriate
// to be consumed in cluster.
func NewAdmission(configLoaderFunc ConfigLoaderFunc) (admission Admission, err error) {
	config, configLoadErr := configLoaderFunc()
	if configLoadErr != nil {
		err = fmt.Errorf("name=%s failed to load configuration - %s", Name, configLoadErr.Error())
		return
	}

	return &runOnceDurationOverrideAdmission{
		config: config,
	}, nil
}

type runOnceDurationOverrideAdmission struct {
	config *Config
}

func (p *runOnceDurationOverrideAdmission) GetConfiguration() *Config {
	return p.config
}

func (p *runOnceDurationOverrideAdmission) IsApplicable(request *admissionv1.AdmissionRequest) bool {
	if request.Resource.Resource == string(corev1.ResourcePods) &&
		request.SubResource == "" &&
		(request.Operation == admissionv1.Create || request.Operation == admissionv1.Update) {

		return true
	}

	return false
}

func (p *runOnceDurationOverrideAdmission) IsExempt(request *admissionv1.AdmissionRequest) (exempt bool, response *admissionv1.AdmissionResponse) {
	pod, err := getPod(request)
	if err != nil {
		return false, admissionresponse.WithBadRequest(request, err)
	}

	return isPodExempt(pod), nil
}

func getPod(request *admissionv1.AdmissionRequest) (pod *corev1.Pod, err error) {
	pod = &corev1.Pod{}
	err = json.Unmarshal(request.Object.Raw, pod)
	return
}

func (p *runOnceDurationOverrideAdmission) Admit(request *admissionv1.AdmissionRequest) *admissionv1.AdmissionResponse {
	klog.V(5).Infof("namespace=%s - admitting resource", request.Namespace)

	pod, err := getPod(request)
	if err != nil {
		return admissionresponse.WithBadRequest(request, err)
	}

	mutator, err := NewMutator(p.config)
	if err != nil {
		return admissionresponse.WithInternalServerError(request, err)
	}

	current, err := mutator.Mutate(pod)
	if err != nil {
		return admissionresponse.WithInternalServerError(request, err)
	}

	if current.Spec.ActiveDeadlineSeconds != nil {
		klog.V(5).Infof("namespace=%s pod activeDeadlineSeconds after mutating is: %v", request.Namespace, *current.Spec.ActiveDeadlineSeconds)
	}

	patch, patchErr := Patch(request.Object, current)
	if patchErr != nil {
		return admissionresponse.WithInternalServerError(request, patchErr)
	}

	return admissionresponse.WithPatch(request, patch)
}
