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
	"math/rand"
	"time"

	operatorv1alpha1 "github.com/ibm/ibm-licensing-hub-operator/pkg/apis/operator/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// cannot set to const due to k8s struct needing pointers to primitive types

var TrueVar = true
var FalseVar = false

const DefaultReceiverImage = "quay.io/opencloudio/ibm-license-advisor-receiver:1.2.0"
const DefaultDatabaseImage = "quay.io/opencloudio/ibm-license-advisor-db:1.2.0"

const DatabaseConfigSecretName = "license-service-hub-db-config"
const PostgresPasswordKey = "POSTGRES_PASSWORD" // #nosec
const PostgresUserKey = "POSTGRES_USER"
const PostgresDatabaseNameKey = "POSTGRES_DATABASE_NAME"
const PostgresPgDataKey = "POSTGRES_PGDATA"

const DatabaseUser = "postgres"
const DatabaseName = "postgres"
const DatabaseMountPoint = "/var/lib/postgresql/data"
const PgData = DatabaseMountPoint + "/pgdata"

const DatabaseContainerName = "database"
const ReceiverContainerName = "receiver"
const ReceiverPort = 8080

const LicensingHubResourceBase = "ibm-licensing-hub-service"
const LicensingHubComponentName = "ibm-licensing-hub-service-svc"
const LicensingHubReleaseName = "ibm-licensing-hub-service"

// Important product values needed for annotations
const LicensingProductName = "IBM Cloud Platform Common Services"
const LicensingProductID = "068a62892a1e4db39641342e592daa25"
const LicensingProductMetric = "FREE"
const LicensingProductVersion = "3.4.0"

const randStringCharset string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const randStringCharsetLength = len(randStringCharset)

var defaultSecretMode int32 = 420
var seconds60 int64 = 60

func RandString(length int) string {
	randFunc := rand.New(rand.NewSource(time.Now().UnixNano()))
	outputStringByte := make([]byte, length)
	for i := 0; i < length; i++ {
		outputStringByte[i] = randStringCharset[randFunc.Intn(randStringCharsetLength)]
	}
	return string(outputStringByte)
}

func GetResourceName(instance *operatorv1alpha1.IBMLicensingHub) string {
	return LicensingHubResourceBase + "-" + instance.GetName()
}

func LabelsForLicensingHubSelector(instance *operatorv1alpha1.IBMLicensingHub) map[string]string {
	return map[string]string{"app": GetResourceName(instance), "component": LicensingHubComponentName, "licensing_cr": instance.GetName()}
}

func LabelsForLicensingHubMeta(instance *operatorv1alpha1.IBMLicensingHub) map[string]string {
	return map[string]string{"app.kubernetes.io/name": GetResourceName(instance), "app.kubernetes.io/component": LicensingHubComponentName,
		"app.kubernetes.io/managed-by": "operator", "app.kubernetes.io/instance": LicensingHubReleaseName, "release": LicensingHubReleaseName}
}

func AnnotationsForPod() map[string]string {
	return map[string]string{"productName": LicensingProductName,
		"productID": LicensingProductID, "productVersion": LicensingProductVersion, "productMetric": LicensingProductMetric,
		"clusterhealth.ibm.com/dependencies": "metering"}
}

func LabelsForLicensingHubPod(instance *operatorv1alpha1.IBMLicensingHub) map[string]string {
	podLabels := LabelsForLicensingHubMeta(instance)
	selectorLabels := LabelsForLicensingHubSelector(instance)
	for key, value := range selectorLabels {
		podLabels[key] = value
	}
	return podLabels
}

func Contains(s []corev1.LocalObjectReference, e corev1.LocalObjectReference) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
