compile:
	go build ./...
all: compile

deploy_refresh:
	gcloud functions deploy Refresh --runtime go113 --trigger-http --set-env-vars TOOKAN_V2_TOKEN=$$TOOKAN_V2_TOKEN --set-env-vars SECRET=$$SECRET

deploy_reopen:
	gcloud functions deploy Reopen --runtime go113 --trigger-http --set-env-vars TOOKAN_V2_TOKEN=$$TOOKAN_V2_TOKEN --set-env-vars SECRET=$$SECRET


deploy_reopen:
	gcloud functions deploy Reopen --runtime go113 --trigger-http --set-env-vars TOOKAN_V2_TOKEN=$$TOOKAN_V2_TOKEN --set-env-vars SECRET=$$SECRET

deploy_change_hour_form:
	gcloud functions deploy ChangeHourForm --runtime go113 --trigger-http --set-env-vars TOOKAN_V2_TOKEN=$$TOOKAN_V2_TOKEN --set-env-vars SECRET=$$SECRET --set-env-vars BASE_URL=$$BASE_URL --set-env-vars SMS_SECRET=$$SMS_SECRET --set-env-vars SMS_SID=$$SMS_SID --set-env-vars SMS_NUMBER=$$SMS_NUMBER

deploy_change_hour:
	gcloud functions deploy ChangeHour --runtime go113 --trigger-http --set-env-vars TOOKAN_V2_TOKEN=$$TOOKAN_V2_TOKEN --set-env-vars SECRET=$$SECRET --set-env-vars BASE_URL=$$BASE_URL --set-env-vars SMS_SECRET=$$SMS_SECRET --set-env-vars SMS_SID=$$SMS_SID --set-env-vars SMS_NUMBER=$$SMS_NUMBER


deploy_send_change_hour_sms:
	gcloud functions deploy SendChangeHourSMS --runtime go113 --trigger-http --set-env-vars TOOKAN_V2_TOKEN=$$TOOKAN_V2_TOKEN --set-env-vars SECRET=$$SECRET --set-env-vars BASE_URL=$$BASE_URL --set-env-vars SMS_SECRET=$$SMS_SECRET --set-env-vars SMS_SID=$$SMS_SID --set-env-vars SMS_NUMBER=$$SMS_NUMBER

deploy_cancel_task:
	gcloud functions deploy CancelTask --runtime go113 --trigger-http --set-env-vars TOOKAN_V2_TOKEN=$$TOOKAN_V2_TOKEN --set-env-vars SECRET=$$SECRET --set-env-vars BASE_URL=$$BASE_URL --set-env-vars SMS_SECRET=$$SMS_SECRET --set-env-vars SMS_SID=$$SMS_SID --set-env-vars SMS_NUMBER=$$SMS_NUMBER
