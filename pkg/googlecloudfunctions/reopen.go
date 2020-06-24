package googlecloudfunctions

import (
	"fmt"
	"net/http"
	"encoding/json"
)

func Reopen(w http.ResponseWriter, r *http.Request) {
	if !checkSecret(w, r, secret) {
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
