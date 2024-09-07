package stl

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"testing"

	"bitbucket.org/cloudcomputer/cloud-mesh-simplifier/cloudmesh"
	"bitbucket.org/cloudcomputer/cloud-mesh-simplifier/mesh"
)

func TestLoadSTLFile(t *testing.T) {

}

func TestSaveBinarySTL(t *testing.T) {

}

func TestCreateMesh(t *testing.T) {
	vertices := []float32{
		1, 2, 3,
		4, 5, 6,
		7, 8, 9}
	triangles := []uint32{
		0, 1, 2,
		1, 1, 1,
		2, 2, 2,
		2, 1, 0,
		0, 0, 0}

	mymesh := cloudmesh.IndexedMesh{Vertices: vertices, Indices: triangles}
	serMesh := mesh.SerializeIndexedMesh(mymesh)
	buf := new(bytes.Buffer)
	normalGarbage := make([]byte, 12)
	garbage := make([]byte, 2)

	for t := 0; t < 5; t++ {
		binary.Write(buf, binary.LittleEndian, normalGarbage)
		for i := 0; i < 9; i++ {
			binary.Write(buf, binary.LittleEndian, serMesh[i+(t*9)])
		}
		binary.Write(buf, binary.LittleEndian, garbage)
	}
	file := bufio.NewReader(buf)
	myMesh, err := CreateMesh(file, 5)

	if err != nil {
		t.Errorf("Error creating mesh: ", err)
	}

	for i := range myMesh.Indices {
		if myMesh.Indices[i] != triangles[i] {
			t.Errorf("Expected %v, but got: %v", triangles[i], myMesh.Indices[i])
		}
	}
	for a := range myMesh.Vertices {
		if myMesh.Vertices[a] != vertices[a] {
			t.Errorf("Expected %v, but got: %v", vertices[a], myMesh.Vertices[a])
		}
	}
}

func TestEmptyMesh(t *testing.T) {
	retVal := emptyMesh()
	if len(retVal.Vertices) != 0 {
		t.Errorf("Expected the vertices length to be 0 but was %v", len(retVal.Vertices))
	}
	if len(retVal.Indices) != 0 {
		t.Errorf("Expected the indices length to be 0 but was %v", len(retVal.Vertices))
	}
}
