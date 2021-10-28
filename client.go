package apigee

import (
	"net/http"
	"net/url"
)

// ApigeeClient manages communication with Apigee V1 Admin API.
type ApigeeClient struct {
	// HTTP client used to communicate with the Edge API.
	client *http.Client

	auth  *AdminAuth
	debug bool

	// defaults to https://login.apigee.com
	LoginBaseUrl string

	// Base URL for API requests.
	BaseURL *url.URL

	// Optional. tells whether to try to obtain a token or not.
	WantToken bool

	// User agent for client
	UserAgent string

	// Services used for communicating with the API
	Proxies      ProxiesService
	SharedFlows  SharedFlowsService
	Products     ProductsService
	Developers   DevelopersService
	Environments EnvironmentsService
	Organization OrganizationService
	Caches       CachesService
	Options      ApigeeClientOptions

	// Account           AccountService
	// Actions           ActionsService
	// Domains           DomainsService
	// DropletActions    DropletActionsService
	// Images            ImagesService
	// ImageActions      ImageActionsService
	// Keys              KeysService
	// Regions           RegionsService
	// Sizes             SizesService
	// FloatingIPs       FloatingIPsService
	// FloatingIPActions FloatingIPActionsService
	// Storage           StorageService
	// StorageActions    StorageActionsService
	// Tags              TagsService

	// Optional function called after every successful request made to the DO APIs
	onRequestCompleted RequestCompletionCallback
}

// wrap the standard http.Response returned from Apigee Edge. (why?)
type Response struct {
	*http.Response
}

// An ErrorResponse reports the error caused by an API request
type ErrorResponse struct {
	// HTTP response that caused this error
	Response *http.Response

	// Error message - maybe the json for this is "fault"
	Message string `json:"message"`
}

type ApigeeClientOptions struct {
	HttpClient *http.Client

	// Optional. The Apigee Admin base URL. For example, if using OPDK this might be
	// http://192.168.10.56:8080 . It defaults to https://api.enterprise.apigee.com
	MgmtUrl string

	// defaults to https://login.apigee.com
	LoginBaseUrl string

	// Specify the Edge organization name.
	Org string

	// Required. Authentication information for the Apigee Management server.
	Auth *AdminAuth

	// Optional. Warning: if set to true, HTTP Basic Auth base64 blobs will appear in output.
	Debug bool

	// Optional. tells whether to try to obtain a token or not.
	WantToken bool
}

// AdminAuth holds information about how to authenticate to the Apigee Management server.
type AdminAuth struct {
	// Optional. The path to the .netrc file that holds credentials for the Edge Management server.
	// By default, this is ${HOME}/.netrc .  If you specify a Password, this option is ignored.
	NetrcPath string

	// Optional. The username to use when authenticating to the Edge Management server.
	// Ignored if you specify a NetrcPath.
	Username string

	// Optional. Used if you explicitly specify a Password.
	Password string

	// Optional. This gets populated with a token by the client
	Token string
}
