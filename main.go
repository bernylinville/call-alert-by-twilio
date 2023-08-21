package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

type Alert struct {
	Status      string      `json:"status"`
	Labels      Labels      `json:"labels"`
	Annotations Annotations `json:"annotations"`
}

type Labels struct {
	Alertname string `json:"alertname"`
}

type Annotations struct {
	Summary string `json:"summary"`
}

type AlertManagerPayload struct {
	Alerts []Alert `json:"alerts"`
}

func main() {
	callTo := os.Getenv("TWILIO_CALL_TO")
	callFrom := os.Getenv("TWILIO_CALL_FROM")
	webPort := os.Getenv("WEB_PORT")

	if webPort == "" {
		webPort = "1337"
	}

	if callTo == "" || callFrom == "" {
		fmt.Println("Please set TWILIO_CALL_TO and TWILIO_CALL_FROM environment variables")
		os.Exit(1)
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
				alertInfo = fmt.Sprintf("The alert name is %s, the summary is %s. The more info please check telegram and alertmanager.", alert.Labels.Alertname, alert.Annotations.Summary)
			}
		}

		twimlResult := fmt.Sprintf("<Response><Say>%s</Say></Response>", alertInfo)

		client := twilio.NewRestClient()

		params := &api.CreateCallParams{}
		params.SetTo(callTo)
		params.SetFrom(callFrom)
		params.SetTwiml(twimlResult)

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
	})

	router.Run(":" + webPort)
}
