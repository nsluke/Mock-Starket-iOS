package middleware

import (
	"context"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

// RealFirebaseVerifier verifies Firebase ID tokens using the Admin SDK.
type RealFirebaseVerifier struct {
	client *auth.Client
}

// NewFirebaseVerifier initializes a Firebase Auth verifier.
// projectID is the Firebase project ID.
// credentialsFile is optional — if empty, uses GOOGLE_APPLICATION_CREDENTIALS env var or default credentials.
func NewFirebaseVerifier(ctx context.Context, projectID string, credentialsFile string) (*RealFirebaseVerifier, error) {
	config := &firebase.Config{ProjectID: projectID}

	var app *firebase.App
	var err error

	if credentialsFile != "" {
		b, readErr := os.ReadFile(credentialsFile)
		if readErr != nil {
			return nil, fmt.Errorf("failed to read firebase credentials file: %w", readErr)
		}
		app, err = firebase.NewApp(ctx, config, option.WithCredentialsJSON(b))
	} else {
		app, err = firebase.NewApp(ctx, config)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase app: %w", err)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase auth client: %w", err)
	}

	return &RealFirebaseVerifier{client: client}, nil
}

// VerifyIDToken verifies a Firebase ID token and returns the user's UID.
func (v *RealFirebaseVerifier) VerifyIDToken(ctx context.Context, idToken string) (string, error) {
	token, err := v.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return "", fmt.Errorf("invalid Firebase token: %w", err)
	}
	return token.UID, nil
}
