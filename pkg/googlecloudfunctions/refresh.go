package googlecloudfunctions

import (
	"fmt"
	"net/http"
	"encoding/json"
)

func Refresh(w http.ResponseWriter, r *http.Request) {
	if !checkSecret(w, r, secret) {
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
