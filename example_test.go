package feather_test

import (
	"log"

	"github.com/feather-id/feather-go"
)

func Example() {
	// Please note this example ignores errors for brevity.
	// You should not ignore errors in production code.

	// Initialize the client with your API key
	client := feather.New("live_...")

	// Create a credential
	credential, _ := client.Credentials.Create(feather.CredentialsCreateParams{
		Type:     feather.CredentialTypeUsernamePassword,
		Username: feather.String("jdoe"),
		Password: feather.String("pa$$w0rd"),
	})

	// Inform the user of their credential status
	switch credential.Status {
	case feather.CredentialStatusRequiresVerificationCode:
		log.Printf("Please check your email for a link to sign in")
		return

	case feather.CredentialStatusInvalid:
		log.Printf("Your username and password did not match")
		return

	case feather.CredentialStatusValid:
		// The username and password were valid!
		break
	}

	// Create an authenticated session
	session, _ := client.Sessions.Create(feather.SessionsCreateParams{
		CredentialToken: credential.Token,
	})

	// Add custom metadata to the user
	user, _ := client.Users.Update(session.UserID, feather.UsersUpdateParams{
		Metadata: &map[string]string{
			"highScore": "123",
		},
	})

	log.Printf("Your high score is: %v", user.Metadata["highScore"])
}
