based on https://github.com/jkroepke/helm-secrets, but I needed Google Secret Manager instead of sops, and I didn't want to install sops if I didn't need it.

secret format:
```
Secret format
gsm:project_id/secret_name/version

Regex:
^gsm:[a-z][a-z0-9-]{4,28}[a-z0-9]\/[a-zA-Z0-9_-]+\/[1-9]?[0-9]+$
```

```
From Google's Documentation:
project_id: "The unique, user-assigned ID of the Project. It must be 6 to 30 lowercase letters, digits, or hyphens. It must start with a letter. Trailing hyphens are prohibited." 
secret_name: "Secret names can only contain English letters (A-Z), numbers (0-9), dashes (-), and underscores (_)"
version: Versions are a monotonically increasing integer starting at 1.
```
