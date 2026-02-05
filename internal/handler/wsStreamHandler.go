package handler

import (
	"butter-socket/internal/hub"
	"butter-socket/internal/llm"
	"butter-socket/models"
	"context"
	"encoding/json"
	"time"
)

func handleChatStreamMessage(client *hub.Client, payload any) {

	// 1. Parse user message
	payloadBytes, _ := json.Marshal(payload)
	var msgIn models.MsgInOut
	json.Unmarshal(payloadBytes, &msgIn)

	// 2. Cancel previous AI if still running
	if client.CancelAI != nil {
		client.CancelAI()
	}

	// 3. Create cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	client.CancelAI = cancel

	// 4. Tell frontend: AI started typing
	sendMessage(client, "typing_start", nil)

	var fullReply string

	// 5. Start streaming AI
	err := llm.StreamButterAI(ctx, msgIn.Content, func(token string) {

		fullReply += token

		// 6. Send token immediately
		sendMessage(client, "message_chunk", models.MsgInOut{
			SenderType:  "AI-AGENT",
			Content:     token,
			ContentType: "text",
			CreatedAt:   time.Now().Format(time.RFC3339),
		})
	})

	if err != nil {
		sendError(client, "AI error")
		return
	}

	// 7. Tell frontend: AI finished
	sendMessage(client, "typing_end", nil)

	// 8. Save fullReply to DB (later)
}
