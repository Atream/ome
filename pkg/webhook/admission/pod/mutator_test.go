package pod

import (
	"context"
	"encoding/json"
	"sort"
	"testing"

	"github.com/google/uuid"
	"github.com/onsi/gomega"
	gomegaTypes "github.com/onsi/gomega/types"
	"gomodules.xyz/jsonpatch/v2"
	"google.golang.org/protobuf/proto"
	admissionv1 "k8s.io/api/admission/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/sgl-project/ome/pkg/constants"
)

func TestMutator_Handle(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	omeNamespace := v1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: constants.OMENamespace,
		},
		Spec:   v1.NamespaceSpec{},
		Status: v1.NamespaceStatus{},
	}

	if err := c.Create(context.TODO(), &omeNamespace); err != nil {
		t.Errorf("failed to create namespace: %v", err)
	}
	mutator := Mutator{Client: c, Clientset: clientset, Decoder: admission.NewDecoder(c.Scheme())}

	cases := map[string]struct {
		configMap v1.ConfigMap
		request   admission.Request
		pod       v1.Pod
		matcher   gomegaTypes.GomegaMatcher
	}{
		"should not mutate non isvc pods": {
			configMap: v1.ConfigMap{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ConfigMap",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      constants.InferenceServiceConfigMapName,
					Namespace: constants.OMENamespace,
				},
				Immutable: nil,
				Data: map[string]string{
					constants.AgentConfigMapKeyName: `{
        				"image" : "ome/agent:latest",
        				"memoryRequest": "100Mi",
        				"memoryLimit": "1Gi",
        				"cpuRequest": "100m",
        				"cpuLimit": "1"
    				}`,
				},
				BinaryData: nil,
			},
			request: admission.Request{
				AdmissionRequest: admissionv1.AdmissionRequest{
					UID: types.UID(uuid.NewString()),
					Kind: metav1.GroupVersionKind{
						Group:   "",
						Version: "v1",
						Kind:    "Pod",
					},
					Resource: metav1.GroupVersionResource{
						Group:    "",
						Version:  "v1",
						Resource: "pods",
					},
					SubResource: "",
					RequestKind: &metav1.GroupVersionKind{
						Group:   "",
						Version: "v1",
						Kind:    "Pod",
					},
					RequestResource: &metav1.GroupVersionResource{
						Group:    "",
						Version:  "v1",
						Resource: "pods",
					},
					RequestSubResource: "",
					Name:               "",
					Namespace:          "default",
					Operation:          admissionv1.Create,
					Object:             runtime.RawExtension{},
					OldObject:          runtime.RawExtension{},
					DryRun:             nil,
					Options:            runtime.RawExtension{},
				},
			},
			pod: v1.Pod{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Pod",
					APIVersion: "v1",
				},
			},
			matcher: gomega.Equal(admission.Response{
				Patches: nil,
				AdmissionResponse: admissionv1.AdmissionResponse{
					UID:     "",
					Allowed: true,
					Result: &metav1.Status{
						TypeMeta: metav1.TypeMeta{},
						ListMeta: metav1.ListMeta{},
						Status:   "",
						Message:  "",
						Reason:   "",
						Details:  nil,
						Code:     200,
					},
					Patch:            nil,
					PatchType:        nil,
					AuditAnnotations: nil,
					Warnings:         nil,
				},
			}),
		},
		"should mutate isvc pods": {
			configMap: v1.ConfigMap{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ConfigMap",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      constants.InferenceServiceConfigMapName,
					Namespace: constants.OMENamespace,
				},
				Immutable: nil,
				Data: map[string]string{
					constants.AgentConfigMapKeyName: `{
        				"image" : "ome/agent:latest",
        				"memoryRequest": "100Mi",
        				"memoryLimit": "1Gi",
        				"cpuRequest": "100m",
        				"cpuLimit": "1"
    				}`,
				},
				BinaryData: nil,
			},
			request: admission.Request{
				AdmissionRequest: admissionv1.AdmissionRequest{
					UID: types.UID(uuid.NewString()),
					Kind: metav1.GroupVersionKind{
						Group:   "",
						Version: "v1",
						Kind:    "Pod",
					},
					Resource: metav1.GroupVersionResource{
						Group:    "",
						Version:  "v1",
						Resource: "pods",
					},
					SubResource: "",
					RequestKind: &metav1.GroupVersionKind{
						Group:   "",
						Version: "v1",
						Kind:    "Pod",
					},
					RequestResource: &metav1.GroupVersionResource{
						Group:    "",
						Version:  "v1",
						Resource: "pods",
					},
					RequestSubResource: "",
					Name:               "",
					Namespace:          "default",
					Operation:          admissionv1.Create,
					Object:             runtime.RawExtension{},
					OldObject:          runtime.RawExtension{},
					DryRun:             nil,
					Options:            runtime.RawExtension{},
				},
			},
			pod: v1.Pod{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Pod",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						constants.InferenceServicePodLabelKey: "isvc-pod",
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: constants.MainContainerName,
						},
					},
				},
			},
			matcher: gomega.BeEquivalentTo(admission.Response{
				Patches: []jsonpatch.JsonPatchOperation{
					{
						Operation: "add",
						Path:      "/metadata/annotations",
						Value: map[string]interface{}{
							"ome.io/enable-metric-aggregation":  "",
							"ome.io/enable-prometheus-scraping": "",
						},
					},
					{
						Operation: "add",
						Path:      "/metadata/namespace",
						Value:     "default",
					},
				},
				AdmissionResponse: admissionv1.AdmissionResponse{
					UID:              "",
					Allowed:          true,
					Result:           nil,
					Patch:            nil,
					PatchType:        (*admissionv1.PatchType)(proto.String(string(admissionv1.PatchTypeJSONPatch))),
					AuditAnnotations: nil,
					Warnings:         nil,
				},
			}),
		},
		"should mutate training pods": {
			configMap: v1.ConfigMap{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ConfigMap",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      constants.InferenceServiceConfigMapName,
					Namespace: constants.OMENamespace,
				},
				Immutable: nil,
				Data: map[string]string{
					constants.AgentConfigMapKeyName: `{
        				"image" : "ome/agent:latest",
        				"memoryRequest": "100Mi",
        				"memoryLimit": "1Gi",
        				"cpuRequest": "100m",
        				"cpuLimit": "1"
    				}`,
				},
				BinaryData: nil,
			},
			request: admission.Request{
				AdmissionRequest: admissionv1.AdmissionRequest{
					UID: types.UID(uuid.NewString()),
					Kind: metav1.GroupVersionKind{
						Group:   "",
						Version: "v1",
						Kind:    "Pod",
					},
					Resource: metav1.GroupVersionResource{
						Group:    "",
						Version:  "v1",
						Resource: "pods",
					},
					SubResource: "",
					RequestKind: &metav1.GroupVersionKind{
						Group:   "",
						Version: "v1",
						Kind:    "Pod",
					},
					RequestResource: &metav1.GroupVersionResource{
						Group:    "",
						Version:  "v1",
						Resource: "pods",
					},
					RequestSubResource: "",
					Name:               "",
					Namespace:          "default",
					Operation:          admissionv1.Create,
					Object:             runtime.RawExtension{},
					OldObject:          runtime.RawExtension{},
					DryRun:             nil,
					Options:            runtime.RawExtension{},
				},
			},
			pod: v1.Pod{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Pod",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						constants.TrainingJobPodLabelKey: "trainingjob",
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: "trainer",
						},
					},
				},
			},
			matcher: gomega.BeEquivalentTo(admission.Response{
				Patches: []jsonpatch.JsonPatchOperation{
					{
						Operation: "add",
						Path:      "/metadata/annotations",
						Value: map[string]interface{}{
							"ome.io/enable-metric-aggregation":  "",
							"ome.io/enable-prometheus-scraping": "",
						},
					},
					{
						Operation: "add",
						Path:      "/metadata/namespace",
						Value:     "default",
					},
				},
				AdmissionResponse: admissionv1.AdmissionResponse{
					UID:              "",
					Allowed:          true,
					Result:           nil,
					Patch:            nil,
					PatchType:        (*admissionv1.PatchType)(proto.String(string(admissionv1.PatchTypeJSONPatch))),
					AuditAnnotations: nil,
					Warnings:         nil,
				},
			}),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if err := c.Create(context.TODO(), &tc.configMap); err != nil {
				t.Errorf("failed to create config map: %v", err)
			}
			byteData, err := json.Marshal(tc.pod)
			if err != nil {
				t.Errorf("failed to marshal pod data: %v", err)
			}
			tc.request.Object.Raw = byteData
			res := mutator.Handle(context.TODO(), tc.request)
			sortPatches(res.Patches)
			g.Expect(res).Should(tc.matcher)
			if err := c.Delete(context.TODO(), &tc.configMap); err != nil {
				t.Errorf("failed to delete configmap %v", err)
			}
		})
	}
}

// sortPatches sorts the slice of patches by Path so that the comparison works
// when there are > 1 patches. Note: make sure the matcher Patches are sorted.
func sortPatches(patches []jsonpatch.JsonPatchOperation) {
	if len(patches) > 1 {
		sort.Slice(patches, func(i, j int) bool {
			return patches[i].Path < patches[j].Path
		})
	}
}
