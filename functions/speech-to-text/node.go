package function

import (
	// "encoding/json"

	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/ermes-labs/api-go/infrastructure"
)

func init() {
	var node infrastructure.Node

	// Get env variable
	encodedEnv, ok := os.LookupEnv("ERMES_NODE")

	if !ok {
		panic("ermes-node env variable is not set")
	}

	decodedEnv, err := base64.StdEncoding.DecodeString(encodedEnv)
	if err != nil {
		fmt.Println("Errore nella decodifica:", err)
		return
	}

	b := []byte(decodedEnv)
	errr := json.Unmarshal(b, &node)

	nodeString := node.String()
	encodedNodeString := base64.StdEncoding.EncodeToString([]byte(nodeString))
	os.Setenv("NODE_STRING", encodedNodeString)

	if errr != nil {
		panic(errr)
	}
}
