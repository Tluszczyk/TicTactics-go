package types

type ErrorResponse struct {
	Title   string `json:"title"`
	Details string `json:"details"`
}

type Request struct {
	Body       string `json:"body"`
	Path       string `json:"path"`
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
	UID UserID
	SID SessionID
}

type UserID string
type User struct {
	UID      UserID
	Username string
	Email    string
	Elo      int
}

type GameState string

const (
	WATING_FOR_OPPONENT GameState = "WAITING_FOR_OPPONENT"
	IN_PROGRESS         GameState = "IN_PROGRESS"
	FINISHED            GameState = "FINISHED"
)

type GameID string
type Game struct {
	GID            GameID    `json:"gid"`
	Board          string    `json:"board"`
	Turn           string    `json:"turn"`
	Winner         string    `json:"winner"`
	MoveHistory    []string  `json:"move_history"`
	AvailableMoves []string  `json:"available_moves"`
	State          GameState `json:"state"`
	TileBoard      string    `json:"tile_board"`
}

type GameSettings struct{}

type GameFilter struct {
	UID    UserID    `json:"uid"`
	Turn   string    `json:"turn"`
	Winner string    `json:"winner"`
	State  GameState `json:"state"`
}

func (filter *GameFilter) Check(game Game) bool {
	return ((filter.State == "" || filter.State == game.State) &&
		(filter.Winner == "" || filter.Winner == game.Winner) &&
		(filter.Turn == "" || filter.Turn == game.Turn))
}

type GameWinner string

const (
	X    GameWinner = "X"
	O    GameWinner = "O"
	TIE  GameWinner = "T"
	NONE GameWinner = "."
)
