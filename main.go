package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
	"github.com/twilio/twilio-go/twiml"
)

type Alert struct {
	Status string `json:"status"`
	Labels Labels `json:"labels"`
	// Additional alert fields...
}

type Labels struct {
	Alertname string `json:"alertname"`
	Hostname  string `json:"hostname"`
	// Additional label fields...
}

type AlertManagerPayload struct {
	Alerts []Alert `json:"alerts"`
	// Additional payload fields...
}

func main() {
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	callUrl := os.Getenv("TWILIO_CALL_URL")
	callTo := os.Getenv("TWILIO_CALL_TO")
	callFrom := os.Getenv("TWILIO_CALL_FROM")

	if accountSid == "" || authToken == "" || callUrl == "" || callTo == "" || callFrom == "" {
		fmt.Println("Twilio account SID and/or auth token not provided")
		os.Exit(1)
	}

	clientParams := twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	}

	router := gin.Default()

	router.POST("/answer", func(context *gin.Context) {
		var payload AlertManagerPayload
		if err := context.ShouldBindJSON(&payload); err != nil {
			context.String(http.StatusBadRequest, "Invalid JSON payload")
			return
		}

		alertInfo := ""
		if len(payload.Alerts) > 0 {
			alert := payload.Alerts[0]
			// Only call if alert is firing
			if alert.Status == "firing" {
				alertInfo = fmt.Sprintf("Alert %s is firing for host %s", alert.Labels.Alertname, alert.Labels.Hostname)
			}

			makeCall(clientParams, callUrl, callTo, callFrom)
		}

		say := &twiml.VoiceSay{
			Message: alertInfo,
		}

		twimlResult, err := twiml.Voice([]twiml.Element{say})
		if err != nil {
			context.String(http.StatusInternalServerError, err.Error())
		} else {
			context.Header("Content-Type", "text/xml")
			context.String(http.StatusOK, twimlResult)
		}
	})

	router.Run(":1337")
}

func makeCall(clientParams twilio.ClientParams, callUrl string, callTo string, callFrom string) {
	client := twilio.NewRestClientWithParams(clientParams)

	params := &api.CreateCallParams{}
	params.SetUrl(callUrl)
	params.SetTo(callTo)
	params.SetFrom(callFrom)

	resp, err := client.Api.CreateCall(params)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		if resp.Sid != nil {
			fmt.Println(*resp.Sid)
		} else {
			fmt.Println(resp.Sid)
		}
	}
}
