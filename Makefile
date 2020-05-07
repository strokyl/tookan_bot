compile:
	go build ./...
all: compile

deploy_refresh:
	gcloud functions deploy Refresh --runtime go113 --trigger-http --set-env-vars TOOKAN_V2_TOKEN=$$TOOKAN_V2_TOKEN --set-env-vars SECRET=$$SECRET

deploy_reopen:
	gcloud functions deploy Reopen --runtime go113 --trigger-http --set-env-vars TOOKAN_V2_TOKEN=$$TOOKAN_V2_TOKEN --set-env-vars SECRET=$$SECRET
