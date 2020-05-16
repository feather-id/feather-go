package feather_test

import (
	"log"

	"github.com/feather-id/feather-go"
)

func Example() {
	// Please note this example ignores errors for brevity.
	// You should not ignore errors in production code.

	// Initialize the client with your API key
	client := feather.New("YOUR_API_KEY")

	// Create an anonymous session
	session, _ := client.Sessions.Create(feather.SessionsCreateParams{})

	// Create a credential
	credential, _ := client.Credentials.Create(feather.CredentialsCreateParams{
		Type:     feather.CredentialTypeUsernamePassword,
		Username: feather.String("jdoe"),
		Password: feather.String("pa$$w0rd"),
	})

	// Inform the user of their credential status
	switch credential.Status {
	case feather.CredentialStatusRequiresOneTimeCode:
		log.Printf("Please check your email for a link to sign in")
		return

	case feather.CredentialStatusInvalid:
		log.Printf("Your username and password did not match")
		return

	case feather.CredentialStatusValid:
		// Life is good :)
		break
	}

	// Upgrade the session
	session, _ = client.Sessions.Upgrade(session.ID, feather.SessionsUpgradeParams{
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
