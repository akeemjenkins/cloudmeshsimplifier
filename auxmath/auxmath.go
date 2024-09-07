// Package auxmath contains math functions that are not part of the standard
// golang math library
package auxmath

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"
)

// Dot : Return the dot product of two []float32
// This will return 0 and an error if the length of the two slices are wrong.
// TODO: benchmark Dot
func Dot(a []float32, b []float32) (f float32, err error) {
	sum := float32(0.0)
	if len(a) != len(b) {
		return sum, fmt.Errorf("auxmath->Dot: The size of the arrays are not the same")
	}
	for i := range a {
		if a[i] >= 1.8446744e+19 && b[i] >= 1.8446744e+19 {
			return float32(math.Inf(1)), errors.New("Dot Product would be infinite")
		}
		sum += a[i] * b[i]
	}
	return sum, nil
}

/* returns a - b*/
func Subtract(a [] float32, b[]float32)(c []float32, err error){
	if len(a) != len(b){
		return []float32{0,0,0}, errors.New("Subtract: vector lengths not equal")
	}
	retVal := make([]float32, len(a))
	for i:= range a{
		retVal[i] = a[i] - b[i]
	}
	return retVal, nil
}

func Add(a [] float32, b[]float32)(c []float32, err error){
	if len(a) != len(b){
		return []float32{0,0,0}, errors.New("Add: vector lengths not equal")
	}
	retVal := make([]float32, len(a))
	for i:= range a{
		retVal[i] = a[i] + b[i]
	}
	return retVal, nil
}

func Scale(v []float32, s float32) []float32{
	retVal := make([]float32, len(v))
	for i:=range v{
		retVal[i] = s*v[i]
	}
	return retVal
}

//Cross returns the cross product between two 3D vectors u x v
func Cross(u []float32, v []float32) []float32 {
	return []float32{u[1]*v[2] - u[2]*v[1], u[2]*v[0] - u[0]*v[2], u[0]*v[1] - u[1]*v[0]}
}

//Magnitude returns the magnitude of the input vector
func Magnitude(u []float32) float32 {
	normSq, _ := Dot(u, u)
	if normSq < 1e-8 {
		return 0
	}
	return float32(math.Sqrt(float64(normSq)))
}

//Parallel determines if two vectors are parallel according to the input tolerance,
//which is the max allowable cosine of the angle between the vectors
func Parallel(u []float32, v []float32, tol float32) bool {
	dot, _ := Dot(u, v)
	deviation := dot - Magnitude(u)*Magnitude(v)
	if math.Abs(float64(deviation)) > float64(tol) {
		return false
	}
	return true

}

// Normalize normalizes a slice of float32
func Normalize(v []float32) []float32 {
	retVal := make([]float32, len(v))
	norm := Magnitude(v)
	if norm < 1e-6 {
		return retVal
	}
	for i := 0; i < len(v); i++ {
		retVal[i] = v[i] / norm
	}
	return retVal
}

// RandomUint32 generates a random uint32 based on a min and max.
// If the min is greater than the max, we will just return 0. Try better next time.
func RandomUint32(min, max uint32) uint32 {
	if min > max {
		return uint32(0)
	}
	if min == max {
		return min
	}
	rand.Seed(time.Now().UnixNano())
	return uint32(rand.Intn(int(max)-int(min))) + min
}
