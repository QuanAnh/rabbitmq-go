package endpoints

type UserRequest struct {
	Id	string				`json:"id"`
}
type UserResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"` // errors don't define JSON marshaling
}
