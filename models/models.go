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
	*MetaData    `json:"metadata"`
	*Customer    `json:"customer"`
	*User        `json:"user"`
	Messages     []Message `json:"messages"`
	Id           string    `json:"id"`
	Status       string    `json:"status"`
	Provider     string    `json:"provider"`
	Summary      string    `json:"summary"`
	Tags         []string  `json:"tags"`
	CompanyId    string    `json:"company_id"`
	DepartmentId string    `json:"department_id"`
	AssignedTo   string    `json:"assigned_to"`
	Source       string    `json:"source"`
}

// WebSocket message types
type WSMessage struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"` //msg in out
}

// payload for -> trigger: message
type MsgInOut struct {
	SenderId    string `json:"sender_id"`
	SenderType  string `json:"sender_type"`
	ReceiverId  string `json:"receiver_id,omitempty"`
	Typing      string `json:"typing,omitempty"`
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
	CreatedAt   string `json:"created_at,omitempty"`
}

// payload for -> trigger: transfer_chat
type TransferChatPayload struct {
}

//payload for -> trigger: accept_chat

// api response for user data
type EssentialResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	User      User   `json:"data"`
	Timestamp string `json:"timestamp"`
	Path      string `json:"path"`
}

type User struct {
	UserID      string       `json:"userId"`
	CompanyID   string       `json:"companyId"`
	Departments []Department `json:"departments"`
}

type Department struct {
	DepartmentID   string `json:"department_id"`
	DepartmentName string `json:"department_name"`
}
