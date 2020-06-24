package main


import (
	"fmt"
	"encoding/json"
	"github.com/strokyl/tookan_bot/pkg/restclient"
	"github.com/strokyl/tookan_bot/pkg/planning"
)

func main() {
	client := restclient.FromEnv()
	teamId := restclient.TeamIdFromEnv()
	tasks := client.GetUnasignedTasks(teamId)
	inTime, outDated, notForToday  := planning.SeparatePlannedAndUnplannedTask(tasks)
	response, err := json.MarshalIndent(map[string][]restclient.Task{
		"inTime" : inTime,
		"outDated" : outDated,
		"notForToday" : notForToday,
	}, "", "  ")

	if err != nil {
		panic(err)
	}

	fmt.Printf(string(response) + "\n")
}
