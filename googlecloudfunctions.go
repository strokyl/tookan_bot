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

	secret = os.Getenv("SECRET")
	if secret == "" {
		panic("Please set SECRET")
	}

	client = restclient.FromEnv()
}


func checkSecret(w http.ResponseWriter, r *http.Request) bool {
	givenSecret, ok := r.URL.Query()["secret"]
	if !ok || len(givenSecret) != 1 || givenSecret[0] != secret {
		fmt.Fprintf(w, "Bad secret: ")
		return false
	}
	return true
}

func getTeamId(w http.ResponseWriter, r *http.Request) (int, error) {
	teamIdStr, ok := r.URL.Query()["teamId"]
	if !ok || len(teamIdStr) != 1 ||teamIdStr[0] == "" {
		fmt.Fprintf(w, "Please teamId")
		return 0, fmt.Errorf("Please set teamId")
	}
	var err error
	teamId, err := strconv.Atoi(teamIdStr[0])
	if err != nil {
		fmt.Fprintf(w, "Please set a teamId that is a integer")
		return 0, err
	}

	return teamId, nil
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	if !checkSecret(w, r) {
		return
	}
	teamId, err := getTeamId(w, r)
	if err != nil {
		return
	}
	tasks := client.GetUnasignedTasks(teamId)
	success, failure := client.AutoAssignTasks(tasks)


	response, err := json.MarshalIndent(map[string]interface{}{
		"num_success": len(success),
		"num_failure": len(failure),
		"success": success,
		"action": "Notify all agent about all unasigned tasks",
		"failure": failure,

	}, "", "  ")

	if err != nil {
		panic(err)
	}

	fmt.Fprintf(w, string(response) + "\n")
}

func Reopen(w http.ResponseWriter, r *http.Request) {
	if !checkSecret(w, r) {
		return
	}
	teamId, err := getTeamId(w, r)
	if err != nil {
		return
	}
	tasks := client.GetCancelledTasks(teamId)
	success, failure := client.MarkTasksUnasigned(tasks)


	response, err := json.MarshalIndent(map[string]interface{}{
		"num_success": len(success),
		"num_failure": len(failure),
		"action": "mark all cancelled tasks as unasigned",
		"success": success,
		"failure": failure,

	}, "", "  ")

	if err != nil {
		panic(err)
	}

	fmt.Fprintf(w, string(response) + "\n")
}
