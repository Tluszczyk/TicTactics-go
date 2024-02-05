package types

type ErrorResponse struct {
	Title   string `json:"title"`
	Details string `json:"details"`
}

type Request struct {
	Body       string `json:"body"`
	HTTPMethod string `json:"httpMethod"`
}

type Response struct {
	Body       interface{} `json:"body"`
	StatusCode int         `json:"statusCode"`
}

type Code int

const (
	OK Code = iota
	ERR
)

type Status struct {
	Code    Code
	Message string
}

type Credentials struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

type Token string
type Session struct {
	Token Token
	UID   UserID
}

type SessionID string
type UserSessionMapping struct {
	UID    UserID
	SID    SessionID
}

type UserID string
type User struct {
	UID      UserID
	Username string
	Email    string
	Elo      int
}

type Game struct {
	GID            string `json:"gid"`
	UID1           string `json:"uid1"`
	UID2           string `json:"uid2"`
	Board          string `json:"board"`
	Turn           string `json:"turn"`
	Winner         string `json:"winner"`
	MoveHistory    string `json:"move_history"`
	AvailableMoves string `json:"available_moves"`
}

type GameSettings struct{}

type GameFilter struct {
	OPlayerUsername string `json:"oplayer_username"`
	XPlayerUsername string `json:"xplayer_username"`
	Turn            string `json:"turn"`
	Winner          string `json:"winner"`
}
