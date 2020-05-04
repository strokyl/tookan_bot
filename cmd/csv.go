package main


import (
	"strconv"
	"os"
	"github.com/strokyl/tookan_bot/pkg/restclient"
	"encoding/csv"
	"strings"
	"regexp"
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

func removeComma(s string) string {
	return strings.ReplaceAll(s, ",", " ")
}

func extractRepasNum(task restclient.Task) string {
	repasRegexp := regexp.MustCompile(`(\d+) Repas`)
	result := repasRegexp.FindSubmatch([]byte(task.Name))
	if len(result) == 2 {
		return string(result[1])
	} else {
		return "?"
	}
}

func main() {
	tasks := client.GetCompletedTasks(teamId)

	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()
	writer.Write([]string{
		"job id",
		"nom",
		"description",
		"date de creation",
		"date de fin",
		"nombre de repas",
	})
	for _, task := range tasks {
		err := writer.Write([]string{
			strconv.Itoa(int(*task.JobId)),
			removeComma(task.Name),
			removeComma(task.Description),
			task.CreationDate,
			task.CompletedDate,
			extractRepasNum(task),
		})

		if err != nil {
			panic(err)
		}
	}
}
