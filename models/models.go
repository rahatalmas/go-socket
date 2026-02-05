package models

type MetaData struct {
	CreatedAt   string `json:"created_at"`
	LastUpdated string `json:"last_updated"`
}

type Customer struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	ProfileUri string `json:"profile_uri"`
	Source     string `json:"source"`
	CompanyId  string `json:"company_id"`
}

type Message struct {
	MetaData
	Id             string `json:"id"`
	ConversationId string `json:"conversation_id"`
	SenderId       string `json:"sender_id"`
	SenderType     string `json:"sender_type"`
	Content        string `json:"content"`
	ContentType    string `json:"content_type"`
}

type Conversation struct {
	MetaData
	Customer
	Messages     []Message `json:"messages"`
	Id           string    `json:"id"`
	Status       string    `json:"status"`
	Provider     string    `json:"provider"`
	Summary      string    `json:"summary"`
	Tags         []string  `json:"tags"`
	CompanyId    string    `json:"company_id"`
	DepartmentId string    `json:"department_id"`
	AssignedTo   string    `json:"assigned_to"`
	Source       string
}

// WebSocket message types
type WSMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"` //msg in out
}

type MsgInOut struct {
	SenderType  string `json:"sender_type"`
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
	CreatedAt   string `json:"created_at"`
}
