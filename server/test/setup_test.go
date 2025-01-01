package test

import (
	"fmt"
	"os"
	"testing"
)

var (
	tmpDatabaseDir string = ""
)

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		fmt.Printf("Failed to setup testing environment: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()

	err = cleanup()
	if err != nil {
		fmt.Printf("Failed to cleanup testing environment: %v\n", err)
		os.Exit(1)
	}

	os.Exit(code)
}

// Sets the testing enviroment.
// IMPORTANT: DO NOT USE THE OTHER SETUP METHODS. THIS SHOULD BE THE GOTO FUNCTION.
func setup() error {
	err := setupEnvironmentVariables()
	if err != nil {
		return err
	}

	return nil
}

func setupEnvironmentVariables() error {
	tmpDatabaseDir, err := os.MkdirTemp("", "libsql-")
	if err != nil {
		return err
	}
	os.Setenv("DATABASE_URL", "file:"+tmpDatabaseDir+"/test.db")

	os.Setenv("DATABASE_AUTH_TOKEN", "empty")

	os.Setenv("BACKEND_URL", "http://127.0.0.1:3000")

	os.Setenv("PORT", "3000")

	os.Setenv("APP_ENV", "testing")

	os.Setenv("RESEND_API_KEY", "key")
	os.Setenv("VERIFY_EMAIL_ADDRESS", "verify@mail.com")

	os.Setenv("USER_TOKEN_KEY", "usertokenkey")
	os.Setenv("BENTO_TOKEN_KEY", "bentotokenkey")
	os.Setenv("EMAIL_TOKEN_KEY", "emailtokenkey")

	return nil
}

// Cleans up the testing enviroment.
func cleanup() error {
	err := os.RemoveAll(tmpDatabaseDir)
	if err != nil {
		return err
	}

	return nil
}
