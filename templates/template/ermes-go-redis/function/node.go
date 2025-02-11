package function

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/my-ermes-labs/api-go/api"
	"github.com/my-ermes-labs/api-go/infrastructure"
	log "github.com/my-ermes-labs/log"
	rc "github.com/my-ermes-labs/storage-redis/packages/go"
	"github.com/redis/go-redis/v9"
)

// The node that the function is running on.
var node *api.Node

// The Redis client.
var redisClient *redis.Client

func init() {
	loggg("\n\node in api\n\n")
	log.MyLog("initializing API")
	// Get the node from the environment variable.
	encodedJsonNode := os.Getenv("ERMES_NODE")

	decodedJsonNode, err := base64.StdEncoding.DecodeString(encodedJsonNode)
	if err != nil {
		fmt.Println("Errore nella decodifica:", err)
		return
	}

	// Unmarshal the environment variable to get the node.
	infraNode, err := infrastructure.UnmarshalNode([]byte(decodedJsonNode))
	// Check if there was an error unmarshalling the node.
	if err != nil {
		panic(err)
	}

	log.MyLog(fmt.Sprintf("NODE ==> Area name = %s Host = %s ", infraNode.AreaName, infraNode.Host))

	// Get the Redis connection details from the environment variables.
	redisHost := envOrDefault("REDIS_HOST", "10.62.0.1")
	redisPort := envOrDefault("REDIS_PORT", "6379")
	redisPassword := envOrDefault("REDIS_PASSWORD", "")

	log.MyLog(fmt.Sprintf("REDIS_HOST = %s", redisHost))
	log.MyLog(fmt.Sprintf("REDIS_PORT = %s", redisPort))
	log.MyLog(fmt.Sprintf("REDIS_PW = %s", redisPassword))

	// Create a new Redis client.
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: redisPassword,
		DB:       0, // use default DB
	})

	if err := checkRedisConnection(redisClient); err != nil {
		log.MyLog(fmt.Sprintf("Errore nella connessione a Redis: %v", err))
		fmt.Println("Errore nella connessione a Redis:", err)
	} else {
		log.MyLog("Connessione a Redis riuscita!")
		fmt.Println("Connessione a Redis riuscita!")
	}

	os.Setenv("REDIS_CLIENT", redisClient.Options().Addr)

	// // The Redis commands.
	var RedisCommands = rc.NewRedisCommands(redisClient)

	// // Create a new node with the Redis commands.
	node = api.NewNode(*infraNode, RedisCommands)

	os.Setenv("NODE", "area name = "+node.AreaName+" Host = "+node.Host)
}

// Get the value of an environment variable or return a default value.
func envOrDefault(key string, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return value
}

// Get the value of an environment variable or panic if it is not set.
func envOrPanic(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		panic(key + " env variable is not set")
	}
	return value
}

func checkRedisConnection(client *redis.Client) error {
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	return err
}

func loggg(bodyContent string) (string, error) {
	url := "http://192.168.64.1:3000/"

	requestBody := bytes.NewBufferString(bodyContent)

	req, err := http.NewRequest("POST", url, requestBody)
	if err != nil {
		return "", fmt.Errorf("error while creating the request: %v", err)
	}

	// Imposta l'header Content-Type per indicare che stiamo inviando testo semplice
	req.Header.Set("Content-Type", "text/plain")

	// Invia la richiesta con un client HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error while sending the request: %v", err)
	}
	defer resp.Body.Close()

	// Legge la risposta del server
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error while reading the response: %v", err)
	}

	return string(responseBody), nil
}
