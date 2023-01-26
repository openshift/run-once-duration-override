package e2e

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	k8sclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"

	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/openshift/library-go/pkg/operator/resource/resourceapply"
	"github.com/openshift/library-go/pkg/operator/resource/resourceread"

	"github.com/openshift/run-once-duration-override/test/e2e/bindata"
)

const (
	webhookNamespace = "run-once-duration-override"
	webhookName      = "runoncedurationoverride"
)

func TestMain(m *testing.M) {
	if os.Getenv("KUBECONFIG") == "" {
		klog.Errorf("KUBECONFIG environment variable not set")
		os.Exit(1)
	}

	if os.Getenv("IMAGE_FORMAT") == "" {
		klog.Errorf("IMAGE_FORMAT environment variable not set")
		os.Exit(1)
	}

	if os.Getenv("NAMESPACE") == "" {
		klog.Errorf("NAMESPACE environment variable not set")
		os.Exit(1)
	}

	kubeClient := getKubeClientOrDie()

	eventRecorder := events.NewKubeRecorder(kubeClient.CoreV1().Events("default"), "test-e2e", &corev1.ObjectReference{})

	ctx, cancelFnc := context.WithCancel(context.TODO())
	defer cancelFnc()

	assets := []struct {
		path           string
		readerAndApply func(objBytes []byte) error
	}{
		{
			path: "assets/100_namespace.yaml",
			readerAndApply: func(objBytes []byte) error {
				_, _, err := resourceapply.ApplyNamespace(ctx, kubeClient.CoreV1(), eventRecorder, resourceread.ReadNamespaceV1OrDie(objBytes))
				return err
			},
		},
		{
			path: "assets/200_sa.yaml",
			readerAndApply: func(objBytes []byte) error {
				_, _, err := resourceapply.ApplyServiceAccount(ctx, kubeClient.CoreV1(), eventRecorder, resourceread.ReadServiceAccountV1OrDie(objBytes))
				return err
			},
		},
		{
			path: "assets/300_clusterrole_aggregated_apiserver.yaml",
			readerAndApply: func(objBytes []byte) error {
				_, _, err := resourceapply.ApplyClusterRole(ctx, kubeClient.RbacV1(), eventRecorder, resourceread.ReadClusterRoleV1OrDie(objBytes))
				return err
			},
		},
		{
			path: "assets/301_clusterrolebinding_aggregated_apiserver.yaml",
			readerAndApply: func(objBytes []byte) error {
				_, _, err := resourceapply.ApplyClusterRoleBinding(ctx, kubeClient.RbacV1(), eventRecorder, resourceread.ReadClusterRoleBindingV1OrDie(objBytes))
				return err
			},
		},
		{
			path: "assets/302_clusterrolebinding_auth_delegator.yaml",
			readerAndApply: func(objBytes []byte) error {
				_, _, err := resourceapply.ApplyClusterRoleBinding(ctx, kubeClient.RbacV1(), eventRecorder, resourceread.ReadClusterRoleBindingV1OrDie(objBytes))
				return err
			},
		},
		{
			path: "assets/303_rolebinding_auth_reader.yaml",
			readerAndApply: func(objBytes []byte) error {
				_, _, err := resourceapply.ApplyRoleBinding(ctx, kubeClient.RbacV1(), eventRecorder, resourceread.ReadRoleBindingV1OrDie(objBytes))
				return err
			},
		},
		{
			path: "assets/304_scc_hostnetwork_role.yaml",
			readerAndApply: func(objBytes []byte) error {
				_, _, err := resourceapply.ApplyRole(ctx, kubeClient.RbacV1(), eventRecorder, resourceread.ReadRoleV1OrDie(objBytes))
				return err
			},
		},
		{
			path: "assets/305_scc_hostnetwork_rolebinding.yaml",
			readerAndApply: func(objBytes []byte) error {
				_, _, err := resourceapply.ApplyRoleBinding(ctx, kubeClient.RbacV1(), eventRecorder, resourceread.ReadRoleBindingV1OrDie(objBytes))
				return err
			},
		},
		{
			path: "assets/400_secret.yaml",
			readerAndApply: func(objBytes []byte) error {
				_, _, err := resourceapply.ApplySecret(ctx, kubeClient.CoreV1(), eventRecorder, resourceread.ReadSecretV1OrDie(objBytes))
				return err
			},
		},
		{
			path: "assets/500_daemonset.yaml",
			readerAndApply: func(objBytes []byte) error {
				required := resourceread.ReadDaemonSetV1OrDie(objBytes)
				// override the operator image with the one built in the CI

				// E.g. IMAGE_FORMAT=registry.build03.ci.openshift.org/ci-op-52fj47p4/stable:${component}
				registry := strings.Split(os.Getenv("IMAGE_FORMAT"), "/")[0]

				required.Spec.Template.Spec.Containers[0].Image = registry + "/" + os.Getenv("NAMESPACE") + "/pipeline:run-once-duration-override-webhook"
				_, _, err := resourceapply.ApplyDaemonSet(
					ctx,
					kubeClient.AppsV1(),
					eventRecorder,
					required,
					1000, // any random high number
				)
				return err
			},
		},
		{
			path: "assets/600_mutating.yaml",
			readerAndApply: func(objBytes []byte) error {
				_, _, err := resourceapply.ApplyMutatingWebhookConfigurationImproved(ctx, kubeClient.AdmissionregistrationV1(), eventRecorder, resourceread.ReadMutatingWebhookConfigurationV1OrDie(objBytes), resourceapply.NewResourceCache())
				return err
			},
		},
	}

	// create required resources, e.g. namespace, crd, roles
	if err := wait.PollImmediate(1*time.Second, 10*time.Second, func() (bool, error) {
		for _, asset := range assets {
			klog.Infof("Creating %v", asset.path)
			if err := asset.readerAndApply(bindata.MustAsset(asset.path)); err != nil {
				klog.Errorf("Unable to create %v: %v", asset.path, err)
				return false, nil
			}
		}

		return true, nil
	}); err != nil {
		klog.Errorf("Unable to create SSO resources: %v", err)
		os.Exit(1)
	}

	webhooksRunning := make(map[string]bool) // node, exists
	// Get the number of master nodes
	if err := wait.PollImmediate(5*time.Second, 1*time.Minute, func() (bool, error) {
		klog.Infof("Listing nodes...")
		nodeItems, err := kubeClient.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
		if err != nil {
			klog.Errorf("Unable to list nodes: %v", err)
			return false, nil
		}
		for _, node := range nodeItems.Items {
			if _, exists := node.Labels["node-role.kubernetes.io/master"]; !exists {
				continue
			}
			webhooksRunning[node.Name] = false
		}
		return true, nil
	}); err != nil {
		klog.Errorf("Unable to collect master nodes")
		os.Exit(1)
	}

	// Wait until the secondary scheduler pod is running
	if err := wait.PollImmediate(5*time.Second, 1*time.Minute, func() (bool, error) {
		klog.Infof("Listing pods...")
		podItems, err := kubeClient.CoreV1().Pods(webhookNamespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			klog.Errorf("Unable to list pods: %v", err)
			return false, nil
		}
		for _, pod := range podItems.Items {
			if !strings.HasPrefix(pod.Name, webhookName+"-") {
				continue
			}
			klog.Infof("Checking pod: %v, phase: %v, deletionTS: %v\n", pod.Name, pod.Status.Phase, pod.GetDeletionTimestamp())
			if pod.Status.Phase == corev1.PodRunning && pod.GetDeletionTimestamp() == nil {
				webhooksRunning[pod.Spec.NodeName] = true
			}
		}

		for node, running := range webhooksRunning {
			if !running {
				klog.Infof("Admission scheduled on node %q is not running", node)
				return false, nil
			}
		}

		return true, nil
	}); err != nil {
		klog.Errorf("Unable to wait for all webhook pods running: %v", err)
		os.Exit(1)
	}

	klog.Infof("All daemonset webhooks are running")
	os.Exit(m.Run())
}

func TestMutation(t *testing.T) {}

func getKubeClientOrDie() *k8sclient.Clientset {
	kubeconfig := os.Getenv("KUBECONFIG")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		klog.Errorf("Unable to build config: %v", err)
		os.Exit(1)
	}
	client, err := k8sclient.NewForConfig(config)
	if err != nil {
		klog.Errorf("Unable to build client: %v", err)
		os.Exit(1)
	}
	return client
}
