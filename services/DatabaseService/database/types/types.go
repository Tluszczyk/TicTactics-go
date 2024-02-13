package types

type FieldType string

const (
	CREATED_AT FieldType = "createdAt"

	USERNAME      FieldType = "USER"
	UID           FieldType = "UID"
	EMAIL         FieldType = "EMAIL"
	HASH_ID       FieldType = "HID"
	PASSWORD_HASH FieldType = "HASH"
	ELO           FieldType = "ELO"

	SID   FieldType = "SID"
	TOKEN FieldType = "TOKEN"

	GID             FieldType = "GID"
	GAME_SETTINGS   FieldType = "GAME_SETTINGS"
	BOARD           FieldType = "BOARD"
	TURN            FieldType = "TURN"
	WINNER          FieldType = "WINNER"
	MOVE_HISTORY    FieldType = "MOVE_HISTORY"
	AVAILABLE_MOVES FieldType = "AVAILABLE_MOVES"
	STATE           FieldType = "STATE"
	TILE_BOARD      FieldType = "TILE_BOARD"

	SYMBOL FieldType = "SYMBOL"
)

type DatabaseItem struct {
	PK         map[FieldType]interface{}
	SK         map[FieldType]interface{}
	Attributes map[FieldType]interface{}
}

func NewDatabaseItem() *DatabaseItem {
	return &DatabaseItem{
		PK:         make(map[FieldType]interface{}),
		SK:         make(map[FieldType]interface{}),
		Attributes: make(map[FieldType]interface{}),
	}
}

func (d *DatabaseItem) IsNil() bool {
	return len(d.PK) == 0 && len(d.SK) == 0 && len(d.Attributes) == 0
}

type DatabaseQuery DatabaseItem

// Inputs
type DatabaseGetItemInput struct {
	TableName string
	Key       DatabaseItem
}

type DatabaseGetItemsInput struct {
	TableName string
	Key       DatabaseItem
}

type DatabasePutItemInput struct {
	TableName string
	Item      DatabaseItem
}

type DatabaseDeleteItemInput struct {
	TableName string
	Key       DatabaseItem
}

type DatabaseUpdateItemInput struct {
	TableName string
	Key       DatabaseItem
	Item      DatabaseItem
}

type DatabaseQueryInput struct {
	TableName string
	Query     DatabaseQuery
}

// Outputs
type DatabaseGetItemOutput struct {
	Item DatabaseItem
}

type DatabaseGetItemsOutput struct {
	Items []DatabaseItem
}

type DatabasePutItemOutput struct {
}

type DatabaseDeleteItemOutput struct {
}

type DatabaseUpdateItemOutput struct {
}

type DatabaseQueryOutput struct {
	Items []DatabaseItem
}
