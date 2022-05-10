package app

import (
	"fmt"
)

// MessageChan ...
var MessageChan = make(chan MessageDetails)

// MessageDetails ...
type MessageDetails struct {
	MessageID string
	Message   string
	ToNumber  string
}

func init() {
	// message := <-Message
	go SendSMSMessages(MessageChan)
}

// SendSMSMessages ...
func SendSMSMessages(message chan MessageDetails) {
	for {
		message := <-message
		//Call the Gateway, and pass the constants here!
		smsService := NewSMSService("prepaidmetering", "9db98f9884d801a3ec25c78d36505f0eeb27395041a5e8fd381a4d4f4018183f", "production")

		//Send SMS - REPLACE Recipient and Message with REAL Values
		recipients, err := smsService.Send("", message.ToNumber, message.Message) //Leave blank, "", if you don't have one)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(recipients)
	}

}
