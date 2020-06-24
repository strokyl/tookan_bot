package planning

import (
	"testing"
	"github.com/strokyl/tookan_bot/pkg/restclient"
	"github.com/stretchr/testify/suite"
	"time"
)

type PlanningTestSuite struct {
	suite.Suite
}

func formatDateWithPadding(t time.Time) string {
	return t.Format("2/1/2006")
}

func formatDateWithoutPadding(t time.Time) string {
	return t.Format("2/1/2006")
}

func (suite *PlanningTestSuite) TestTaskForTodayShouldReturnTrueForTaskWithoutDate() {
	suite.True(TaskForToday("yolo 15:00 à 18:00"))
}

func (suite *PlanningTestSuite) TestTaskForTodayShouldReturnTrueForTaskWithTodayDate() {
	suite.True(TaskForToday("yolo 15:00 à 18:00 " + formatDateWithPadding(time.Now()) + " !"))
}

func (suite *PlanningTestSuite) TestTaskForTodayShouldReturnTrueForTaskWithUnpaddedTodayDate() {
	suite.True(TaskForToday("yolo 15:00 à 18:00 " + formatDateWithoutPadding(time.Now()) + " !"))
}

func (suite *PlanningTestSuite) TestTaskForTodayShouldReturnFalseForTaskWithTomorrowDate() {
	suite.False(TaskForToday("yolo 15:00 à 18:00 " + formatDateWithPadding(time.Now().Add(24*time.Hour)) + " !"))
}

func (suite *PlanningTestSuite) TestTaskForTodayShouldReturnFalseForTaskWithUnpaddedTomorrowDate() {
	suite.False(TaskForToday("yolo 15:00 à 18:00 " + formatDateWithoutPadding(time.Now().Add(24*time.Hour)) + " !"))
}

var aDate = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
func (suite *PlanningTestSuite) TestReplacePlanLabelOnTaskNameWhenNoDateBefore() {
	suite.Equal(ReplacePlanLabelOnTaskName("yolo 9:2 à 19:23 2 repas", UnsafeNewPlan("12:00 à 18:00"), aDate),
		"yolo 12:00 à 18:00 10/11/2009 2 repas")
}

func (suite *PlanningTestSuite) TestReplacePlanLabelOnTaskNameWhenDateBefore() {
	suite.Equal(ReplacePlanLabelOnTaskName("yolo 9:2 à 19:23 23/10/2009 2 repas", UnsafeNewPlan("12:00 à 18:00"), aDate),
		"yolo 12:00 à 18:00 10/11/2009 2 repas")
}

func (suite *PlanningTestSuite) TestReplacePlanLabelOnTaskNameWhenDateBeforeButMisplaced() {
	suite.Equal(ReplacePlanLabelOnTaskName("23/10/2009 yolo 9:2 à 19:23 2 repas", UnsafeNewPlan("12:00 à 18:00"), aDate),
		"yolo 12:00 à 18:00 10/11/2009 2 repas")
}

func (suite *PlanningTestSuite) TestReplacePlanLabelOnTaskNameWhenDateBeforeButMisplacedInAnOtherWay() {
	suite.Equal(ReplacePlanLabelOnTaskName("yolo 09:2 à 9:23 2 repas 23/10/2009", UnsafeNewPlan("12:00 à 18:00"), aDate),
		"yolo 12:00 à 18:00 10/11/2009 2 repas")
}

func (suite *PlanningTestSuite) TestReplacePlanLabelOnTaskNameWhenMultipleLabel() {
	suite.Equal(ReplacePlanLabelOnTaskName("11:00 à 14:00, \n14:00 à 18:00, \n18:00 à 20:30 #PourEux François 2 Repas ", UnsafeNewPlan("12:00 à 18:00"), aDate),
		"12:00 à 18:00 10/11/2009 #PourEux François 2 Repas")
}

func (suite *PlanningTestSuite) TestReplacePlanLabelWhenNoLabel() {
	suite.Equal(ReplacePlanLabelOnTaskName("yolo 2 repas", UnsafeNewPlan("12:00 à 18:00"), aDate),
		"yolo 2 repas")
}

func (suite *PlanningTestSuite) TestFindTodayPlanTask() {
	suite.Equal(FindTodayPlanTask(restclient.Task{Name: "yolo 2 repas"}),
		[]APlan{})
	suite.Equal(FindTodayPlanTask(restclient.Task{Name: "yolo 2 repas 12:00 à 18:00"}),
		[]APlan{
			APlan{label: "12:00 à 18:00"},
		})

	suite.Equal(FindTodayPlanTask(restclient.Task{Name: "yolo 2 repas 12:00 à 18:00\n 13:00 à 14:00,18:00 à 23:00"}),
		[]APlan{
			APlan{label: "12:00 à 18:00"},
			APlan{label: "13:00 à 14:00"},
			APlan{label: "18:00 à 23:00"},
		})
}

func (suite *PlanningTestSuite) TestExtractEndInTotalMinutes() {
	myPlan := APlan{label: "18:00 à 19:24"}
	suite.Equal(myPlan.extractEndInTotalMinutes(), 19*60 + 24)
	suite.Equal(myPlan.extractEndInTotalMinutes(), 19*60 + 24)
}

func (suite *PlanningTestSuite) TestExtractEndInTotalMinutesWhenNoPadding() {
	myPlan := APlan{label: "18:00 à 9:4"}
	suite.Equal(myPlan.extractEndInTotalMinutes(), 9*60 + 4)
}

func getADateForGivenHour(hour, minute int) time.Time {
	location, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		panic(err)
	}
	return time.Date(2020, time.May, 19, hour, minute, 0, 0, location)
}

func (suite *PlanningTestSuite) TestPlanEndBefore() {
	myPlan := APlan{label: "18:00 à 19:24"}
	suite.True(myPlan.EndBefore(getADateForGivenHour(19,24)))
	suite.True(myPlan.EndBefore(getADateForGivenHour(19,25)))
	suite.True(myPlan.EndBefore(getADateForGivenHour(20,00)))

	suite.False(myPlan.EndBefore(getADateForGivenHour(19,23)))
	suite.False(myPlan.EndBefore(getADateForGivenHour(18,55)))
	suite.False(myPlan.EndBefore(getADateForGivenHour(18,00)))
}

func (suite *PlanningTestSuite) TestAllEndBefore() {
	myPlans := []APlan{
		APlan{label: "15:00 à 16:00"},
		APlan{label: "18:00 à 19:24"},
	}
	suite.True(AllEndBefore(myPlans, getADateForGivenHour(19,24)))
	suite.True(AllEndBefore(myPlans, getADateForGivenHour(19,25)))
	suite.True(AllEndBefore(myPlans, getADateForGivenHour(20,00)))

	suite.False(AllEndBefore(myPlans, getADateForGivenHour(19,23)))
	suite.False(AllEndBefore(myPlans, getADateForGivenHour(18,55)))
	suite.False(AllEndBefore(myPlans, getADateForGivenHour(18,00)))
}

func TestPlanning(t *testing.T) {
	suite.Run(t, new(PlanningTestSuite))
}
