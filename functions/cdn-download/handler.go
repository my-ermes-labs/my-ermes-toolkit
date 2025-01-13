package function

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	handle "github.com/my-ermes-labs/api-go/http"

	"github.com/my-ermes-labs/api-go/api"
	rc "github.com/my-ermes-labs/storage-redis/packages/go"
	"github.com/sashabaranov/go-openai"
)

var client = openai.NewClient(os.Getenv("OPENAI_API_KEY"))

func cnd_download(
	w http.ResponseWriter,
	r *http.Request,
	sessionToken api.SessionToken,
) (err error) {
	// Check that a session exists.
	// if sessionToken == nil {
	// 	http.Error(w, "No session available", http.StatusBadRequest)
	// 	return nil
	// }

	// Read the file name from the request query parameters.
	filename := r.URL.Query().Get("filename")
	if filename == "" {
		http.Error(w, "Missing ‘filename’ parameter", http.StatusBadRequest)
		return nil
	}

	// Acquire the session. Ermes will handle the returned error with a
	// “500 Internal Server Error” response or a Redirect in case the error
	// is an “ErrMigratedTo” instance.
	_, err = Node.AcquireSession(
		r.Context(),
		sessionToken,
		api.NewAcquireSessionOptionsBuilder().AllowOffloading().Build(),
		func() error {
			ks := rc.NewErmesKeySpaces(sessionToken.SessionId)
			// Get the file from the session.
			file, err := redisClient.Get(
				r.Context(), ks.Session(filename)).Result()
			if err != nil {
				return err
			}

			if file == "" {
				// Return an error if the file is not found.
				http.Error(w, "File not found", http.StatusNotFound)
			} else {
				// Derive the IO usage from the file size. Early returns and
				// unset resources will default to the average(!) usage of 1.
				Node.UpdateSessionResourcesUsage(r.Context(), sessionToken.SessionId, map[string]float64{
					"io": deriveIOUsage(len(file)),
				})

				// Write the file in the response.
				w.Header().Set("Content-Type", "application/octet-stream")
				w.Write([]byte(file))
			}

			return nil
		})
	return err
}

func deriveIOUsage(size int) float64 {
	// Derive the IO usage from the file size. Early returns and unset
	// resources will default to the average(!) usage of 1.
	return 0.03 * (1 + float64(size)/1024/5)
}

func Handle(w http.ResponseWriter, r *http.Request) {
	// Check input...
	// Create options dynamically...
	// Do something...
	// Finally, Call the handler...
	myLog("Hello from the API handler!")
	handle.Handle(Node, w, r, options, cnd_download)
}

func Ask(prompt string) string {
	// TODO: This part is commented out to remove the little uncertainty given by
	// The use of an external API. Initial tests are done in a more predictable way
	// To ease the comparison between the different setups.
	/*ctx := context.Background()

	req := openai.CompletionRequest{
		Model:     openai.GPT3Ada,
		MaxTokens: 50,
		Prompt:    prompt,
	}
	resp, err := client.CreateCompletion(ctx, req)
	if err != nil {
		fmt.Printf("Completion error: %v\n", err)
		return ""
	}
	return resp.Choices[0].Text
	*/

	start := time.Now()
	for time.Since(start) < time.Duration((50+int(len(prompt)/1024/100)))*time.Millisecond {
		// Keep the CPU busy
	}
	response := prompt
	return response
}

func deriveCPUUsage(size int) float64 {
	// Derive the CPU usage from the response size. Early returns and unset
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
