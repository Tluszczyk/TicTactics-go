package cmd

import (
	"strings"
	"services/lib/types"
)

// How the board is serialized:
// 1. The board is a 3x3 grid of tiles
// 2. Each tile is a 3x3 grid of cells
// 3. Each cell has a value of either "empty", "X", or "O"

// The board llok like this:

// A1|A2|A3 A4|A5|A6 A7|A8|A9 
// B1|B2|B3 B4|B5|B6 B7|B8|B9 
// C1|C2|C3 C4|C5|C6 C7|C8|C9 

// D1|D2|D3 D4|D5|D6 D7|D8|D9 
// E1|E2|E3 E4|E5|E6 E7|E8|E9 
// F1|F2|F3 F4|F5|F6 F7|F8|F9
 
// G1|G2|G3 G4|G5|G6 G7|G8|G9 
// H1|H2|H3 H4|H5|H6 H7|H8|H9 
// I1|I2|I3 I4|I5|I6 I7|I8|I9

// The board is serialized by concatenating the cells of each tile, then the tiles of the board, like this:
// A1A2A3B1B2B3C1C2C3A4A5A6B4B5B6C4C5C6A7A8A9B7B8B9C7C8C9D1D2D3E1E2E3F1F2F3D4D5D6E4E5E6F4F5F6D7D8D9E7E8E9F7F8F9G1G2G3H1H2H3I1I2I3G4G5G6H4H5H6I4I5I6G7G8G9H7H8H9I7I8I9

// The serialized board looks like this:
// XO.XO.XO.XO.XO.XO.XO.XO.XO.XO.XO.XO.XO.XO.XO.XO.XO.XO.XO.XO.XO.XO.XO.XO.XO.XO.XO.

func (gameLogicService GameLogicService) SerializeBoard(board Board) string {
	serializedBoard := strings.Builder{}

	for _, tileRow := range board.Body {
		for _, tile := range tileRow {
			for _, cellRow := range tile.Body {
				for _, cell := range cellRow {
					serializedBoard.WriteString(string(cell.Body))
				}
			}
		}
	}

	return serializedBoard.String()
}

func (gameLogicService GameLogicService) DeserializeBoard(serializedBoard string) Board {
	var board Board

	for y := 0; y < 3; y++ {
		row := []Tile{}

		for x := 0; x < 3; x++ {
			iterator := y * 27 + x * 9
			tile := gameLogicService.DeserializeTile(serializedBoard[iterator : iterator + 9], y, x)
			row = append(row, tile)
		}

		board.Body = append(board.Body, row)
	}

	return board
}

func (gameLogicService GameLogicService) DeserializeTile(serializedTile string, row int, col int) Tile {
	tile := Tile{
		Row: row,
		Col: col,
		Body: [][]Cell{},
	}

	for y := 0; y < 3; y++ {
		row := []Cell{}

		for x := 0; x < 3; x++ {
			cell := Cell{
				Row: y,
				Col: x,
				Body: CellValue(serializedTile[y * 3 + x]),
			}

			row = append(row, cell)
		}

		tile.Body = append(tile.Body, row)
	}

	return tile
}

func (gameLogicService GameLogicService) SerializeTileBoard(tileBoard TileBoard) string {
	serializedTileBoard := strings.Builder{}

	for _, row := range tileBoard.Body {
		for _, cell := range row {
			serializedTileBoard.WriteString(string(cell))
		}
	}

	return serializedTileBoard.String()
}

func (gameLogicService GameLogicService) DeserializeTileBoard(serializedTileBoard string) TileBoard {
	tileBoard := TileBoard{}

	for y := 0; y < 3; y++ {
		row := []types.GameWinner{}

		for x := 0; x < 3; x++ {
			cell := types.GameWinner(serializedTileBoard[y * 3 + x])
			row = append(row, cell)
		}

		tileBoard.Body = append(tileBoard.Body, row)
	}

	return tileBoard
}