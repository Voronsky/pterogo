package pterogo

import (
	"log"
	"log/slog"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestListServers(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	//s, err := listServers()
	err := godotenv.Load()
	if err != nil {
		log.Fatalf(`No env file found`)
	}
	bearer_auth_token := os.Getenv("PTERO_API_KEY")
	base_url := os.Getenv("BASE_URL")
	s, err := listServers(bearer_auth_token, base_url)
	if err != nil {
		log.Fatalf(`ListServers() = %q, %v, want nil, error`, s, err)
	}
	logger.Info("Servers queried", "Servers", s)
}

func TestListServersBadUrl_neg(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("ListServersBadUrl negative test")
	err := godotenv.Load()
	if err != nil {
		log.Fatalf(`No env file found`)
	}

	bearer_auth_token := os.Getenv("PTERO_API_KEY")
	s, err := listServers(bearer_auth_token, "https://example.com")
	if s != nil {
		log.Fatalf("Function returned a map, when it should have failed.")
	}
	if err == nil {
		logger.Info("Received an error with server variable set to nil.")
	}
	logger.Info("ListServersBadUrl negative test complete")

}

func TestListServersBadAuth_neg(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	err := godotenv.Load()
	logger.Info("ListServersBadAuth negative test")
	if err != nil {
		log.Fatalf(`No env file found`)
	}

	bearer_auth_token := "example"
	base_url := os.Getenv("BASE_URL")
	s, err := listServers(bearer_auth_token, base_url)
	if s != nil {
		log.Fatalf("Function returned a map, when it should have failed.")
	}
	logger.Info("Received an error with server variable set to nil.", "Error", err)
	logger.Info("ListServersBadAuth negative test complete")

}

func TestServerDetails(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	err := godotenv.Load()
	logger.Info("TestServerDetails()")
	if err != nil {
		log.Fatalf(`No env file found`)
	}

	// Get detail about the server passed
}
