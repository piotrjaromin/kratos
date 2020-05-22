package attack

type AttackMode string

const (
	json     AttackMode = "json"
	jsScript AttackMode = "jsScript"
)

type AttachBody struct {
	AttackMode AttackMode
}

type AttackResponse struct {
	Message string `json:"message"`
}
