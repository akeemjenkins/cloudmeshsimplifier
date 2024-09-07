package mesh

import (
	"errors"
	"fmt"
	"log"
	"testing"

	"bitbucket.org/cloudcomputer/cloud-mesh-simplifier/auxmath"
	"bitbucket.org/cloudcomputer/cloud-mesh-simplifier/cloudmesh"
)

//1.  Create single triangle mesh, check that it has an empty neighbor list with no error
//2.  Create a basic octahedron, check specific neighbors
//3.  Create a triStrip, check that all interior triangles have exactly two neighbors

//TODO: we need to refactor the CloudMesh stuff so that it's accessible in test modules without creating circular dependencies

type myIndexedMesh struct {
	Indices  []uint32
	Vertices []float32
}

func (m myIndexedMesh) GetNumFacets() uint32 {
	return uint32(len(m.Indices) / 3)
}

func (m myIndexedMesh) GetNumVertices() uint32 {
	return uint32(len(m.Vertices) / 3)
}

func (m myIndexedMesh) GetVertices(triangle uint32) ([]uint32, error) {

	if triangle >= m.GetNumFacets() {
		return []uint32{0, 0, 0}, errors.New("GetVertices:requested index is out of bounds")
	}
	return []uint32{m.Indices[3*triangle+0],
		m.Indices[3*triangle+1],
		m.Indices[3*triangle+2]}, nil
}

func (m myIndexedMesh) GetPoint(vertex uint32) ([]float32, error) {
	if vertex > m.GetNumVertices() {
		return []float32{0, 0, 0}, errors.New("GetPoint:requested index is out of bounds")
	}

	return []float32{m.Vertices[3*vertex+0],
		m.Vertices[3*vertex+1],
		m.Vertices[3*vertex+2]}, nil
}

func TestGetNeighbors1(t *testing.T) {
	theseVertices := []float32{0, 0, 0, 0, 0, 1, 0, 1, 0}
	theseTriangles := []uint32{0, 1, 2}
	testMesh := myIndexedMesh{Indices: theseTriangles, Vertices: theseVertices}

	neighborhood := CreateNeighborhood(testMesh)

	neighbs, err := neighborhood.GetTriangleNeighborsOfTriangle(0)

	if len(neighbs) != 0 {
		t.Error("Expected 0 neighbors in a mesh with only one triangle.")
	}

	if err != nil {
		t.Error("Expected no error for empty neighbor list.")
	}
}

func createOctahedronMesh() Mesh {
	theseVertices := []float32{
		0, 0, 1,
		1, 0, 0,
		0, 1, 0,
		-1, 0, 0,
		0, -1, 0,
		0, 0, -1}

	theseTriangles := []uint32{
		1, 2, 0,
		2, 3, 0,
		3, 4, 0,
		4, 1, 0,
		5, 1, 4,
		5, 2, 1,
		5, 3, 2,
		5, 4, 3}
	fullTriangles := make([]uint32, 8*3)
	copy(fullTriangles, theseTriangles)

	fullVertices := make([]float32, 6*3)
	copy(fullVertices, theseVertices)

	return myIndexedMesh{Indices: fullTriangles, Vertices: fullVertices}
}

func TestGetNeighbors2(t *testing.T) {
	octahedron := createOctahedronMesh()
	neighborhood := CreateNeighborhood(octahedron)

	neighbs0, _ := neighborhood.GetTriangleNeighborsOfTriangle(0)
	neighbs1, _ := neighborhood.GetTriangleNeighborsOfTriangle(1)
	neighbs2, _ := neighborhood.GetTriangleNeighborsOfTriangle(2)

	fmt.Printf("neighbors0: %d, %d, %d \n", neighbs0[0], neighbs0[1], neighbs0[2])
	fmt.Printf("neighbors1: %d, %d, %d \n", neighbs1[0], neighbs1[1], neighbs1[2])
	fmt.Printf("neighbors2: %d, %d, %d \n", neighbs2[0], neighbs2[1], neighbs2[2])

	if len(neighbs0) != 3 {
		t.Error("Expected number of neighbors to be 3 for all triangle in octahedron.")
	}
	if len(neighbs1) != 3 {
		t.Error("Expected number of neighbors to be 3 for all triangle in octahedron.")
	}
	if len(neighbs2) != 3 {
		t.Error("Expected number of neighbors to be 3 for all triangle in octahedron.")
	}

	if neighbs0[0] != uint32(5) {
		t.Error("Expected first neighbor of triangle 0 to be triangle 5")
	}

	if neighbs0[1] != uint32(1) {
		t.Error("Expected second neighbor of triangle 0 to be triangle 1")
	}

	if neighbs0[2] != uint32(3) {
		t.Error("Expected third neighbor of triangle 0 to be triangle 3")
	}

	if neighbs1[0] != uint32(6) {
		t.Error("Expected first neighbor of triangle 1 to be triangle 6")
	}

	if neighbs1[1] != uint32(2) {
		t.Error("Expected second neighbor of triangle 1 to be triangle 2")
	}

	if neighbs1[2] != uint32(0) {
		t.Error("Expected third neighbor of triangle 1 to be triangle 0")
	}

	if neighbs2[0] != uint32(7) {
		t.Error("Expected first neighbor of triangle 2 to be triangle 7")
	}

	if neighbs2[1] != uint32(3) {
		t.Error("Expected second neighbor of triangle 2 to be triangle 3")
	}

	if neighbs2[2] != uint32(1) {
		t.Error("Expected third neighbor of triangle 2 to be triangle 1")
	}

}

