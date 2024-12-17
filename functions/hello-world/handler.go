package function

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ermes-labs/api-go/infrastructure"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	var input []byte
	var node infrastructure.Node = infrastructure.Node{}

	if r.Body != nil {
		defer r.Body.Close()

		body, _ := io.ReadAll(r.Body)

		input = body
	}

	encodedEnv, ok := os.LookupEnv("ERMES_NODE")
	if !ok {
		panic("ermes-node env variable is not set")
	}

	decodedEnv, err := base64.StdEncoding.DecodeString(encodedEnv)
	if err != nil {
		fmt.Println("Errore nella decodifica:", err)
		return
	}

	encodedNode, okk := os.LookupEnv("NODE_STRING")
	if !okk {
		panic("node_string env variable is not set")
	}

	decodedNode, errrr := base64.StdEncoding.DecodeString(encodedNode)
	if errrr != nil {
		fmt.Println("Errore nella decodifica:", err)
		return
	}

	// QUESTO DA INTERNAL SERVER ERROR QUI E NON FA RUNNING FUNCTIONE IN NODE.
	errr := json.Unmarshal([]byte(decodedEnv), &node)

	if errr != nil {
		panic(fmt.Sprintf("Error unmarshalling ERMES_NODE: %v", errr))
	}

	w.WriteHeader(http.StatusOK)
	// w.Write([]byte(fmt.Sprintf("Body: %s \nENV= %v", string(input), node)))
	w.Write([]byte(fmt.Sprintf("Body: %s \nDecoded env = %v, \nEncoded env = %s \nNODE: %s", string(input), decodedEnv, encodedEnv, string(decodedNode))))

}
