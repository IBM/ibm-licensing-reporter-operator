//
// Copyright 2020 IBM Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package resources

import (
	operatorv1alpha1 "github.com/ibm/ibm-licensing-hub-operator/pkg/apis/operator/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var replicas = int32(1)

var log = logf.Log.WithName("controller_licenseadvisor_deployment")

func GetLicensingDeployment(instance *operatorv1alpha1.IBMLicensingHub) *appsv1.Deployment {
	metaLabels := LabelsForLicensingHubMeta(instance)
	selectorLabels := LabelsForLicensingHubSelector(instance)
	podLabels := LabelsForLicensingHubPod(instance)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      GetResourceName(instance),
			Namespace: instance.GetNamespace(),
			Labels:    metaLabels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: selectorLabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      podLabels,
					Annotations: AnnotationsForPod(),
				},
				Spec: corev1.PodSpec{
					Volumes: getLicensingHubVolumes(instance.Spec),
					Containers: []corev1.Container{
						GetDatabaseContainer(instance.Spec),
						GetReceiverContainer(instance.Spec),
					},
					TerminationGracePeriodSeconds: &seconds60,
					ServiceAccountName:            GetServiceAccountName(instance),
					Affinity: &corev1.Affinity{
						NodeAffinity: &corev1.NodeAffinity{
							RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
								NodeSelectorTerms: []corev1.NodeSelectorTerm{
									{
										MatchExpressions: []corev1.NodeSelectorRequirement{
											{
												Key:      "beta.kubernetes.io/arch",
												Operator: corev1.NodeSelectorOpIn,
												Values:   []string{"amd64"},
											},
										},
									},
								},
							},
						},
					},
					Tolerations: []corev1.Toleration{
						{
							Key:      "dedicated",
							Operator: corev1.TolerationOpExists,
							Effect:   corev1.TaintEffectNoSchedule,
						},
						{
							Key:      "CriticalAddonsOnly",
							Operator: corev1.TolerationOpExists,
						},
					},
				},
			},
		},
	}
	return deployment
}

func ShouldUpdateDeployment(expectedDeployment *appsv1.Deployment, foundDeployment *appsv1.Deployment) bool {

	shouldUpdate := false
	expectedDatabase := findContainer(expectedDeployment.Spec.Template.Spec.Containers, DatabaseContainerName)
	expectedReceiver := findContainer(expectedDeployment.Spec.Template.Spec.Containers, ReceiverContainerName)

	foundDatabase := findContainer(foundDeployment.Spec.Template.Spec.Containers, DatabaseContainerName)
	foundReceiver := findContainer(foundDeployment.Spec.Template.Spec.Containers, ReceiverContainerName)

	if expectedDatabase == nil || expectedReceiver == nil || foundDatabase == nil || foundReceiver == nil {
		return true
	}

	containers := map[*corev1.Container]*corev1.Container{
		foundDatabase: expectedDatabase,
		foundReceiver: expectedReceiver,
	}

	// Checks for every container
	for foundContainer, expectedContainer := range containers {
		containerLogger := log.WithValues("Container.Name", foundContainer.Name)
		if expectedContainer.Image != foundContainer.Image {
			foundContainer.Image = expectedContainer.Image
			containerLogger.Info("Image needs update")
			shouldUpdate = true
		}
		if expectedContainer.ImagePullPolicy != foundContainer.ImagePullPolicy {
			foundContainer.ImagePullPolicy = expectedContainer.ImagePullPolicy
			containerLogger.Info("Image Pull Policy needs update")
			shouldUpdate = true
		}
		if !expectedContainer.Resources.Limits.Cpu().Equal(*foundContainer.Resources.Limits.Cpu()) {
			foundContainer.Resources.Limits[corev1.ResourceCPU] = *expectedContainer.Resources.Limits.Cpu()
			containerLogger.Info("CPU Limit needs update")
			shouldUpdate = true
		}
		if !expectedContainer.Resources.Requests.Cpu().Equal(*foundContainer.Resources.Requests.Cpu()) {
			foundContainer.Resources.Requests[corev1.ResourceCPU] = *expectedContainer.Resources.Requests.Cpu()
			containerLogger.Info("CPU Request needs update")
			shouldUpdate = true
		}
		if !expectedContainer.Resources.Limits.Memory().Equal(*foundContainer.Resources.Limits.Memory()) {
			foundContainer.Resources.Limits[corev1.ResourceMemory] = *expectedContainer.Resources.Limits.Memory()
			containerLogger.Info("Memory Limit needs update")
			shouldUpdate = true
		}
		if !expectedContainer.Resources.Requests.Memory().Equal(*foundContainer.Resources.Requests.Memory()) {
			foundContainer.Resources.Requests[corev1.ResourceMemory] = *expectedContainer.Resources.Requests.Memory()
			containerLogger.Info("Memory Request needs update")
			shouldUpdate = true
		}

	}

	expectedVolume := findVolume(expectedDeployment.Spec.Template.Spec.Volumes, APISecretTokenVolumeName)
	foundVolume := findVolume(foundDeployment.Spec.Template.Spec.Volumes, APISecretTokenVolumeName)

	if expectedVolume.VolumeSource.Secret.SecretName != foundVolume.VolumeSource.Secret.SecretName {
		foundVolume.VolumeSource.Secret.SecretName = expectedVolume.VolumeSource.Secret.SecretName
		log.Info("Api Token Secret needs update")
		shouldUpdate = true
	}

	return shouldUpdate
}

func findContainer(containers []corev1.Container, name string) *corev1.Container {

	for _, container := range containers {
		if container.Name == name {
			return &container
		}
	}
	return nil
}

func findVolume(volumes []corev1.Volume, name string) *corev1.Volume {

	for _, volume := range volumes {
		if volume.Name == name {
			return &volume
		}
	}
	return nil
}
