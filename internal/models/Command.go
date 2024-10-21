package models

const (
	CommandDelete  RabbitCommand = "DELETE"
	CommandPurge   RabbitCommand = "PURGE"
	CommandAdd     RabbitCommand = "ADD"
	CommandGet     RabbitCommand = "Get"
	CommandGetByID RabbitCommand = "GetById"
	CommandUpdate  RabbitCommand = "Update"
)

type RabbitCommand string

type Msg struct {
	Command RabbitCommand `json:"command"`
	Key     interface{}   `json:"key"`
}
type Filter struct {
	Query *string `json:"query"`
}
