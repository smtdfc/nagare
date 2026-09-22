package rest

const (
	CheckAuthStateEndpoint = "/api/v1/user/auth/state"
)

type CheckAuthStateResponse struct {
	IsAuth bool `json:"isAuth"`
}
