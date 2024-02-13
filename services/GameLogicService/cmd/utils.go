package cmd

import (
	"fmt"
	"services/lib/types"
)

func (gameLogicService GameLogicService) CreateBoard() Board {
	board := Board{
		Body: [][]Tile{},
	}

	for tileRowID := 0; tileRowID < 3; tileRowID++ {
		tileRow := []Tile{}

		for tileID := 0; tileID < 3; tileID++ {
			tile := Tile{
				Row:  tileRowID,
				Col:  tileID,
				Body: [][]Cell{},
			}

			for cellRowID := 0; cellRowID < 3; cellRowID++ {
				cellRow := []Cell{}

				for cellID := 0; cellID < 3; cellID++ {
					cell := Cell{
						Row:  cellRowID,
						Col:  cellID,
						Body: Empty,
					}

					cellRow = append(cellRow, cell)
				}

				tile.Body = append(tile.Body, cellRow)
			}

			tileRow = append(tileRow, tile)
		}

		board.Body = append(board.Body, tileRow)
	}

	return board
}

func (gameLogicService GameLogicService) CreateTileBoard() TileBoard {
	return TileBoard{
		Body: [][]types.GameWinner{
			{types.NONE, types.NONE, types.NONE},
			{types.NONE, types.NONE, types.NONE},
			{types.NONE, types.NONE, types.NONE},
		},
	}
}

func (gameLogicService GameLogicService) GetDefaultGameSettings() types.GameSettings {
	return types.GameSettings{}
}

func getCellPosition(cellID types.CellPosition) (int, int) {
	return int(cellID[0]) - 65, int(cellID[1]) - 49
}

func (tile Tile) GetCellID(cell Cell) types.CellPosition {
	return types.CellPosition(string(rune(3*tile.Row + cell.Row + 65)) + string(rune(3*tile.Col + cell.Col + 49)))
}

func getAllFreeCells(board Board) []types.CellPosition {
	freeCells := []types.CellPosition{}

	for _, tileRow := range board.Body {
		for _, tile := range tileRow {
			for _, cellRow := range tile.Body {
				for _, cell := range cellRow {
					if cell.Body == Empty {
						freeCells = append(freeCells, tile.GetCellID(cell))
					}
				}
			}
		}
	}

	return freeCells
}

func getTileAvailableMoves(tile Tile, lastTileRow int, lastTileCol int) []types.CellPosition {
	freeCells := []types.CellPosition{}

	for row, cellRow := range tile.Body {
		for col, cell := range cellRow {
			if row == lastTileRow && col == lastTileCol {
				continue
			}

			if cell.Body == Empty {
				freeCells = append(freeCells, tile.GetCellID(cell))
			}
		}
	}

	return freeCells
}

func (gameLogicService GameLogicService) UpdateAvailableMoves(game *types.Game) {
	board := gameLogicService.DeserializeBoard(game.Board)

	// If the game has not started yet, all cells are available
	if len(game.MoveHistory) == 0 {
		game.AvailableMoves = getAllFreeCells(board)
		return
	}

	lastMove := game.MoveHistory[len(game.MoveHistory)-1]

	lastMoveRow, lastMoveCol := getCellPosition(lastMove)
	lastTileRow, lastTileCol := lastMoveRow/3, lastMoveCol/3
	nextTileRow, nextTileCol := lastMoveRow%3, lastMoveCol%3

	tileAvailableMoves := getTileAvailableMoves(board.Body[nextTileRow][nextTileCol], lastTileRow, lastTileCol)

	// If the tile is full, all cells are available
	if len(tileAvailableMoves) == 0 {
		game.AvailableMoves = getAllFreeCells(board)
		return
	}

	game.AvailableMoves = tileAvailableMoves
}

func (tile Tile) IsFull() bool {
	for _, cellRow := range tile.Body {
		for _, cell := range cellRow {
			if cell.Body == Empty {
				return false
			}
		}
	}

	return true
}

func (tile Tile) CheckWinner() types.GameWinner {
	for row := 0; row < 3; row++ {
		if tile.Body[row][0].Body != Empty && tile.Body[row][0].Body == tile.Body[row][1].Body && tile.Body[row][1].Body == tile.Body[row][2].Body {
			return types.GameWinner(tile.Body[row][0].Body)
		}
	}

	for col := 0; col < 3; col++ {
		if tile.Body[0][col].Body != Empty && tile.Body[0][col].Body == tile.Body[1][col].Body && tile.Body[1][col].Body == tile.Body[2][col].Body {
			return types.GameWinner(tile.Body[0][col].Body)
		}
	}

	if tile.Body[0][0].Body != Empty && tile.Body[0][0].Body == tile.Body[1][1].Body && tile.Body[1][1].Body == tile.Body[2][2].Body {
		return types.GameWinner(tile.Body[0][0].Body)
	}

	if tile.Body[0][2].Body != Empty && tile.Body[0][2].Body == tile.Body[1][1].Body && tile.Body[1][1].Body == tile.Body[2][0].Body {
		return types.GameWinner(tile.Body[0][2].Body)
	}

	if tile.IsFull() {
		return types.GameWinner(types.TIE)
	}

	return types.NONE
}

