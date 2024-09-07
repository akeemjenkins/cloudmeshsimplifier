//Package mesh contains functions and data models related to a mesh
package mesh

// Mesh -- basic mesh interface
// This this is the public interface that all meshes adhere to
type Mesh interface {
	//Given an index of a triangle, returns the vertices of the triangle
	// as  a []unit32 of size 3, also returns an error if the vertices
	// cannot be returned.
	GetVertices(facet uint32) ([]uint32, error)

	// Given a vertex return []float32 representing the three points in
	// x, y, z of the point. This returns nil if there is no error.
	GetPoint(vertex uint32) ([]float32, error)

	// Return the number of triangles in the mesh.
	GetNumFacets() uint32

	// Returnt the number of vertices in the mesh.
	GetNumVertices() uint32
}
