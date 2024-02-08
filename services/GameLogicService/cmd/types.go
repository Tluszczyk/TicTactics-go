package cmd

import (
	"services/lib/types"
)

type CellValue string

const (
	Empty CellValue = "."
	X     CellValue = "X"
	O     CellValue = "O"
)

type Cell struct {
	Row  int
	Col  int
	Body CellValue
}

type Tile struct {
	Row  int
	Col  int
	Body [][]Cell
}

type Board struct {
	Body [][]Tile
}

type TileBoard struct {
	Body [][]types.GameWinner
}