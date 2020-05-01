package googlecloudfunctions

import (
	"fmt"
	"net/http"
	"os"
	"github.com/strokyl/tookan_bot/pkg/restclient"
	"encoding/json"
	"strconv"
)

var client *restclient.TookanClient
var secret string

func init() {
	token := os.Getenv("TOOKAN_V2_TOKEN")
	if token == "" {
		panic("Please set TOOKAN_V2_TOKEN")
	}

	secret = os.Getenv("SECRET")
	if secret == "" {
		panic("Please set SECRET")
	}

	client = restclient.New(token)
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	givenSecret, ok := r.URL.Query()["secret"]
	if !ok || len(givenSecret) != 1 || givenSecret[0] != secret {
		fmt.Fprintf(w, "Bad secret: ")
		return
	}
	teamIdStr, ok := r.URL.Query()["teamId"]
	if !ok || len(teamIdStr) != 1 ||teamIdStr[0] == "" {
		fmt.Fprintf(w, "Please teamId")
		return
	}
	var err error
	teamId, err := strconv.Atoi(teamIdStr[0])
	if err != nil {
		panic("Please set a TOOKAN_TEAM_ID that is an integer")
	}
	tasks := client.GetUnasignedTask(teamId)
	success, failure := client.AutoAssignTasks(tasks)


	response, err := json.MarshalIndent(map[string]interface{}{
		"num_success": len(success),
		"num_failure": len(failure),
		"success": success,
		"failure": failure,

	}, "", "  ")

	if err != nil {
		panic(err)
	}

	fmt.Fprintf(w, string(response) + "\n")
}
