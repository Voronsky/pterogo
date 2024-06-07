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
	s, err := listServers(bearer_auth_token)
	if err != nil {
		log.Fatalf(`ListServers() = %q, %v, want nil, error`, s, err)
	}
	logger.Info("Servers queried", "Servers", s)
}
