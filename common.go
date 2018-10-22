package apigee

type Attribute struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type CredentialApiProduct struct {
	ApiProduct string `json:"apiproduct,omitempty"`
	Status     string `json:"status,omitempty"`
}

type Credential struct {
	ApiProducts    []CredentialApiProduct `json:"apiProducts,omitempty"`
	Attributes     []Attribute            `json:"attributes,omitempty"`
	ConsumerKey    string                 `json:"consumerKey,omitempty"`
	ConsumerSecret string                 `json:"consumerSecret,omitempty"`
	ExpiresAt      int                    `json:"expiresAt,omitempty"`
	IssuedAt       int                    `json:"issuedAt,omitempty"`
	Scopes         []string               `json:"scopes,omitempty"`
}

//This is just a placeholder
type App struct {
	ApigeeId string
}
