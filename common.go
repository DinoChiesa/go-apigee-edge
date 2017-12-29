package apigee

type Attribute struct {
	Name						string		`json:"name,omitempty"`
	Value						string		`json:"value,omitempty"`
}

//This is just a placeholder
type App struct {
	ApigeeId					string		`json:"apigee_id,omitempty"`
}
