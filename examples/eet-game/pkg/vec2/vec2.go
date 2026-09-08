// Package vec2 provides a simple 2D vector type and common math utilities.
package vec2

import "math"

// Vec2 represents a 2D vector with X and Y components.
type Vec2 struct {
	X float64
	Y float64
}

// Add returns the sum of v and other.
func (v Vec2) Add(other Vec2) Vec2 {
	return Vec2{X: v.X + other.X, Y: v.Y + other.Y}
}

// Sub returns v minus other.
func (v Vec2) Sub(other Vec2) Vec2 {
	return Vec2{X: v.X - other.X, Y: v.Y - other.Y}
}

// Scale returns v multiplied by scalar s.
func (v Vec2) Scale(s float64) Vec2 {
	return Vec2{X: v.X * s, Y: v.Y * s}
}

// Length returns the magnitude of v.
func (v Vec2) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// Distance returns the Euclidean distance between v and other.
func (v Vec2) Distance(other Vec2) float64 {
	return v.Sub(other).Length()
}

// Normalize returns a unit vector in the direction of v.
// Returns Vec2{0,0} if the length is zero.
func (v Vec2) Normalize() Vec2 {
	l := v.Length()
	if l == 0 {
		return Vec2{}
	}
	return v.Scale(1 / l)
}
