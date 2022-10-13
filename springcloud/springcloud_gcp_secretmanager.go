package springcloud

import (
	"context"
	"fmt"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

const (
	defaultName                                     = "projects/%s/secrets/%s/versions/%s"
	gcpSecretManagerShortVersionedLen               = 2
	gcpSecretManagerShortProjectsScopedVersionedLen = 3
	gcpSecretManagerLongProjectScopedLen            = 4
)

func getSecretManagerValue(ctx context.Context, client *secretmanager.Client, name string) (string, error) {
	finalName, err := parseName(ctx, name)
	if err != nil {
		return "", errors.Wrapf(err, "springcloud: failed parsing name %s", name)
	}

	// Build the request
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: finalName,
	}

	// Call the API
	val, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", errors.Wrapf(err, "springcloud: failed to access secret %s", name)
	}

	return string(val.Payload.Data), nil
}

func parseName(ctx context.Context, name string) (string, error) {
	tempName := strings.Split(name, "/")

	if strings.HasPrefix(name, "projects") && len(tempName) == gcpSecretManagerLongProjectScopedLen {
		return fmt.Sprintf("%s/%s/secrets/%s/versions/%s", tempName[0], tempName[1], tempName[2], tempName[3]), nil
	}

	// Get default credentials gcloud.
	credential, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		return "", errors.Wrapf(err, "springcloud: failed to access secret %s", name)
	}

	if len(tempName) == 1 {
		return fmt.Sprintf(defaultName, credential.ProjectID, tempName[0], "latest"), nil
	}

	if len(tempName) == gcpSecretManagerShortVersionedLen {
		return fmt.Sprintf(defaultName, credential.ProjectID, tempName[0], tempName[1]), nil
	}

	if len(tempName) == gcpSecretManagerShortProjectsScopedVersionedLen {
		return fmt.Sprintf(defaultName, tempName[0], tempName[1], tempName[2]), nil
	}

	return name, nil
}
