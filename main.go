package main

import "github.com/go-resty/resty/v2"
import "strconv"
import "fmt"
import "os"
import "encoding/json"

type Response struct {
	Status int `json:"status"`
	Message string `json:"message"`
	Data *[]Task `json:"data"`
}

type JobId int

type Task struct {
	JobId *JobId `json:"job_id"`
}

func autoAssignATask(token string, teamId int, id JobId) bool {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetBody(map[string]interface{}{
			"api_key": token,
			"job_id": id,
		}).
		Post("https://api.tookanapp.com/v2/re_autoassign_task")

	if err != nil {
		panic(err)
	}

	if resp.IsError() {
		panic(fmt.Sprintf("Http error %s: %s", resp.Status(), resp.String()))
	}

	result := Response{}
	json.Unmarshal(resp.Body(), &result)

	if result.Status != 200 {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Could not auto affect task %d: %+v\n", id, result))
		return false
	}
	return true
}

func autoAssignTasks(token string, teamId int, ids []JobId) {
	for _,id := range ids {
		autoAssignATask(token, teamId, id)
	}
}

func getUnasignedTask(token string, teamId int) []JobId {
	const unasigned_status = 6
	const pickup_type = 0

	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetBody(map[string]interface{}{
			"api_key": token,
			"job_type": pickup_type,
			"job_status": unasigned_status,
			"team_id": teamId,
		}).
		Post("https://api.tookanapp.com/v2/get_all_tasks")

	if err != nil {
		panic(err)
	}

	if resp.IsError() {
		panic(fmt.Sprintf("Http error %s: %s", resp.Status(), resp.String()))
	}

	result := Response{}
	json.Unmarshal(resp.Body(), &result)

	if result.Status != 200 {
		panic(fmt.Sprintf("Could not fetch unasigned tasks: %+v\n", result))
	}

	if result.Data == nil {
		panic("Result has no data field: " + resp.String())
	}

	ids := make([]JobId, len(*result.Data))
	for i, data := range *result.Data {
		if data.JobId == nil {
			panic(fmt.Sprintf("Task has no job_id", result))
		}
		ids[i] = *data.JobId
	}
	return ids
}

func main() {

	token := os.Getenv("TOOKAN_V2_TOKEN")
	if token == "" {
		panic("Please set TOOKAN_V2_TOKEN")
	}
	teamIdStr := os.Getenv("TOOKAN_TEAM_ID")
	if teamIdStr == "" {
		panic("Please set TOOKAN_TEAM_ID")
	}
	teamId, err := strconv.Atoi(teamIdStr)
	if err != nil {
		panic("Please set a TOOKAN_TEAM_ID that is an integer")
	}


	ids := getUnasignedTask(token, teamId)
	autoAssignTasks(token, teamId, ids)
	fmt.Printf("%+v\n", ids)
}
