package function

import (
	"fmt"
	"net/http"

	"github.com/my-ermes-labs/api-go/api"

	handle "github.com/my-ermes-labs/api-go/http"

	log "github.com/my-ermes-labs/log"
	rc "github.com/my-ermes-labs/storage-redis/packages/go"
	"github.com/redis/go-redis/v9"
)

// "github.com/redis/go-redis/v9"

// Request handler.
func Handler(w http.ResponseWriter, r *http.Request, sessionToken api.SessionToken) error {
	log.MyLog("\n\nHandler in api\n\n")
	// Create the keyspace.
	ks := rc.NewErmesKeySpaces(sessionToken.SessionId)

	// Execute the Redis commands in a transaction.
	cmds, err := redisClient.TxPipelined(r.Context(), func(pipe redis.Pipeliner) error {
		// Increment "counter" in the node keyspace.
		pipe.Incr(r.Context(), ks.Node("counter"))
		// Increment "counter" in the session keyspace.
		pipe.Incr(r.Context(), ks.Session("counter"))
		// No error.
		return nil
	})

	// Check if there was an error and return it. Error handling is done by the
	// Handle function.
	if err != nil {
		return err
	}

	// Get the node key.
	nodeKey := cmds[0].(*redis.IntCmd).Val()
	// Get the session key.
	sessionKey := cmds[1].(*redis.IntCmd).Val()
	// Create the response body.
	body := fmt.Sprintf("Node key: %d, Session key: %d", nodeKey, sessionKey)

	// Write the response.
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))
	// Return no error.
	return nil
}

func Handle(w http.ResponseWriter, r *http.Request) {
	// Check input...
	// Create options dynamically...
	// Do something...
	// Finally, Call the handler...
	fmt.Fprintln(w, "Hello from the API handler!")
	handle.Handle(node, w, r, options, Handler)
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
