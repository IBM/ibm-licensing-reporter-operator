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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func getReceiverEnvironmentVariables() []corev1.EnvVar {
	return []corev1.EnvVar{
		{
			Name:  "POSTGRESQL_USER",
			Value: DatabaseUser,
		},
		{
			Name: "POSTGRESQL_PASSWORD",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: DatabaseConfigSecretName,
					},
					Key: PostgresPasswordKey,
				},
			},
		},
		{
			Name:  "POSTGRESQL_DATABASE",
			Value: DatabaseName,
		},
	}

}

func getReceiverProbeHandler() corev1.Handler {
	return corev1.Handler{
		HTTPGet: &corev1.HTTPGetAction{
			Path: "/",
			Port: intstr.IntOrString{
				Type:   intstr.Int,
				IntVal: receiverServicePort.IntVal,
			},
			Scheme: "HTTPS",
		},
	}
}

func GetReceiverContainer(spec operatorv1alpha1.IBMLicensingHubSpec, instance *operatorv1alpha1.IBMLicensingHub) corev1.Container {
	container := getContainerBase(spec.ReceiverContainer)
	container.ImagePullPolicy = corev1.PullAlways
	container.Env = getReceiverEnvironmentVariables()
	container.Resources = instance.Spec.ReceiverContainer.Resources
	container.VolumeMounts = getLicensingHubVolumeMounts()
	container.Name = ReceiverContainerName
	container.Ports = []corev1.ContainerPort{
		{
			ContainerPort: ReceiverPort,
			Protocol:      corev1.ProtocolTCP,
		},
	}
	container.LivenessProbe = getLivenessProbe(getReceiverProbeHandler())
	container.ReadinessProbe = getReadinessProbe(getReceiverProbeHandler())
	return container
}
