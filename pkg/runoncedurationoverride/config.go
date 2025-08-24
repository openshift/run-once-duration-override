/*
Copyright 2023 The run-once-duration-operator Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
	ActiveDeadlineSeconds int64 `json:"activeDeadlineSeconds"`
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
	if openErr != nil {
		err = fmt.Errorf("unable to load file %s: %s", path, openErr)
		return
	}

	object, err = Decode(reader)
	return
}
