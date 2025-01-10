package function

import (
	"bytes"
	"fmt"
	"io"
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

func speech_to_text(
	w http.ResponseWriter,
	r *http.Request,
	sessionToken api.SessionToken,
) (err error) {
	myLog("SPEECH TO TEXT")
	// Read the file from the request body.
	fileWav := r.Body
	defer r.Body.Close()
	if err != nil {
		myLog("\n\nError!\n\n" + err.Error())
		return err
	}

	myLog("Filewav found \n\n")

	ss, err := Node.CreateAndAcquireSession(
		r.Context(),
		api.CreateAndAcquireSessionOptions{},
		func(sessionToken api.SessionToken) error {
			ks := rc.NewErmesKeySpaces(sessionToken.SessionId)
			// Code to handle speech to text and AIChat.
			prompt := SpeechToText(fileWav)
			aiResponse := Ask(prompt)

			// Derive the CPU usage from the response size. Early returns and
			// unset resources will default to the average(!) usage of 1.
			Node.UpdateSessionResourcesUsage(r.Context(), sessionToken.SessionId, map[string]float64{
				"cpu": deriveCPUUsage(len(aiResponse)),
			})

			// If there is an error it will be lifted up to the main scope.
			err = redisClient.RPush(r.Context(), ks.Session("chat"), aiResponse).Err()
			if err != nil {
				return err
			}

			// Write the response.
			w.Write([]byte(aiResponse))
			return nil
		})
	myLog("DONE!!!" + ss.SessionId)
	// Create a session if it does not exists and acquire it.
	return err
}

func SpeechToText(fileReader io.Reader) string {
	// TODO: This part is commented out to remove the little uncertainty given by
	// The use of an external API. Initial tests are done in a more predictable way
	// To ease the comparison between the different setups.
	/*
		ctx := context.Background()
		audioRequest := openai.AudioRequest{
			Model:  "whisper-1",
			Reader: fileReader,
		}

		response, err := client.CreateTranscription(ctx, audioRequest)
		if err != nil {
			log.Printf("Failed to transcribe audio: %v", err)
			return ""
		}

		return (response.Text)
	*/

	stringFile, err := io.ReadAll(fileReader)
	if err != nil {
		panic(err)
	}

	myLog("\n\n\n FILEEE " + string(stringFile))

	start := time.Now()
	for time.Since(start) < time.Duration((50+int(len(stringFile)/1024/100)))*time.Millisecond {
		// Keep the CPU busy
	}

	return string(stringFile)
}

func Handle(w http.ResponseWriter, r *http.Request) {
	// Check input...
	// Create options dynamically...
	// Do something...
	// Finally, Call the handler...
	myLog("Hello from the API handler!")
	handle.Handle(Node, w, r, options, speech_to_text)
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
