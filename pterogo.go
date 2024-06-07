package pterogo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
)

type PteroResp struct {
	Object string      `json:"object, omitempty"`
	Data   []PteroData `json:"data, omitempty"`
}

type PteroData struct {
	Object     string     `json:"object, omitempty"`
	Attributes Attributes `json:"attributes, omitempty"`
}

type Attributes struct {
	Name        string `json:"name, omitempty"`
	Identifier  string `json:"identifier, omitempty"`
	Description string `json:"description, omitempty"`
}

type Server struct {
	Name        string
	Description string
}

// Grabs the list of servers from Pterodactyl
// Taken from the Pterodactyl API page. This will return an error if it fails at any point
// Otherwise, it will return a map of unique servers , based off their identifier
// A Bearer Auth token is required
func listServers(auth_token string, url string) (map[string]Server, error) {
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
		slog.Error("An error occurred trying to issue the request", "Error", err)
		return nil, err
	}

	if resp.StatusCode >= 300 {
		// Create custom error for this
		err := errors.New("request failed")
		slog.Error("Non-200 status code was returned", "StatusCode", resp.StatusCode)
		return nil, err
	}

	slog.Info("Request successful", "Resp", resp)

	servers := map[string]Server{}
	r := PteroResp{}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Failed to read body", "Error", err)
		return nil, err
	}

	// Decode the JSON body into the appropriate interface
	json.Unmarshal(body, &r)
	slog.Info("listServers()", "pteroResp", r)
	for i := 0; i < len(r.Data); i++ {
		attrs := r.Data[i]
		servers[attrs.Attributes.Identifier] = Server{attrs.Attributes.Name, attrs.Attributes.Description}
		//slog.Info("Server identifier", "Server=", attrs.Attributes.Identifier)
		logger.Info("Server identifer", "Server=", attrs.Attributes.Identifier)
	}

	return servers, nil
}
