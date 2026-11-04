package rest

const (
	ListCredentialsEndpoint      = "/api/v1/user/credentials/list"
	GetCredentialDetailsEndpoint = "/api/v1/user/credentials/details"
	AddCredentialEndpoint        = "/api/v1/user/credentials/add"
	UpdateCredentialEndpoint     = "/api/v1/user/credentials/update"
	DeleteCredentialEndpoint     = "/api/v1/user/credentials/delete"
)

type Credential struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetListCredentialResponse struct {
	Credentials []*Credential `json:"credentials"`
}

type GetCredentialDetailsResponse struct {
	Credential *Credential `json:"credential"`
}

type AddCredentialRequest struct {
	Name   string `json:"name"`
	ApiKey string `json:"apiKey"`
}

type AddCredentialResponse struct {
	Credential *Credential `json:"credential"`
}

type UpdateCredentialRequest struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	ApiKey string `json:"apiKey"`
}

type UpdateCredentialResponse struct {
	Credential *Credential `json:"credential"`
}

type DeleteCredentialRequest struct {
	ID string `json:"id"`
}
