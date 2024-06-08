package pterogo

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
)

// PteroRequestHeaders keeps track of the auth token and base url for all requests
// Its methods allow to make a request using the auth token and base url
type PteroRequestHeaders struct {
	auth_token string
	url        string
}

// A PterodactylClient implements methods for all client API routes
type PterodactylClient struct {
	Request PteroRequestHeaders //underlying PteroRequestHeaders needed for client Requests
}

// PteroResp will hold the JSON decoded body sent from the Pterodactyl Server.
type PteroResp struct {
	Object string      `json:"object"`
	Data   []PteroData `json:"data"`
}

// PteroData holds all the JSON decoded the nested Pterodactyl 'object' and data found in the Response
type PteroData struct {
	Object     string     `json:"object"`
	Attributes Attributes `json:"attributes"`
}

// Attributes holds the attributes of the Pterodactly Object found in the Data JSON object
type Attributes struct {
	Name        string `json:"name"`
	Identifier  string `json:"identifier"`
	Description string `json:"description"`
}

// Holds necessary information about a server
type Server struct {
	Name        string
	Description string
}

// Builds the custom headers needed for Pterodactyl API routes
// Executes the Request based on the method and route passed
func (prh PteroRequestHeaders) PteroRequest(method string, route string) ([]byte, error) {
	client := &http.Client{}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	//Build Method Request
	req, err := http.NewRequest(method, route, nil)
	if err != nil {
		slog.Error("Failed to make a new request", "Error", err)
		return nil, err
	}

	//Add Pterodactyl Headers
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+prh.auth_token)

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
		err := fmt.Errorf("received internal server error=%d, please report this trace to the github", resp.StatusCode)
		logger.Error("Internal server code was returned", "StatusCode", resp.StatusCode)
		return nil, err
	}

	logger.Info("Request successful", "Resp", resp.StatusCode)

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Failed to read body", "Error", err)
		return nil, err
	}
	return body, nil
}

// Grabs the list of servers from Pterodactyl
// Taken from the Pterodactyl API page. This will return an error if it fails at any point
// Otherwise, it will return a map of unique servers , based off their identifier
// A Bearer Auth token is required
func (pc PterodactylClient) ListServers() (map[string]Server, error) {
	r := PteroResp{}
	servers := map[string]Server{}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	//Build GET Request
	route := fmt.Sprintf("%s/api/client", pc.Request.url)

	// Decode the JSON body into the appropriate interface
	body, err := pc.Request.PteroRequest("GET", route)
	if err != nil {
		logger.Error("Error received making request to Pterodactyl", "Error", err)
		return nil, err
	}

	json.Unmarshal(body, &r)
	slog.Info("Decoded JSON body", "pteroResp=", r)
	for i := 0; i < len(r.Data); i++ {
		attrs := r.Data[i]
		servers[attrs.Attributes.Identifier] = Server{attrs.Attributes.Name, attrs.Attributes.Description}
		logger.Info("Server identifer", "Server", attrs.Attributes.Identifier)
	}

	return servers, nil
}

// Return server details for the specific identifier
func (pc PterodactylClient) ServerDetails(identifier string) (*Server, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	server := &Server{}
	data := PteroData{}

	//Build GET route and make Request
	route := fmt.Sprintf("%s/api/client/servers/%s", pc.Request.url, identifier)

	body, err := pc.Request.PteroRequest("GET", route)
	if err != nil {
		logger.Error("Error received making request to Pterodactyl", "Error", err)
		return nil, err
	}

	json.Unmarshal(body, &data)
	slog.Info("Decoded JSON body", "PteroResp", data)

	server.Name = data.Attributes.Name
	server.Description = data.Attributes.Description

	return server, nil
}
