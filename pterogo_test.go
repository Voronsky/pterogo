package pterogo

import (
	"log"
	"log/slog"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestListServers(t *testing.T) {

	// Parse the env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf(`No env file found`)
	}
	bearer_auth_token := os.Getenv("PTERO_API_KEY")
	base_url := os.Getenv("BASE_URL")

	// Test method
	client := PterodactylClient{
		Request: PteroRequestHeaders{bearer_auth_token, base_url},
	}

	s, err := client.ListServers()
	if err != nil {
		log.Fatalf(`ListServers() = %q, %v, want nil, error`, s, err)
	}
	logger.Info("Servers queried", "Servers", s)
}

func TestListServersBadUrl_neg(t *testing.T) {
	logger.Info("=== ListServersBadUrl negative test ===")

	// Parse env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf(`No env file found`)
	}

	bearer_auth_token := os.Getenv("PTERO_API_KEY")

	// Test the bad route
	client := PterodactylClient{
		Request: PteroRequestHeaders{bearer_auth_token, "https://example.com"},
	}

	s, err := client.ListServers()
	if s != nil {
		log.Fatalf("Function returned a map, when it should have failed.")
	}
	if err == nil {
		logger.Info("Received an error with server variable set to nil.")
	}
	logger.Info("=== ListServersBadUrl negative test complete ===")

}

func TestListServersBadAuth_neg(t *testing.T) {
	logger.Info("=== ListServersBadAuth negative test ===")
	err := godotenv.Load()
	if err != nil {
		log.Fatalf(`No env file found`)
	}

	// Test with bad auth
	bearer_auth_token := "example"
	base_url := os.Getenv("BASE_URL")

	client := PterodactylClient{
		Request: PteroRequestHeaders{bearer_auth_token, base_url},
	}
	s, err := client.ListServers()
	if s != nil {
		log.Fatalf("Function returned a map, when it should have failed.")
	}
	logger.Info("Received an error with server variable set to nil.", "Error", err)
	logger.Info("=== ListServersBadAuth negative test complete ===")

}

func TestServerDetails(t *testing.T) {
	logger.Info("=== TestServerDetails() begin === ")

	// Parse env file
	err := godotenv.Load()
	bearer_auth_token := os.Getenv("PTERO_API_KEY")
	base_url := os.Getenv("BASE_URL")
	if err != nil {
		log.Fatalf(`No env file found`)
	}

	client := PterodactylClient{
		Request: PteroRequestHeaders{bearer_auth_token, base_url},
	}

	s, err := client.ServerDetails("102248be")
	if err != nil {
		log.Fatalf(`Error retrieving server details, wanted non-nil error`)
	}

	if s.Name == "" && s.Description == "" {
		log.Fatalf(`Pterodactly Response returned an empty response, wanted server name and desc`)
	}

	// Get detail about the server passed

	logger.Info("Server info received", "Server Info", s)
	logger.Info(" === TestServerDetails() complete === ")

}

func TestChangePowerState(t *testing.T) {
	logger.Info(" === TestChangePowerState() begin ===")

	// Parse env file
	err := godotenv.Load()
	bearer_auth_token := os.Getenv("PTERO_API_KEY")
	base_url := os.Getenv("BASE_URL")

	if err != nil {
		log.Fatalf(`No env file found`)
	}

	client := PterodactylClient{
		Request: PteroRequestHeaders{bearer_auth_token, base_url},
	}

	success, err := client.ChangePowerState("102248be", "start")
	if err != nil {
		log.Fatalf("Error trying to change power state")
	}

	if success != 0 {
		logger.Error("Error in trying to change the power state of the server", "SuccessCode", success)
		log.Fatalf("Error in trying to change the power state of the server")
	}

	logger.Info("Change State succeeded", "SuccessCode", success)
	logger.Info("=== TestChangePowerState() complete ===")
}

func TestGetPowerState(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("== TestGetPowerState() begin ==")
	identifier := "102248be"

	// Parse env file
	err := godotenv.Load()
	bearer_auth_token := os.Getenv("PTERO_API_KEY")
	base_url := os.Getenv("BASE_URL")

	if err != nil {
		log.Fatalf(`No env file found`)
	}

	client := PterodactylClient{
		Request: PteroRequestHeaders{bearer_auth_token, base_url},
	}

	logger.Info("Getting power state for server", "Identifier", identifier)
	state, err := client.GetPowerState(identifier)
	if err != nil {
		log.Fatalf("Error trying to get power state")
	}

	if state == "" {
		logger.Error("State unable to be retrieved", "State", state)
		log.Fatalf("State was empty string, expected non-empty string")
	}

	logger.Info("Succeeded retrieving power state", "PowerState", state)
	logger.Info("== TestGetPowerState() complete ==")
}
