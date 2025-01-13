package function

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"time"

	handle "github.com/my-ermes-labs/api-go/http"

	"github.com/my-ermes-labs/api-go/api"
	rc "github.com/my-ermes-labs/storage-redis/packages/go"
)

func cdn_upload(
	w http.ResponseWriter,
	r *http.Request,
	sessionToken api.SessionToken,
) (err error) {
	// Read the file name from the request query parameters.

	filename := r.URL.Query().Get("filename")
	if filename == "" {
		http.Error(w, "Missing 'filename' parameter", http.StatusBadRequest)
		return nil
	}

	// Read the file from the request body.
	file, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil || len(file) == 0 {
		http.Error(w, "Missing body in http message", http.StatusBadRequest)
		return nil
	}

	// Create a session if it does not exists and acquire it. Ermes
	// will handle the returned error with a “500 Internal Server Error”
	// response or a Redirect in case the error is an “ErrMigratedTo”
	// instance.
	_, err = Node.CreateAndAcquireSession(
		r.Context(),
		api.CreateAndAcquireSessionOptions{},
		func(sessionToken api.SessionToken) error {
			ks := rc.NewErmesKeySpaces(sessionToken.SessionId)
			// Derive the IO usage from the file size. Early returns and unset
			// resources will default to the average(!) usage of 1.
			Node.UpdateSessionResourcesUsage(r.Context(), sessionToken.SessionId, map[string]float64{
				"io": deriveIOUsage(len(file)),
			})

			// Set the file in the session.
			return redisClient.Set(
				r.Context(), ks.Session(filename), file, time.Hour).Err()
		})

	return err
}

func Handle(w http.ResponseWriter, r *http.Request) {
	// Check input...
	// Create options dynamically...
	// Do something...
	// Finally, Call the handler...
	myLog("Hello from the API handler!")
	handle.Handle(Node, w, r, options, cdn_upload)
}

func deriveIOUsage(size int) float64 {
	// Derive the IO usage from the file size. Early returns and unset
	// resources will default to the average(!) usage of 1.
	return 0.03 * (1 + float64(size)/1024/5)
}

// Create the handler options.
var options = handle.NewHandlerOptionsBuilder().
	AcquireSessionOptions(
		// Options to acquire a session.
		api.NewAcquireSessionOptionsBuilder().
			// Allow the session to be offloaded after it is acquired.
			AllowOffloading().
			// Allow the session to be acquired while it is offloading.
			AllowWhileOffloading().
			// Create the options.
			Build(),
	).
	CreateSessionOptions(
		// Options to create a session.
		api.NewCreateSessionOptionsBuilder().
			// Set the session to expire in 1 hour.
			Expires(3600).
			// Create the options.
			Build(),
	).
	// Set the session token header name.
	SessionTokenHeaderName("X-Ermes-Token").
	// Custom redirect logic.
	RedirectNewRequest(func(req *http.Request, node *api.Node) bool {
		// Access the node resources.
		sessions, _, _ := node.Cmd.GetNodeResourcesUsage(req.Context(), node.Host)
		// Redirect after reaching a threshold.
		return sessions > uint(node.Resources["max_sessions"])
	}).
	// Set the session token cookie name.
	Build()

func myLog(bodyContent string) (string, error) {
	url := "http://192.168.64.1:3000/handlerGo"

	requestBody := bytes.NewBufferString(bodyContent)

	req, err := http.NewRequest("POST", url, requestBody)
	if err != nil {
		return "", fmt.Errorf("error while creating the request: %v", err)
	}

	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error while sending the request: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error while reading the response: %v", err)
	}

	return string(responseBody), nil
}
