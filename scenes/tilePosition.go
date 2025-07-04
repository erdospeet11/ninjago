// Package scenes provides utility functions for managing game scenes and their components.
package scenes

import "math"

// TilePosition calculates the tile position (column and row) for the given x and y coordinates.
// It uses the tile's width and height to determine the corresponding tile indices.
//
// Parameters:
//  - x: The x coordinate in the game world.
//  - y: The y coordinate in the game world.
//  - tileWidth: The width of a single tile.
//  - tileHeight: The height of a single tile.
//
// Returns:
//  - int: The column index of the tile.
//  - int: The row index of the tile.
func TilePosition(x, y, tileWith, tileHeight float64) (int, int) {

	return int(math.Floor(x / tileWith)), int(math.Floor(y / tileHeight))
}
