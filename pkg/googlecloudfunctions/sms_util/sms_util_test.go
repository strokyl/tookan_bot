package sms_util

import (
	"testing"
	"github.com/stretchr/testify/suite"
	"github.com/strokyl/tookan_bot/pkg/restclient"
	"github.com/strokyl/tookan_bot/pkg/planning"
	"time"
)

type SendSmsUtilTestSuite struct {
	suite.Suite
}

func (suite *SendSmsUtilTestSuite) TestGenerateChangeTaskPlanURL() {
	suite.Equal(
		GenerateChangeTaskPlanURL(
			"secret",
			"http://localhost:8080/",
			restclient.JobId(42),
			time.Date(2020, time.May, 22, 0, 0, 0, 0, time.UTC),
			planning.UnsafeNewPlan("16:42 à 17:42"),
		),
		"http://localhost:8080/ChangeHour?jobId=42&plan=16%3A42+%C3%A0+17%3A42&date=22%2F05%2F2020&secret=FsynSpNVvQupS1dr8Ud3OmkEvM8%3D",
	)
}

func TestSendChangeHourEmail(t *testing.T) {
	suite.Run(t, new(SendSmsUtilTestSuite))
}
