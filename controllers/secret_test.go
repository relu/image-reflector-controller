package controllers

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	corev1 "k8s.io/api/core/v1"
)

func TestExtractAuthn(t *testing.T) {
	// the secret in testdata/secret.json was created with kubectl
	// create secret docker-registry. Test that it can be decoded to
	// get an authentication value.
	b, err := ioutil.ReadFile("testdata/secret.json")
	if err != nil {
		t.Fatal(err)
	}
	var secret corev1.Secret
	if err = json.Unmarshal(b, &secret); err != nil {
		t.Fatal(err)
	}
	auth, err := authFromSecret(secret, "https://index.docker.io/v1/")
	if err != nil {
		t.Fatal(err)
	}
	authConfig, err := auth.Authorization()
	if err != nil {
		t.Fatal()
	}
	if authConfig.Username != "fooser" || authConfig.Password != "foopass" {
		t.Errorf("expected username/password to be fooser/foopass, got %s/%s",
			authConfig.Username, authConfig.Password)
	}
}
