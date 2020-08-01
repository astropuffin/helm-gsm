// Sample quickstart is a basic program that uses Secret Manager.
package main

import (
	"context"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"regexp"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

func main() {
	//valid, err := checkValidSecret("gsm:project-id/secret_name/0")
	//if !valid {
	//	log.Printf("%s", err)
	//}
	var s secrets
	values := s.loadSecretYaml()

	log.Printf("%s", values.Secrets["a"])
}

type secrets struct {
	Secrets map[string]string `yaml:"secrets,omitempty"`
}

func (s *secrets) loadSecretYaml() *secrets {
	yamlFile, err := ioutil.ReadFile("secrets.yaml")
	if err != nil {
		log.Printf("error loading secrets.yaml #%v", err)
	}
	err = yaml.Unmarshal(yamlFile, s)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return s
}

func getSecret(s string) {
	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to setup client: %v", err)
	}

	// Build the request.
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: "projects/dronedeploy-code-delivery-0/secrets/tiles_prod_common_secrets_json/versions/1",
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		log.Fatalf("failed to access secret version: %v", err)
	}

	// Print the secret payload.
	//
	// WARNING: Do not print the secret in a production environment - this
	// snippet is showing how to access the secret material.
	log.Printf("Plaintext: %s", result.Payload.Data)
}

func checkValidSecret(s string) (bool, string) {
	regex_string := `^gsm:[a-z][a-z0-9-]{4,28}[a-z0-9]\/[a-zA-Z0-9_-]+\/[1-9][0-9]*$`
	matched, err := regexp.Match(regex_string, []byte(s))
	if err != nil {
		log.Fatalf("error matching regex")
	}
	msg := ""
	if !matched {
		msg = "Secret did not match required format.\n" +
			"\n" +
			"Must be in the form 'gsm:project_id/secret_name/version'.\n" +
			"project_id: 'The unique, user-assigned ID of the Project. It must be 6 to 30 lowercase letters, digits, or hyphens. It must start with a letter. Trailing hyphens are prohibited.'\n" +
			"secret_name: 'Secret names can only contain English letters (A-Z), numbers (0-9), dashes (-), and underscores (_)'\n" +
			"version: 'Versions are a monotonically increasing integer starting at 1.'\n" +
			"\n" +
			"example: 'gsm:project-id/secret_name/1'\n" +
			"regex: '" + regex_string + "'\n" +
			"\n"
	}

	return matched, msg
}
