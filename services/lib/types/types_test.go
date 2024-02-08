package types

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
)

func TestErrorResponse(t *testing.T) {
	errResp := ErrorResponse{
		Title:   "Error",
		Details: "Something went wrong",
	}

	jsonData, err := json.Marshal(errResp)
	if err != nil {
		t.Errorf("Error marshaling ErrorResponse: %v", err)
	}

	expectedJSON := `{"title":"Error","details":"Something went wrong"}`

	if string(jsonData) != expectedJSON {
		t.Errorf("Expected JSON: %s, got: %s", expectedJSON, string(jsonData))
	}
}

func TestSession(t *testing.T) {
	userID := UserID(uuid.New().String())
	user := User{
		UID:      userID,
		Username: "testuser",
	}

	token := Token(uuid.New().String())

	session := Session{
		Token: token,
		UID:   user.UID,
	}

	if session.UID != userID {
		t.Errorf("Expected UserID: %v, got: %v", userID, session.UID)
	}

	if session.Token != token {
		t.Errorf("Expected Token: %v, got: %v", token, session.Token)
	}
}

func TestGame(t *testing.T) {
	game := Game{
		GID:            "123",
		UID1:           "user1",
		UID2:           "user2",
		Board:          "XO_OX___",
		Turn:           "X",
		Winner:         "",
		MoveHistory:    []string{"A1","B2","C3"},
		AvailableMoves: []string{"A1","B2","C3"},
	}

	jsonData, err := json.Marshal(game)
	if err != nil {
		t.Errorf("Error marshaling Game: %v", err)
	}

	expectedJSON := `{"gid":"123","uid1":"user1","uid2":"user2","board":"XO_OX___","turn":"X","winner":"","move_history":"012345678","available_moves":"3,4,6,7,8"}`

	if string(jsonData) != expectedJSON {
		t.Errorf("Expected JSON: %s, got: %s", expectedJSON, string(jsonData))
	}
}