func TestGetNormal1(t *testing.T) {
	p1 := []float32{0, 0, 0}
	p2 := []float32{2, 0, 0}
	p3 := []float32{0, 1, 0}

	normal, errNormal := ComputeTriangleNormal(p1, p2, p3)
	if errNormal != nil {
		log.Fatal("Expected there to be no problem with computing triangle normal")
	}
	expectedAnswer := []float32{0, 0, 1}
	dot, err := auxmath.Dot(expectedAnswer, normal)
	if err != nil {
		log.Fatalf("Problem with dot product: %s", err)
	}
	if dot < .999 || dot > 1.001 {
		t.Error("Expected the normal to be in the z-direction with magnitude 1 ")
	}
}

func TestGetNormal2(t *testing.T) {
	p1 := []float32{0, 0, 0}
	p2 := []float32{0, 1, 0}
	p3 := []float32{2, 0, 0}

	normal, errNormal := ComputeTriangleNormal(p1, p2, p3)
	if errNormal != nil {
		log.Fatal("Expected there to be no problem with computing triangle normal")
	}
	expectedAnswer := []float32{0, 0, -1}
	dot, err := auxmath.Dot(expectedAnswer, normal)
	if err != nil {
		log.Fatalf("Problem with dot product: %s", err)
	}
	if dot < .999 || dot > 1.001 {
		t.Error("Expected the normal to be in the negative z-direction with magnitude 1 ")
	}
}

func TestGetNormal3(t *testing.T) {
	p1 := []float32{0, 0, 10}
	p2 := []float32{0, 0, 1}
	p3 := []float32{0, 0, 5}

	_, errNormal := ComputeTriangleNormal(p1, p2, p3)
	if errNormal == nil {
		t.Error("Expected there to be an error for a degenerate triangle:")
	}
}

func TestGetArea1(t *testing.T) {
	p1 := []float32{0, 0, 1}
	p2 := []float32{0, 1, 0}
	p3 := []float32{0, 0, 0}

	area := ComputeTriangleArea(p1, p2, p3)
	if area > .50001 || area < .4999 {
		t.Error("Expected area to be exactly .5.")
	}
}

// TODO: fix a set of triangles that are out of bounds for their Vertices
// serialize Indexed Mesh should return some sort of error with out of bounds
func TestSerializeIndexedMesh(t *testing.T) {
	theseVertices := []float32{1, 2, 3, 4, 5, 6, 7, 8, 9}
	theseTriangles := []uint32{0, 1, 2}
	testMesh := cloudmesh.IndexedMesh{Indices: theseTriangles, Vertices: theseVertices}
	retVal := SerializeIndexedMesh(testMesh)
	// here we reuse theseVertices to compare againts the retVal
	for i := range retVal {
		if theseVertices[i] != retVal[i] {
			t.Errorf("Expected %v, got %v\n", theseVertices[i], retVal[i])
		}
	}
}

func TestComputeArea(t *testing.T) {
	theseVertices := []float32{0, 0, 0, 0, 1, 0, 1, 1, 0}
	theseTriangles := []uint32{0, 1, 2}
	testMesh := cloudmesh.IndexedMesh{Indices: theseTriangles, Vertices: theseVertices}
	retVal := ComputeArea(testMesh, 0)
	if retVal != .5 {
		t.Errorf("Expected %v, got %v\n", .5, retVal)
	}
}

func TestComputeCentroid(t *testing.T) {
	theseVertices := []float32{0, 0, 0, 0, 3, 0, 3, 3, 0}
	theseTriangles := []uint32{0, 1, 2}
	testMesh := cloudmesh.IndexedMesh{Indices: theseTriangles, Vertices: theseVertices}
	retVal := ComputeCentroid(testMesh, 0)
	center := []float32{1, 2, 0}
	for i := range center {
		if center[i] != retVal[i] {
			t.Errorf("Expected %v, got %v\n", center[i], retVal[i])
		}
	}
	//retVal = ComputeCentroid(testMesh, 1)
}
