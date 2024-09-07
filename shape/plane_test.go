package shape

import "testing"

func TestCreatePlane(t *testing.T) {
	singlevertices := []float32{
		0, 1, 0,
		0, 0, 0,
		1, 1, 0}
	singleindices := []uint32{0, 1, 2}
	mesh := CreatePlane(1)
	for i := range mesh.Vertices {
		if mesh.Vertices[i] != singlevertices[i] {
			t.Errorf("Vertices: expected %v and got %v\n", singlevertices[i], mesh.Vertices[i])
		}
	}
	for x := range mesh.Indices {
		if mesh.Indices[x] != singleindices[x] {
			t.Errorf("Indexes: expect %v and got %v\n", singleindices[x], mesh.Indices[x])
		}
	}

	badmesh := CreatePlane(-1)
	if len(badmesh.Vertices) != 0 && len(badmesh.Indices) != 0 {
		t.Error("CreatePlane didn't return an empty mesh for triangles less than 0 in number")
	}

}
