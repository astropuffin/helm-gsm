package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

func main() {
	parseSecrets()
}

type secrets struct {
	Secrets map[string]string `yaml:"secrets,omitempty"`
}

func parseSecrets() {
	var raw secrets
	raw_values := raw.loadSecretYaml()
	if raw_values == nil {
		return
	}

	var plaintext secrets
	plaintext.Secrets = make(map[string]string)

	for k, v := range raw_values.Secrets {
		if valid, _ := checkValidSecret(v); valid {
			plaintext.Secrets[k] = base64.StdEncoding.EncodeToString(getSecret(transformStringToCanonicalName(v)))
		} else {
			plaintext.Secrets[k] = v
		}
	}

	data, err := yaml.Marshal(plaintext)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	f, err := os.Create("secrets.yaml.dec")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("secrets.yaml.dec", data, 0644)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()
}

func (s *secrets) loadSecretYaml() *secrets {
	secretsFile := "secrets.yaml"
	if fileExists(secretsFile) {
		yamlFile, err := ioutil.ReadFile(secretsFile)
		if err != nil {
			log.Printf("error loading secrets.yaml #%v", err)
		}
		err = yaml.Unmarshal(yamlFile, s)
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
		}
		return s
	} else {
		log.Printf("%s does not exist. Skipping decryption", secretsFile)
		return nil
	}
}

func transformStringToCanonicalName(raw string) string {
	s := strings.Split(raw[4:], "/")
	return fmt.Sprintf("projects/%s/secrets/%s/versions/%s", s[0], s[1], s[2])
}

func getSecret(s string) []byte {
	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to setup client: %v", err)
	}

	// Build the request.
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		//Name: "projects/dronedeploy-code-delivery-0/secrets/tiles_prod_common_secrets_json/versions/1",
		Name: s,
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		log.Fatalf("failed to access secret version: %v", err)
	}

	return result.Payload.Data
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

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
