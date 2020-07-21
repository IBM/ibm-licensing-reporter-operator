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
)

const APISecretTokenVolumeName = "api-token"
const persistentVolumeClaimVolumeName = "data"

func getLicensingHubVolumeMounts() []corev1.VolumeMount {
	return []corev1.VolumeMount{
		{
			Name:      APISecretTokenVolumeName,
			MountPath: "/opt/ibm/licensing",
			ReadOnly:  true,
		},
	}
}

func getLicensingHubDatabaseVolumeMounts() []corev1.VolumeMount {
	return []corev1.VolumeMount{
		{
			Name:      persistentVolumeClaimVolumeName,
			MountPath: DatabaseMountPoint,
		},
	}
}

func getLicensingHubVolumes(spec operatorv1alpha1.IBMLicensingHubSpec) []corev1.Volume {
	volumes := []corev1.Volume{

		{
			Name: APISecretTokenVolumeName,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  spec.APISecretToken,
					DefaultMode: &defaultSecretMode,
				},
			},
		},
		{
			Name: persistentVolumeClaimVolumeName,
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: PersistenceVolumeClaimName,
				},
			},
		},
	}
	return volumes
}
