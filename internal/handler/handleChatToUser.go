package handler

import (
	"butter-socket/internal/hub"
	"butter-socket/models"
	"encoding/json"
	"fmt"
)

// when user want to connect with a human
// trigger name: transfer_chat
func handleChatTransferToUser(client *hub.Client, payload any) {
	fmt.Println("customer id ", client.Customer.Id)
	fmt.Println("chat transfer request")
	fmt.Println("customer-> sos flag waving to for connection: ", client.SosFlag)
	fmt.Println("payload: ", payload)
	fmt.Println("conversation: ", *client.Conversation)
	fmt.Println("conversation user: ", *client.Conversation.Customer)
	client.SosFlag = true
	fmt.Println("flag waving: ", client.SosFlag)

	//------------------------logics---------------------/////
	human := client.Hub.GetAllUsers()
	if len(human) > 0 {
		for _, user := range human {
			sendMessage(user, "transfer_chat", client.Conversation)
		}
	}
}

// trigger name: accept_chat (for users)
func handleHumanAcceptTheChat(client *hub.Client, payload any) {
	//-------paylad construction-------////
	payloadBytes, _ := json.Marshal(payload)
	var transferPayload models.Conversation
	fmt.Println(transferPayload)
	json.Unmarshal(payloadBytes, &transferPayload)
	fmt.Println("Accepted Payload(user conversation): ", *transferPayload.Customer)
	//-------------------------------------------//
	fmt.Println("user id ", client.User.UserID)
	customer := client.Hub.GetAllClients()[transferPayload.Customer.Id]
	customer.FlagRevealed = true
	customer.User = client.User
	fmt.Println("assigned user ", customer.User)
	sendMessage(client.Hub.GetAllClients()[transferPayload.Customer.Id], "connection_event", "human communication started")
}

// trigger name: message
func handleConversationWithHuman(client *hub.Client, payload any) {
	fmt.Println("customer id ", client.Customer.Id)
	fmt.Println(payload)
	fmt.Println(client.SosFlag)
	fmt.Println(client.FlagRevealed)
	fmt.Println("assigned user:", *client.User)
}
