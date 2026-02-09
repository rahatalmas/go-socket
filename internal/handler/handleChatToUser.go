package handler

import (
	"butter-socket/internal/hub"
	"butter-socket/models"
	"encoding/json"
	"fmt"
)

// trigger name: transfer_chat
func handleChatTransferToUser(client *hub.Client, payload any) {
	if !client.SosFlag {
		client.SosFlag = true
		fmt.Println("flag waving: ", client.SosFlag)
		connList := client.Hub.GetAllUserConnByCompanyId(client.Conversation.CompanyId)
		fmt.Println(len(connList))
		if len(connList) == 0 {
			unavilableMsgPayload := models.MsgInOut{
				SenderType: "system",
				SenderId:   "system",
				ReceiverId: client.Customer.Id,
				Content:    "no one is available to chat",
			}
			sendMessage(client, "connection_event", unavilableMsgPayload)
		} else {
			for _, conn := range connList {
				sendMessage(conn, "transfer_chat", payload)
			}
		}
	}

}

// trigger name: accept_chat (for users)
func handleHumanAcceptTheChat(client *hub.Client, payload any) {
	payloadBytes, _ := json.Marshal(payload)
	var transferPayload models.Conversation
	json.Unmarshal(payloadBytes, &transferPayload)

	if client.Hub.GetAllClients()[transferPayload.Customer.Id] != nil {
		client.SosFlag = true
		client.FlagRevealed = true
		customer := client.Hub.GetAllClients()[transferPayload.Customer.Id]
		customer.FlagRevealed = true
		customer.User = client.User

		unavilableMsgPayload := models.MsgInOut{
			SenderId:   "system",
			ReceiverId: transferPayload.Customer.Id,
			Content:    "human communication started",
		}
		sendMessage(client.Hub.GetAllClients()[transferPayload.Customer.Id], "connection_event", unavilableMsgPayload)
	}
}

// trigger name: message
func handleConversationWithHuman(client *hub.Client, payload any) {
	if client.Type == "customer" {
		user := client.User
		sendMessage(client.Hub.GetUserConnByUserId(user.UserID), "message", payload)
	} else {
		payloadBytes, _ := json.Marshal(payload)
		var msgPayload models.MsgInOut
		json.Unmarshal(payloadBytes, &msgPayload)
		customer := client.Hub.GetAllClients()[msgPayload.ReceiverId]
		sendMessage(customer, "message", payload)
	}
}

// trigger name: typing
// func typingUpdate(client *hub.Client, payload any) {
// 	if client.Type == "customer" {
// 		user := client.User
// 		sendMessage(client.Hub.GetUserConnByUserId(user.UserID), "message", payload)
// 	} else {
// 		payloadBytes, _ := json.Marshal(payload)
// 		var msgPayload models.MsgInOut
// 		json.Unmarshal(payloadBytes, &msgPayload)
// 		customer := client.Hub.GetAllClients()[msgPayload.ReceiverId]
// 		sendMessage(customer, "typing", payload)
// 	}
// }
