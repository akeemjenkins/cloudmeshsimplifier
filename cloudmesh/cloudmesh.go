// Package cloudmesh implements a Mesh
package cloudmesh

import (
	"errors"
)

// A data structure that houses a mesh that has been indexed
type IndexedMesh struct {
	Indices  []uint32
	Vertices []float32
}

// NewMesh returns a pointer to an initialized mesh
func NewMesh() (m *IndexedMesh) {
	i := make([]uint32, 0, 0)
	v := make([]float32, 0, 0)
	mesh := IndexedMesh{Indices: i, Vertices: v}

	return &mesh
}

// Get the number of triangles in an indexed mesh
func (m IndexedMesh) GetNumFacets() uint32 {
	return uint32(len(m.Indices) / 3)
}

// Get the number of Vertices in an indexed mesh
func (m IndexedMesh) GetNumVertices() uint32 {
	return uint32(len(m.Vertices) / 3)
}

// Given an index of a triangle, return the Vertices []uint32 of that triangle
// Indexed from 0
func (m IndexedMesh) GetVertices(triangle uint32) ([]uint32, error) {

	if triangle >= m.GetNumFacets() {
		return []uint32{0, 0, 0}, errors.New("GetVertices:requested index is out of bounds")
	}
	return []uint32{m.Indices[3*triangle+0],
		m.Indices[3*triangle+1],
		m.Indices[3*triangle+2]}, nil
}

// Given a vertex return a slice of points
func (m IndexedMesh) GetPoint(vertex uint32) ([]float32, error) {
	if vertex > m.GetNumVertices() {
		return []float32{0, 0, 0}, errors.New("GetPoint:requested index is out of bounds")
	}

	return []float32{m.Vertices[3*vertex+0],
		m.Vertices[3*vertex+1],
		m.Vertices[3*vertex+2]}, nil
}

// Given an index(triangle) remove a triangle from the mesh
func (m *IndexedMesh) RemoveTriangle(triangle uint32) { //O(1), 0-based index
	//swap with last element and pop_back
	arrayIndex := 3 * triangle

	m.Indices[arrayIndex], m.Indices[len(m.Indices)-3] = m.Indices[len(m.Indices)-3], m.Indices[arrayIndex]
	m.Indices[arrayIndex+1], m.Indices[len(m.Indices)-2] = m.Indices[len(m.Indices)-2], m.Indices[arrayIndex+1]
	m.Indices[arrayIndex+2], m.Indices[len(m.Indices)-1] = m.Indices[len(m.Indices)-1], m.Indices[arrayIndex+2]

	for i := 0; i < 3; i++ {
		m.Indices = m.Indices[:len(m.Indices)-1] //pop last element of a slice
	}
}

func (m *IndexedMesh) Pop() {
	for i := 0; i < 3; i++ {
		m.Indices = m.Indices[1:]
	}
}

// Append a triangle to the mesh
func (m *IndexedMesh) AddTriangle(v1 uint32, v2 uint32, v3 uint32) { //O(1)
	m.Indices = append(m.Indices, v1, v2, v3)
}

// Return the index of the last triangle in the mesh
func (m IndexedMesh) LastTriangle() uint32 {
	return uint32((len(m.Indices) - 3) / 3)
}
