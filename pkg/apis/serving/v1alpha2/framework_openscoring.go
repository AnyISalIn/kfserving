/*
Copyright 2019 kubeflow.org.
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

package v1alpha2

import (
	"fmt"
	"strings"

	"github.com/kubeflow/kfserving/pkg/constants"
	v1 "k8s.io/api/core/v1"
)

var _ Predictor = (*OpenScoringSpec)(nil)

func (p *OpenScoringSpec) GetStorageUri() string {
	return p.StorageURI
}

func (p *OpenScoringSpec) GetResourceRequirements() *v1.ResourceRequirements {
	// return the ResourceRequirements value if set on the spec
	return &p.Resources
}

func (p *OpenScoringSpec) GetContainer(modelName string, parallelism int, config *InferenceServicesConfig) *v1.Container {
	arguments := []string{
		"--model-dir", constants.DefaultModelLocalMountPath,
	}
	return &v1.Container{
		Image:     config.Predictors.OpenScoring.ContainerImage + ":" + p.RuntimeVersion,
		Name:      constants.InferenceServiceContainerName,
		Resources: p.Resources,
		Args:      arguments,
	}
}

func (p *OpenScoringSpec) ApplyDefaults(config *InferenceServicesConfig) {
	if p.RuntimeVersion == "" {
		p.RuntimeVersion = config.Predictors.OpenScoring.DefaultImageVersion
	}
	setResourceRequirementDefaults(&p.Resources)
}

func (p *OpenScoringSpec) Validate(config *InferenceServicesConfig) error {

	if isGPUEnabled(p.Resources) && !strings.Contains(p.RuntimeVersion, PyTorchServingGPUSuffix) {
		return fmt.Errorf(InvalidPyTorchRuntimeIncludesGPU)
	}

	if !isGPUEnabled(p.Resources) && strings.Contains(p.RuntimeVersion, PyTorchServingGPUSuffix) {
		return fmt.Errorf(InvalidPyTorchRuntimeExcludesGPU)
	}
	return nil
}
