// ------------------------------------------------------------
// 📁 File: pkg/token/google.go
// 🧠 Handles Google ID token validation and extraction.
//     This file verifies the token against Google's servers
//     and parses the payload to extract user identity.

package token

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// 📦 Payload returned by Google's tokeninfo endpoint
type GoogleIDTokenPayload struct {
	Sub     string `json:"sub"`     // Google user ID
	Email   string `json:"email"`   // User email
	Name    string `json:"name"`    // Full name
	Picture string `json:"picture"` // Profile picture URL
}

// func VerifyGoogleIDToken(idToken string) (*GoogleIDTokenPayload, error) {
// 	return &GoogleIDTokenPayload{
// 		Sub:     "dev-google-id-123",
// 		Email:   "testuser@example.com",
// 		Name:    "Test User",
// 		Picture: "https://picsum.photos/200",
// 	}, nil
// }

// func VerifyGoogleIDToken(idToken string) (*GoogleIDTokenPayload, error) {
// 	switch idToken {
// 	case "user-1-token":
// 		return &GoogleIDTokenPayload{Sub: "google-1", Email: "user1@test.com", Name: "User One", Picture: "https://picsum.photos/id/1/200"}, nil
// 	case "user-2-token":
// 		return &GoogleIDTokenPayload{Sub: "google-2", Email: "user2@test.com", Name: "User Two", Picture: "https://picsum.photos/id/2/200"}, nil
// 	case "user-3-token":
// 		return &GoogleIDTokenPayload{Sub: "google-3", Email: "user3@test.com", Name: "User Three", Picture: "https://picsum.photos/id/3/200"}, nil
// 	case "user-4-token":
// 		return &GoogleIDTokenPayload{Sub: "google-4", Email: "user4@test.com", Name: "User Four", Picture: "https://picsum.photos/id/4/200"}, nil
// 	case "user-5-token":
// 		return &GoogleIDTokenPayload{Sub: "google-5", Email: "user5@test.com", Name: "User Five", Picture: "https://picsum.photos/id/5/200"}, nil
// 	case "user-6-token":
// 		return &GoogleIDTokenPayload{Sub: "google-6", Email: "user6@test.com", Name: "User Six", Picture: "https://picsum.photos/id/6/200"}, nil
// 	case "user-7-token":
// 		return &GoogleIDTokenPayload{Sub: "google-7", Email: "user7@test.com", Name: "User Seven", Picture: "https://picsum.photos/id/7/200"}, nil
// 	case "user-8-token":
// 		return &GoogleIDTokenPayload{Sub: "google-8", Email: "user8@test.com", Name: "User Eight", Picture: "https://picsum.photos/id/8/200"}, nil
// 	case "user-9-token":
// 		return &GoogleIDTokenPayload{Sub: "google-9", Email: "user9@test.com", Name: "User Nine", Picture: "https://picsum.photos/id/9/200"}, nil
// 	case "user-10-token":
// 		return &GoogleIDTokenPayload{Sub: "google-10", Email: "user10@test.com", Name: "User Ten", Picture: "https://picsum.photos/id/10/200"}, nil
// 	case "user-11-token":
// 		return &GoogleIDTokenPayload{Sub: "google-11", Email: "user11@test.com", Name: "User Eleven", Picture: "https://picsum.photos/id/11/200"}, nil
// 	case "user-12-token":
// 		return &GoogleIDTokenPayload{Sub: "google-12", Email: "user12@test.com", Name: "User Twelve", Picture: "https://picsum.photos/id/12/200"}, nil
// 	case "user-13-token":
// 		return &GoogleIDTokenPayload{Sub: "google-13", Email: "user13@test.com", Name: "User Thirteen", Picture: "https://picsum.photos/id/13/200"}, nil
// 	case "user-14-token":
// 		return &GoogleIDTokenPayload{Sub: "google-14", Email: "user14@test.com", Name: "User Fourteen", Picture: "https://picsum.photos/id/14/200"}, nil
// 	case "user-15-token":
// 		return &GoogleIDTokenPayload{Sub: "google-15", Email: "user15@test.com", Name: "User Fifteen", Picture: "https://picsum.photos/id/15/200"}, nil
// 	case "user-16-token":
// 		return &GoogleIDTokenPayload{Sub: "google-16", Email: "user16@test.com", Name: "User Sixteen", Picture: "https://picsum.photos/id/16/200"}, nil
// 	case "user-17-token":
// 		return &GoogleIDTokenPayload{Sub: "google-17", Email: "user17@test.com", Name: "User Seventeen", Picture: "https://picsum.photos/id/17/200"}, nil
// 	case "user-18-token":
// 		return &GoogleIDTokenPayload{Sub: "google-18", Email: "user18@test.com", Name: "User Eighteen", Picture: "https://picsum.photos/id/18/200"}, nil
// 	case "user-19-token":
// 		return &GoogleIDTokenPayload{Sub: "google-19", Email: "user19@test.com", Name: "User Nineteen", Picture: "https://picsum.photos/id/19/200"}, nil
// 	case "user-20-token":
// 		return &GoogleIDTokenPayload{Sub: "google-20", Email: "user20@test.com", Name: "User Twenty", Picture: "https://picsum.photos/id/20/200"}, nil
// 	default:
// 		return &GoogleIDTokenPayload{Sub: "google-default", Email: "default@test.com", Name: "Default User", Picture: "https://picsum.photos/200"}, nil
// 	}
// }

// 🔐 Verifies Google ID token using tokeninfo endpoint
func VerifyGoogleIDToken(idToken string) (*GoogleIDTokenPayload, error) {
	resp, err := http.Get(fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", idToken))
	if err != nil {
		return nil, fmt.Errorf("failed to contact Google: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid Google ID token: %s", resp.Status)
	}

	var payload GoogleIDTokenPayload
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("failed to decode token payload: %w", err)
	}

	// Basic checks (expand if needed)
	if payload.Sub == "" || payload.Email == "" {
		return nil, fmt.Errorf("invalid Google token payload")
	}

	return &payload, nil
}
