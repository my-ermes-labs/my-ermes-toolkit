package function

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/my-ermes-labs/api-go/api"
	"github.com/my-ermes-labs/api-go/infrastructure"
	rc "github.com/my-ermes-labs/storage-redis/packages/go"
	"github.com/redis/go-redis/v9"
)

// The node that the function is running on.
var Node *api.Node

// The Redis client.
var redisClient *redis.Client

func init() {
	// Get the node from the environment variable.
	encodedJsonNode := os.Getenv("ERMES_NODE")

	decodedJsonNode, err := base64.StdEncoding.DecodeString(encodedJsonNode)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Unmarshal the environment variable to get the node.
	infraNode, err := infrastructure.UnmarshalNode([]byte(decodedJsonNode))
	// Check if there was an error unmarshalling the node.
	if err != nil {
		panic(err)
	}

	// Get the Redis connection details from the environment variables.
	redisHost := envOrDefault("REDIS_HOST", "10.62.0.1")
	redisPort := envOrDefault("REDIS_PORT", "6379")
	redisPassword := envOrDefault("REDIS_PASSWORD", "")
	// Create a new Redis client.
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: redisPassword,
		DB:       0, // use default DB
	})

	if err := checkRedisConnection(redisClient); err != nil {
		// myLog(fmt.Sprintf("Errore nella connessione a Redis: %v", err))
		fmt.Println("Errore nella connessione a Redis:", err)
	} else {
		// myLog("Connessione a Redis riuscita!")
		fmt.Println("Connessione a Redis riuscita!")
	}

	// The Redis commands.
	var RedisCommands = rc.NewRedisCommands(redisClient)
	// Create a new node with the Redis commands.
	Node = api.NewNode(*infraNode, RedisCommands)
}

// Get the value of an environment variable or return a default value.
func envOrDefault(key string, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return value
}

func checkRedisConnection(client *redis.Client) error {
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	return err
}