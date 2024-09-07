// Package sphere has functions to generate a sphere in a mesh
package shape

import (
	"log"

	"bitbucket.org/cloudcomputer/cloud-mesh-simplifier/auxmath"
	"bitbucket.org/cloudcomputer/cloud-mesh-simplifier/cloudmesh"
	"bitbucket.org/cloudcomputer/cloud-mesh-simplifier/mesh"
)

// Create an Octahedron mesh
func createOctahedronMesh(capacityInTriangles uint32) cloudmesh.IndexedMesh {
	theseVertices := []float32{
		0, 0, 100,
		100, 0, 0,
		0, 100, 0,
		-100, 0, 0,
		0, -100, 0,
		0, 0, -100}

	theseTriangles := []uint32{
		1, 2, 0,
		2, 3, 0,
		3, 4, 0,
		4, 1, 0,
		5, 1, 4,
		5, 2, 1,
		5, 3, 2,
		5, 4, 3}
	fullTriangles := make([]uint32, 8*3, 3*capacityInTriangles+8*3)
	copy(fullTriangles, theseTriangles)

	fullVertices := make([]float32, 6*3, 3*capacityInTriangles+6*3)
	copy(fullVertices, theseVertices)

	return cloudmesh.IndexedMesh{Indices: fullTriangles, Vertices: fullVertices}
}

//returns a list of triangles from a sphere centered at the origin with radius r
func FacetSphere(numTriangles uint32) cloudmesh.IndexedMesh {
	//Algorithm: Initialization: start with an octahedron with points lying on the unit sphere
	//loop: for each triangle, get the point at its barycenter and add it to the mesh,
	//except normalize it.  remove the older triangle from the mesh and add three new ones
	//that connect to the new vertex
	//use a stack to order triangles from big to small
	//terminate when #triangles is within 3 of the target #.
	capacity := numTriangles
	baseSphere := createOctahedronMesh(capacity) //a simple faceted "unit sphere"
	currNumTris := baseSphere.GetNumFacets()

	//debugI := 0

	for currNumTris < numTriangles {
		// if debugI == 0 || debugI == 1 || debugI == 2 || debugI == 3 || debugI == 4 {
		// 	baseSphere.PrintTriangleIndices()
		// }

		tri := uint32(0) //baseSphere.LastTriangle()

		vertices, err := baseSphere.GetVertices(tri)
		if err != nil {
			log.Fatalf("Error getting Vertices: %s", err)

		}

		center := mesh.ComputeCentroid(baseSphere, tri)
		centerNorm := auxmath.Normalize(center)
		baseSphere.Vertices = append(baseSphere.Vertices, centerNorm[0], centerNorm[1], centerNorm[2])
		newVertex := uint32(len(baseSphere.Vertices)-3) / 3

		//DEBUG
		//if debugI == 0 || debugI == 1 || debugI == 2 || debugI == 3 || debugI == 4 {
		//fmt.Printf("New point: %g, %g, %g \n\n", centerNorm[0], centerNorm[1], centerNorm[2])
		//}
		//END DEBUG

		baseSphere.Pop()
		baseSphere.AddTriangle(newVertex, vertices[2], vertices[1])
		baseSphere.AddTriangle(newVertex, vertices[1], vertices[0])
		baseSphere.AddTriangle(newVertex, vertices[0], vertices[2])
		currNumTris = baseSphere.GetNumFacets()

		//debugI++
	}
	return baseSphere //serializecloudmesh.IndexedMesh(baseSphere)
}
