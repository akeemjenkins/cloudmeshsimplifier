package shape

import (
	"bitbucket.org/cloudcomputer/cloud-mesh-simplifier/mesh"
	"bitbucket.org/cloudcomputer/cloud-mesh-simplifier/cloudmesh"
	"log"
)
//TODO: refactor this algorithm used for sphere generation
func Octahedron(numTriangles uint32) cloudmesh.IndexedMesh {
	//Algorithm: Initialization: start with an octahedron with points lying on the unit sphere
	//loop: for each triangle, get the point at its barycenter and add it to the mesh,
	//except normalize it.  remove the older triangle from the mesh and add three new ones
	//that connect to the new vertex
	//use a queue to order triangles from big to small
	//terminate when #triangles is within 3 of the target #.
	capacity := numTriangles
	baseSphere := createOctahedronMesh(capacity) //a simple faceted "unit sphere"
	currNumTris := baseSphere.GetNumFacets()

	//debugI := 0

	for currNumTris < numTriangles {


		tri := uint32(0) //baseSphere.LastTriangle()

		vertices, err := baseSphere.GetVertices(tri)
		if err != nil {
			log.Fatalf("Error getting Vertices: %s", err)

		}

		center := mesh.ComputeCentroid(baseSphere, tri)
		//centerNorm := auxmath.Normalize(center)
		//baseSphere.Vertices = append(baseSphere.Vertices, centerNorm[0], centerNorm[1], centerNorm[2])
		baseSphere.Vertices = append(baseSphere.Vertices, center[0], center[1], center[2])
		newVertex := uint32(len(baseSphere.Vertices)-3) / 3

		baseSphere.Pop()
		baseSphere.AddTriangle(vertices[1], newVertex, vertices[0])
		baseSphere.AddTriangle(vertices[0], newVertex, vertices[2])
		baseSphere.AddTriangle(vertices[2], newVertex, vertices[1])

		currNumTris = baseSphere.GetNumFacets()

		//debugI++
	}
	return baseSphere //serializecloudmesh.IndexedMesh(baseSphere)
}

