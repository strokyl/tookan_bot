package main


import (
	"strconv"
	"fmt"
	"os"
	"github.com/strokyl/tookan_bot/pkg/restclient"
	"encoding/json"
)

var client *restclient.TookanClient
var teamId int

func init() {
	token := os.Getenv("TOOKAN_V2_TOKEN")
	if token == "" {
		panic("Please set TOOKAN_V2_TOKEN")
	}
	teamIdStr := os.Getenv("TOOKAN_TEAM_ID")
	if teamIdStr == "" {
		panic("Please set TOOKAN_TEAM_ID")
	}
	var err error
	teamId, err = strconv.Atoi(teamIdStr)
	if err != nil {
		panic("Please set a TOOKAN_TEAM_ID that is an integer")
	}

	client = restclient.New(token)
}

func main() {
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

	fmt.Printf(string(response) + "\n")
}
