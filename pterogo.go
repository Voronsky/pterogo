package pterogo

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
)

// Collection of methods that pertain to the Pterodactyl Client API
type ClientAPI interface {
	ListServers() (map[string]Server, error)
	ServerDetails() (Server, error)
}
type PterodactylClient struct{}

type PteroResp struct {
	Object string      `json:"object"`
	Data   []PteroData `json:"data"`
}

type PteroData struct {
	Object     string     `json:"object"`
	Attributes Attributes `json:"attributes"`
}

type Attributes struct {
	Name        string `json:"name"`
	Identifier  string `json:"identifier"`
	Description string `json:"description"`
}

type Server struct {
	Name        string
	Description string
}

// Grabs the list of servers from Pterodactyl
// Taken from the Pterodactyl API page. This will return an error if it fails at any point
// Otherwise, it will return a map of unique servers , based off their identifier
// A Bearer Auth token is required
func (pc PterodactylClient) ListServers(auth_token string, url string) (map[string]Server, error) {
	client := &http.Client{}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	//Build GET Request
	route := fmt.Sprintf("%s/api/client", url)
	req, err := http.NewRequest("GET", route, nil)
	if err != nil {
		slog.Error("Failed to make a new request", "Error", err)
		return nil, err
	}

	//Add Pterodactyl Headers
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+auth_token)

	//Issue GET request
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("An error occurred trying to issue the request", "Error", err)
		return nil, err
	}

	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		// Create custom error for this
		err := fmt.Errorf("received redirection error=%d", resp.StatusCode)
		logger.Error("Redirect error code was returned", "StatusCode", resp.StatusCode)
		return nil, err
	}

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		// Create custom error for this
		err := fmt.Errorf("received client error=%d", resp.StatusCode)
		logger.Error("Client error code was returned", "StatusCode", resp.StatusCode)
		return nil, err
	}

	if resp.StatusCode >= 500 {
		err := fmt.Errorf("received internal server error=%d, please report this to the github", resp.StatusCode)
		logger.Error("Internal server code was returned", "StatusCode", resp.StatusCode)
		return nil, err
	}

	logger.Info("Request successful", "Resp", resp.StatusCode)

	servers := map[string]Server{}
	r := PteroResp{}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Failed to read body", "Error", err)
		return nil, err
	}

	// Decode the JSON body into the appropriate interface
	json.Unmarshal(body, &r)
	slog.Info("Decoded JSON body", "pteroResp=", r)
	for i := 0; i < len(r.Data); i++ {
		attrs := r.Data[i]
		servers[attrs.Attributes.Identifier] = Server{attrs.Attributes.Name, attrs.Attributes.Description}
		//slog.Info("Server identifier", "Server=", attrs.Attributes.Identifier)
		logger.Info("Server identifer", "Server=", attrs.Attributes.Identifier)
	}

	return servers, nil
}

func (pc PterodactylClient) ServerDetails(identifier string, auth_token string, url string) (*Server, error) {
	client := &http.Client{}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	//Build GET Request
	route := fmt.Sprintf("%s/api/client/servers/%s", url, identifier)
	req, err := http.NewRequest("GET", route, nil)
	server := &Server{}
	data := PteroData{}

	if err != nil {
		slog.Error("Failed to make a new request", "Error", err)
		return nil, err
	}

	//Add Pterodactyl Headers
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+auth_token)

	//Issue GET request
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("An error occurred trying to issue the request", "Error", err)
		return nil, err
	}

	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		// Create custom error for this
		err := fmt.Errorf("received redirection error=%d", resp.StatusCode)
		logger.Error("Redirect error code was returned", "StatusCode", resp.StatusCode)
		return nil, err
	}

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		// Create custom error for this
		err := fmt.Errorf("received client error=%d", resp.StatusCode)
		logger.Error("Client error code was returned", "StatusCode", resp.StatusCode)
		return nil, err
	}

	if resp.StatusCode >= 500 {
		err := fmt.Errorf("received internal server error=%d, please report this to the github", resp.StatusCode)
		logger.Error("Internal server code was returned", "StatusCode", resp.StatusCode)
		return nil, err
	}

	logger.Info("Request successful", "Resp", resp.Status)

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		logger.Error("Failed to read body", "Error", err)
		return nil, err
	}

	json.Unmarshal(body, &data)
	slog.Info("Decoded JSON body", "pteroResp=", data)

	server.Name = data.Attributes.Name
	server.Description = data.Attributes.Description

	return server, nil
}
