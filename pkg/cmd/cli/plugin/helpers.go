/*
Copyright The Velero contributors.

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

package plugin

import (
	"context"

	"github.com/pkg/errors"
	appsv1api "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"

	"github.com/adi-bhardwaj/velero-modified/pkg/install"
)

// veleroDeployment returns a Velero deployment object, selected using a label.
func veleroDeployment(ctx context.Context, kubeClient kubernetes.Interface, namespace string) (*appsv1api.Deployment, error) {
	veleroLabels := labels.FormatLabels(install.Labels())

	deployList, err := kubeClient.
		AppsV1().
		Deployments(namespace).
		List(ctx, metav1.ListOptions{
			LabelSelector: veleroLabels,
		})
	if err != nil {
		return nil, err
	}

	if len(deployList.Items) < 1 {
		return nil, errors.New("Velero deployment not found")
	}

	return &deployList.Items[0], nil
}
