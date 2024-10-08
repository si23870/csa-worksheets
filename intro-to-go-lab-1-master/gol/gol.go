package main

func calculateNextState(p golParams, world [][]byte) [][]byte {
	var alive int
	var newWorld = make([][]byte, p.imageHeight)
	for i := range world {
		newWorld[i] = make([]byte, p.imageWidth)
		for j := range world[i] {
			alive = checkLiveNeighbours(p, world, i, j)
			if world[i][j] == 0 {
				if alive == 3 {
					newWorld[i][j] = 255
				} else {
					newWorld[i][j] = 0
				}
			} else if world[i][j] == 255 {
				if alive < 2 || alive > 3 {
					newWorld[i][j] = 0
				} else {
					newWorld[i][j] = 255
				}
			}
		}
	}
	return newWorld
}

func checkLiveNeighbours(p golParams, world [][]byte, i int, j int) (alive int) {
	var below, left, above, right int
	if i == 0 {
		above = p.imageHeight - 1
	} else {
		above = i - 1
	}
	if i == p.imageHeight-1 {
		below = 0
	} else {
		below = i + 1
	}
	if j == 0 {
		left = p.imageWidth - 1
	} else {
		left = j - 1
	}
	if j == p.imageWidth-1 {
		right = 0
	} else {
		right = j + 1
	}

	if world[above][j] == 255 {
		alive++
	}
	if world[below][j] == 255 {
		alive++
	}
	if world[i][left] == 255 {
		alive++
	}
	if world[i][right] == 255 {
		alive++
	}
	return alive
}

func calculateAliveCells(p golParams, world [][]byte) []cell {
	var aliveCells []cell
	for i := range world {
		for j := range world[i] {
			if world[i][j] == 255 {
				aliveCells = append(aliveCells, cell{j, i})
			}
		}
	}

	return aliveCells
}
