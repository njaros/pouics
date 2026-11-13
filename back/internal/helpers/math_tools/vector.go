package mathtools

import "math"

// Vector direction is A -> B
func PointsToVector(xA float64, yA float64, xB float64, yB float64) (float64, float64) {
	return xB - xA, yB - yA
}

// Knowing his direction, decompose a norm to his x and y values.
func NormDecomposition(norm float64, dirX float64, dirY float64) (float64, float64) {
	if dirX == 0 && dirY == 0 {
		return 0, 0
	}
	ratio := norm / math.Sqrt(dirX * dirX + dirY * dirY)
	return ratio * dirX, ratio * dirY
}
