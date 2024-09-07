package auxmath

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

func TestDot(t *testing.T) {
	a := []float32{1, 1, 1}
	b := []float32{1, 1, 1}
	result, _ := Dot(a, b)
	if result != 3 {
		t.Error("Expected 1 and got ", result)
	}

	// This is the limit of a 32 bit float in multiplying
	// 1.844^2 will overflow a float32
	a = []float32{1.8446744e+19, 1.8446744e+19, 1.8446744e+19}
	b = []float32{1.8446744e+19, 1.8446744e+19, 1.8446744e+19}
	result1, err1 := Dot(a, b)
	if err1 == nil && result1 == float32(math.Inf(1)) {
		t.Errorf("Failed to make sure the dot product wouldn't be infinite.")
	}

	//
	a = []float32{1, 2, 3}
	b = []float32{1}
	result2, err2 := Dot(a, b)
	if err2 != nil && result2 != float32(0.0) {
		t.Errorf("Failed to detect different sized slices")
	}
}

func BenchmarkDot(b *testing.B) {
	var a []float32
	var c []float32
	for i := 0; i < 200000000; i++ {
		a = append(a, rand.Float32())
		c = append(c, rand.Float32())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Dot(a, c)
	}
}

func TestNormalize(t *testing.T) {
	v := []float32{1, 0, 0}
	result := Normalize(v)
	expected := []float32{1, 0, 0}
	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("Expected, %v and got %v", expected[i], result[i])
		}
	}
	z := []float32{1e-7, 1e-7, 1e-7}
	result = Normalize(z)
	expectedz := make([]float32, len(z))
	for index := range result {
		if result[index] != expectedz[index] {
			t.Errorf("Expected, %v and got %v", expectedz[index], result[index])
		}
	}
}

func TestRandomUint32(t *testing.T) {
	t1 := RandomUint32(0, 0)
	if t1 != 0 {
		t.Error("Expected 0 and got %v\n", t1)
	}
	t1 = RandomUint32(9, 4)
	if t1 != 0 {
		t.Error("Failed passing in min > max\n")
	}
	t1 = RandomUint32(0, 5000)
	if t1 < 0 || t1 > 5000 {
		t.Error("The random number generated was outside of the bounds given: %v\n", t1)
	}
}

func TestCross(t *testing.T) {
	u := []float32{float32(0), float32(0), float32(0)}
	c := Cross(u, u)
	fmt.Printf("c: %v", c)
}
