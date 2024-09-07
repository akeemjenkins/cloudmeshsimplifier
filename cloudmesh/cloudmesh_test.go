package cloudmesh

import (
	"log"
	"testing"
)

func TestGetNumTriangles(t *testing.T) {
	//Indices1:= []uint32{1,1,1}
	//Vertices1:= []float32{1,1,1}

	theseVertices := []float32{0, 0, 1}
	theseTriangles := []uint32{2, 1, 0}
	testMesh := IndexedMesh{Indices: theseTriangles, Vertices: theseVertices}
	numTriangles := testMesh.GetNumFacets()
	if numTriangles != 1 {
		t.Error("Expected 1 triangle and got ", numTriangles)
	}
}

//check for Vertices that aren't divisible by 3 only
func TestGetNumVertices(t *testing.T) {
	theseVertices := []float32{0, 0, 1}
	theseTriangles := []uint32{2, 1, 0}
	testMesh := IndexedMesh{Indices: theseTriangles, Vertices: theseVertices}
	numTriangles := testMesh.GetNumVertices()
	if numTriangles != 1 {
		t.Error("Expected 1 vertex and got ", numTriangles)
	}
}

func TestGetVertices(t *testing.T) {
	theseVertices := []float32{0, 0, 1, 1, 2, 111, 5, 6, 7}
	theseTriangles := []uint32{2, 1, 9, 2, 88, 22, 8, 9, 10}
	testMesh := IndexedMesh{Indices: theseTriangles, Vertices: theseVertices}
	triangles, err := testMesh.GetVertices(1)
	if err != nil {
		log.Fatal(err)
	}
	index := 3
	for i := range triangles {
		if triangles[i] != testMesh.Indices[index] {
			t.Errorf("Expected %v and got %v", testMesh.Indices[index], triangles[i])
		}
		index++
	}
	// test case where requested triangle is out of bounds
	triangles, err = testMesh.GetVertices(60)
	if err == nil {
		t.Errorf("Should have returned an error as the index is out of bounds.")
	}
}

func TestGetPoint(t *testing.T) {
	theseVertices := []float32{0, 0, 1, 1, 2, 111, 5, 6, 7}
	theseTriangles := []uint32{2, 1, 9, 2, 2, 2, 8, 9, 10}
	testMesh := IndexedMesh{Indices: theseTriangles, Vertices: theseVertices}
	point, err := testMesh.GetPoint(1)
	if err != nil {
		log.Fatal(err)
	}
	index := 3
	for i := range point {
		if point[i] != testMesh.Vertices[index] {
			t.Errorf("Expected %v and got %v", testMesh.Vertices[index], point[i])
		}
		index++
	}
	// test where requested point is out of bounds
	point, err = testMesh.GetPoint(50)
	if err == nil {
		t.Errorf("Expected an error as the point requested is out of bounds.")
	}
}

func TestRemoveTriangle(t *testing.T) {
	theseVertices := []float32{0, 0, 1, 1, 1, 1, 5, 6, 7}
	theseTriangles := []uint32{2, 1, 9, 2, 2, 2, 8, 9, 10}
	testMesh := IndexedMesh{Indices: theseTriangles, Vertices: theseVertices}
	testMesh.RemoveTriangle(1)
	if testMesh.Indices[2] != 9 {
		t.Error("Expected 9 at index 2 and got: ", testMesh.Indices[2])
	}
	if testMesh.Indices[3] != 8 {
		t.Error("Expcted 8 at index 3 and got: ", testMesh.Indices[3])
	}
}

func TestPop(t *testing.T) {
	theseVertices := []float32{0, 0, 1, 1, 1, 1, 5, 6, 7}
	theseTriangles := []uint32{2, 1, 9, 2, 2, 2, 8, 9, 10}
	testMesh := IndexedMesh{Indices: theseTriangles, Vertices: theseVertices}
	testMesh.Pop()
	for i := range testMesh.Indices {
		if theseTriangles[i+3] != testMesh.Indices[i] {
			t.Errorf("Expected %v and got %v", theseTriangles[i+3], testMesh.Indices[i])
		}
	}
}

func TestAddTriangle(t *testing.T) {
	theseVertices := []float32{0, 0, 1, 1, 1, 1}
	theseTriangles := []uint32{2, 1, 9, 2, 2, 2}
	testMesh := IndexedMesh{Indices: theseTriangles, Vertices: theseVertices}

	testMesh.AddTriangle(1, 10, 11)
	if testMesh.Indices[6] != 1 {
		t.Error("Expected 1, got: ", testMesh.Indices[6])
	}
	if testMesh.Indices[7] != 10 {
		t.Error("Expected 10, got: ", testMesh.Indices[7])
	}
	if testMesh.Indices[8] != 11 {
		t.Error("Expected 11, got: ", testMesh.Indices[8])
	}
}

/*
func TestAddVertex(t *testing.T) {
	theseVertices := []float32{0, 0, 1, 1, 1, 1}
	theseTriangles := []uint32{2, 1, 9, 2, 2, 2}
	testMesh := IndexedMesh{Indices: theseTriangles, Vertices: theseVertices}

	vertexIndex := testMesh.AddVertex(1, 1, 1)
	if vertexIndex != 1 {
		t.Error("Expected 1, got:", vertexIndex)
	}
	vertexIndex = testMesh.AddVertex(0, 1, 0)
	if vertexIndex != 2 {
		t.Error("Expected 2, got:", vertexIndex)
	}
}

func BenchmarkAddVertex(b *testing.B) {
	theseVertices := []float32{0, 0, 1, 1, 1, 1}
	theseTriangles := []uint32{2, 1, 9, 2, 2, 2}
	testMesh := IndexedMesh{Indices: theseTriangles, Vertices: theseVertices}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testMesh.AddVertex(rand.Float32(), 1, 1)
	}
}
*/
func TestLastTriangle(t *testing.T) {
	theseVertices := []float32{0, 0, 1, 1, 1, 1}
	theseTriangles := []uint32{2, 1, 9, 2, 2, 2}
	testMesh := IndexedMesh{Indices: theseTriangles, Vertices: theseVertices}

	lastTri := testMesh.LastTriangle()
	if lastTri != 1 {
		t.Error("got:", lastTri)
	}
}
