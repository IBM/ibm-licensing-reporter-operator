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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var (
	receiverServicePort    = intstr.FromInt(ReceiverPort)
	receiverTargetPort     = intstr.FromInt(ReceiverPort)
	receiverTargetPortName = intstr.FromString("receiver-port")
)

func getServiceSpec(instance *operatorv1alpha1.IBMLicensingHub) corev1.ServiceSpec {
	return corev1.ServiceSpec{
		Type: corev1.ServiceTypeClusterIP,
		Ports: []corev1.ServicePort{
			{
				Name:       receiverTargetPortName.String(),
				Port:       receiverServicePort.IntVal,
				TargetPort: receiverTargetPort,
				Protocol:   corev1.ProtocolTCP,
			},
		},
		Selector: LabelsForLicensingHubSelector(instance),
	}
}

func GetReceiverService(instance *operatorv1alpha1.IBMLicensingHub) *corev1.Service {
	metaLabels := LabelsForLicensingHubMeta(instance)
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      LicensingHubResourceBase,
			Namespace: instance.GetNamespace(),
			Labels:    metaLabels,
		},
		Spec: getServiceSpec(instance),
	}
}