func (tileBoard TileBoard) IsFull() bool {
	for _, row := range tileBoard.Body {
		for _, cell := range row {
			if cell == types.NONE {
				return false
			}
		}
	}

	return true
}

func (tileBoard TileBoard) CheckWinner() types.GameWinner {
	for row := 0; row < 3; row++ {
		if tileBoard.Body[row][0] != types.NONE && tileBoard.Body[row][0] == tileBoard.Body[row][1] && tileBoard.Body[row][1] == tileBoard.Body[row][2] {
			return tileBoard.Body[row][0]
		}
	}

	for col := 0; col < 3; col++ {
		if tileBoard.Body[0][col] != types.NONE && tileBoard.Body[0][col] == tileBoard.Body[1][col] && tileBoard.Body[1][col] == tileBoard.Body[2][col] {
			return tileBoard.Body[0][col]
		}
	}

	if tileBoard.Body[0][0] != types.NONE && tileBoard.Body[0][0] == tileBoard.Body[1][1] && tileBoard.Body[1][1] == tileBoard.Body[2][2] {
		return tileBoard.Body[0][0]
	}

	if tileBoard.Body[0][2] != types.NONE && tileBoard.Body[0][2] == tileBoard.Body[1][1] && tileBoard.Body[1][1] == tileBoard.Body[2][0] {
		return tileBoard.Body[0][2]
	}

	if tileBoard.IsFull() {
		return types.TIE
	}

	return types.NONE
}

func (gameLogicService GameLogicService) UpdateTileBoard(game *types.Game) {
	board := gameLogicService.DeserializeBoard(game.Board)
	tileBoard := TileBoard{}

	for _, tileRow := range board.Body {
		row := []types.GameWinner{}

		for _, tile := range tileRow {
			row = append(
				row, types.GameWinner(tile.CheckWinner()),
			)
		}

		tileBoard.Body = append(tileBoard.Body, row)
	}

	game.TileBoard = gameLogicService.SerializeTileBoard(tileBoard)
}

func (gameLogicService GameLogicService) UpdateWinner(game *types.Game) {
	tileBoard := gameLogicService.DeserializeTileBoard(game.TileBoard)

	winner := tileBoard.CheckWinner()

	if winner != types.NONE {
		game.Winner = string(winner)
		return
	}

	if len(game.AvailableMoves) == 0 {
		game.Winner = string(types.TIE)
		return
	}
}

func (gameLogicService GameLogicService) UpdateGameState(game *types.Game) {
	if game.Winner != string(types.NONE) {
		game.State = types.FINISHED
		return
	}
}

func (gameLogicService GameLogicService) UpdateTurn(game *types.Game) {
	if game.State == types.IN_PROGRESS {
		if game.Turn == string(types.X) {
			game.Turn = string(types.O)
		} else {
			game.Turn = string(types.X)
		}
	}
}

func (gameLogicService GameLogicService) UpdateGame(game *types.Game) {
	gameLogicService.UpdateAvailableMoves(game)
	gameLogicService.UpdateTileBoard(game)
	gameLogicService.UpdateWinner(game)
	gameLogicService.UpdateGameState(game)
	gameLogicService.UpdateTurn(game)
}

func (gameLogicService GameLogicService) CanPutMove(game *types.Game, move types.Move) (bool, error) {
	if game.Turn != move.Symbol {
		return false, fmt.Errorf("it's not %s's turn", move.Symbol)
	}

	if game.State != types.IN_PROGRESS {
		return false, fmt.Errorf("the game is not in progress")
	}

	if len(game.AvailableMoves) == 0 {
		return false, fmt.Errorf("there are no available moves")
	}

	for _, availableMove := range game.AvailableMoves {
		if availableMove == move.Cell {
			return true, nil
		}
	}

	return false, fmt.Errorf("the cell is not available")
}

func (gameLogicService GameLogicService) PutMove(game *types.Game, move types.Move) {
	board := gameLogicService.DeserializeBoard(game.Board)

	cellRow, cellCol := getCellPosition(move.Cell)

	tileRow, tileCol := cellRow/3, cellCol/3

	board.Body[tileRow][tileCol].Body[cellRow%3][cellCol%3].Body = CellValue(move.Symbol)

	game.Board = gameLogicService.SerializeBoard(board)
	game.MoveHistory = append(game.MoveHistory, move.Cell)

	gameLogicService.UpdateGame(game)
}
