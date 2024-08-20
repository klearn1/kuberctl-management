/*
Copyright 2023 The Kubernetes Authors.

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
package e2enode

import (
	"bytes"
	"context"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	watchtools "k8s.io/client-go/tools/watch"
	"k8s.io/kubernetes/test/e2e/framework"
	e2epod "k8s.io/kubernetes/test/e2e/framework/pod"
	admissionapi "k8s.io/pod-security-admission/api"
	"strings"
	"time"
)

var _ = SIGDescribe(framework.WithNodeConformance(), "Shortened Grace Period", func() {
	f := framework.NewDefaultFramework("shortened-grace-period")
	f.NamespacePodSecurityEnforceLevel = admissionapi.LevelPrivileged
	ginkgo.Context("When repeatedly deleting pods", func() {
		var podClient *e2epod.PodClient
		var dc dynamic.Interface
		var ns string
		var podName = "test-shortened-grace"
		var ctx = context.Background()
		var rcResource = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
		const (
			gracePeriod      = 10000
			gracePeriodShort = 100
		)
		ginkgo.BeforeEach(func() {
			ns = f.Namespace.Name
			dc = f.DynamicClient
			podClient = e2epod.NewPodClient(f)
		})
		ginkgo.It("shorter grace period of a second command overrides the longer grace period of a first command", func() {
			testRcNamespace := ns
			expectedWatchEvents := []watch.Event{
				{Type: watch.Modified},
				{Type: watch.Modified},
				{Type: watch.Deleted},
			}
			eventFound := false
			callback := func(retryWatcher *watchtools.RetryWatcher) (actualWatchEvents []watch.Event) {
				podClient.CreateSync(ctx, getGracePeriodTestPod(podName, testRcNamespace, gracePeriod))

				w, err := podClient.Watch(context.TODO(), metav1.ListOptions{LabelSelector: "test-shortened-grace=true"})
				framework.ExpectNoError(err, "failed to watch")
				err = podClient.Delete(ctx, podName, *metav1.NewDeleteOptions(gracePeriod))
				time.Sleep(3 * time.Second)
				framework.ExpectNoError(err, "failed to delete pod")
				ctxUntil, cancel := context.WithTimeout(ctx, 20*time.Second)
				defer cancel()
				_, err = watchtools.UntilWithoutRetry(ctxUntil, w, func(watchEvent watch.Event) (bool, error) {
					if watchEvent.Type != watch.Modified {
						return false, nil
					}
					actualWatchEvents = append(actualWatchEvents, watchEvent)
					eventFound = true
					return true, nil
				})
				framework.ExpectNoError(err, "Wait until condition with watch events should not return an error")
				if !eventFound {
					framework.Failf("failed to find %v event", watch.Modified)
				}
				w, err = podClient.Watch(context.TODO(), metav1.ListOptions{LabelSelector: "test-shortened-grace=true"})
				framework.ExpectNoError(err, "failed to watch")
				err = podClient.Delete(ctx, podName, *metav1.NewDeleteOptions(gracePeriodShort))
				framework.ExpectNoError(err, "failed to delete pod")
				time.Sleep(3 * time.Second)
				ctxUntil, cancel = context.WithTimeout(ctx, 20*time.Second)
				defer cancel()
				eventFound = false
				_, err = watchtools.UntilWithoutRetry(ctxUntil, w, func(watchEvent watch.Event) (bool, error) {
					if watchEvent.Type != watch.Modified {
						return false, nil
					}
					actualWatchEvents = append(actualWatchEvents, watchEvent)
					eventFound = true
					return true, nil
				})
				framework.ExpectNoError(err, "Wait until condition with watch events should not return an error")
				if !eventFound {
					framework.Failf("failed to find %v event", watch.Modified)
				}
				// Get pod logs.
				logs, err := podClient.GetLogs(podName, &v1.PodLogOptions{}).Stream(ctx)
				framework.ExpectNoError(err, "failed to get pod logs")
				defer func() {
					if err := logs.Close(); err != nil {
						framework.ExpectNoError(err, "failed to log close")
					}
				}()
				buf := new(bytes.Buffer)
				_, err = buf.ReadFrom(logs)
				if err != nil {
					framework.ExpectNoError(err, "failed to read from")
				}
				podLogs := buf.String()
				// Verify the number of SIGINT
				gomega.Expect(strings.Count(podLogs, "SIGINT 1")).To(gomega.Equal(1), fmt.Sprintf("unexpected number of SIGINT 1 entries in pod logs.logs is :%s", podLogs))
				gomega.Expect(strings.Count(podLogs, "SIGINT 2")).To(gomega.Equal(1), fmt.Sprintf("unexpected number of SIGINT 2 entries in pod logs.logs is :%s", podLogs))
				w, err = podClient.Watch(context.TODO(), metav1.ListOptions{LabelSelector: "test-shortened-grace=true"})
				framework.ExpectNoError(err, "failed to watch")
				ctxUntil, cancel = context.WithTimeout(ctx, 15*time.Second)
				defer cancel()
				eventFound = false
				_, err = watchtools.UntilWithoutRetry(ctxUntil, w, func(watchEvent watch.Event) (bool, error) {
					if watchEvent.Type != watch.Deleted {
						return false, nil
					}
					actualWatchEvents = append(actualWatchEvents, watchEvent)
					eventFound = true
					return true, nil
				})
				framework.ExpectNoError(err, "Wait until condition with watch events should not return an error")
				if !eventFound {
					framework.Failf("failed to find %v event", watch.Deleted)
				}
				return expectedWatchEvents
			}
			framework.WatchEventSequenceVerifier(ctx, dc, rcResource, ns, podName, metav1.ListOptions{LabelSelector: "test-shortened-grace=true"}, expectedWatchEvents, callback, func() (err error) {
				return err
			})
		})
	})
})

func getGracePeriodTestPod(name, testRcNamespace string, gracePeriod int64) *v1.Pod {
	pod := &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"test-shortened-grace": "true",
			},
			Namespace: testRcNamespace,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:    name,
					Image:   busyboxImage,
					Command: []string{"sh", "-c"},
					Args: []string{`
_term() {
 if [ "$COUNT" -eq 0 ]; then
   echo "SIGINT 1"
 elif [ "$COUNT" -eq 1 ]; then
   echo "SIGINT 2"
   sleep 80
   exit 0
 fi
 COUNT=$((COUNT + 1))
}
COUNT=0
trap _term SIGTERM
while true; do
 sleep 1
done
`},
				},
			},
			TerminationGracePeriodSeconds: &gracePeriod,
		},
	}
	return pod
}
