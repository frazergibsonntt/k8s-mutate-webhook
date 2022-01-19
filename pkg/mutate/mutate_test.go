package mutate

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	v1beta1 "k8s.io/api/admission/v1beta1"
)

func TestMutatesValidRequest(t *testing.T) {
	rawJSON := `{
		"kind": "AdmissionReview",
		"apiVersion": "admission.k8s.io/v1beta1",
		"request": {
			"uid": "7f0b2891-916f-4ed6-b7cd-27bff1815a8c",
			"kind": {
				"group": "",
				"version": "v1",
				"kind": "Pod"
			},
			"resource": {
				"group": "",
				"version": "v1",
				"resource": "pods"
			},
			"requestKind": {
				"group": "",
				"version": "v1",
				"kind": "Pod"
			},
			"requestResource": {
				"group": "",
				"version": "v1",
				"resource": "pods"
			},
			"namespace": "yolo",
			"operation": "CREATE",
			"userInfo": {
				"username": "kubernetes-admin",
				"groups": [
					"system:masters",
					"system:authenticated"
				]
			},
			"object": {
				"apiVersion": "networking.k8s.io/v1",
				"kind": "NetworkPolicy",
				"metadata": {
					"name": "default-netpol",
					"namespace": "mutatingwebhooktest",
				},
				"spec": {
					"egress": [
						{
							"ports": [
								{
									"port": 443,
									"protocol": "TCP"
								}
							],
							"to": [
								{
									"ipBlock": {
										"cidr": "10.105.66.53/32"
									}
								},
								{
									"ipBlock": {
										"cidr": "10.105.74.214/32"
									}
								}
							]
						}
					],
					"podSelector": {},
					"policyTypes": [
						"Egress"
					]
				}
			},
			"oldObject": null,
			"dryRun": false,
			"options": {
				"kind": "CreateOptions",
				"apiVersion": "meta.k8s.io/v1"
			}
		}
	}`
	response, err := Mutate([]byte(rawJSON), false)
	if err != nil {
		t.Errorf("failed to mutate AdmissionRequest %s with error %s", string(response), err)
	}

	r := v1beta1.AdmissionReview{}
	err = json.Unmarshal(response, &r)
	assert.NoError(t, err, "failed to unmarshal with error %s", err)

	rr := r.Response
	assert.Equal(t, `[{"op":"replace","path":"/spec/egress/0/to/0/ipBlock/CIDR","value":"10.105.66.53"}]`, string(rr.Patch))
	assert.Contains(t, rr.AuditAnnotations, "mutateme")

}

func TestErrorsOnInvalidJson(t *testing.T) {
	rawJSON := `Wut ?`
	_, err := Mutate([]byte(rawJSON), false)
	if err == nil {
		t.Error("did not fail when sending invalid json")
	}
}

func TestErrorsOnInvalidPod(t *testing.T) {
	rawJSON := `{
		"request": {
			"object": 111
		}
	}`
	_, err := Mutate([]byte(rawJSON), false)
	if err == nil {
		t.Error("did not fail when sending invalid pod")
	}
}
