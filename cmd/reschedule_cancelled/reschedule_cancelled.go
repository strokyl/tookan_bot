package main


import (
	"fmt"
	"encoding/json"
	"github.com/strokyl/tookan_bot/pkg/restclient"
)

func main() {
	client := restclient.FromEnv()
	teamId := restclient.TeamIdFromEnv()
	tasks := client.GetCancelledTasks(teamId)
	success, failure := client.MarkTasksUnasigned(tasks)
	response, err := json.MarshalIndent(map[string]interface{}{
		"num_success": len(success),
		"num_failure": len(failure),
		"success": success,
		"failure": failure,

	}, "", "  ")

	if err != nil {
		panic(err)
	}

	fmt.Printf(string(response) + "\n")
}
