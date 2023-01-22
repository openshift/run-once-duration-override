package runoncedurationoverride

import (
	"fmt"
	"io"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// RunOnceDurationOverride is the configuration for the RunOnceDurationOverride
// admission controller which overrides activeDeadlineSeconds for run-once pods.
type RunOnceDurationOverride struct {
	metav1.TypeMeta `json:",inline"`
	Spec            RunOnceDurationOverrideSpec `json:"spec,omitempty"`
}

type RunOnceDurationOverrideSpec struct {
	ActiveDeadlineSeconds int64 `json:"ActiveDeadlineSeconds"`
}

type Config struct {
	ActiveDeadlineSeconds int64
}

func (c *Config) String() string {
	return fmt.Sprintf("ActiveDeadlineSeconds=%d", c.ActiveDeadlineSeconds)
}

func ConvertExternalConfig(object *RunOnceDurationOverride) *Config {
	return &Config{
		ActiveDeadlineSeconds: object.Spec.ActiveDeadlineSeconds,
	}
}

// DecodeUnstructured decodes a raw stream into a an
// unstructured.Unstructured instance.
func Decode(reader io.Reader) (object *RunOnceDurationOverride, err error) {
	decoder := yaml.NewYAMLOrJSONDecoder(reader, 30)

	c := &RunOnceDurationOverride{}
	if err = decoder.Decode(c); err != nil {
		return
	}

	object = c
	return
}

func DecodeWithFile(path string) (object *RunOnceDurationOverride, err error) {
	reader, openErr := os.Open(path)
	if err != nil {
		err = fmt.Errorf("unable to load file %s: %s", path, openErr)
		return
	}

	object, err = Decode(reader)
	return
}
