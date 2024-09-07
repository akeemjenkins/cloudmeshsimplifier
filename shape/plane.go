package shape

import (
	"fmt"
	"math"

	"bitbucket.org/cloudcomputer/cloud-mesh-simplifier/cloudmesh"
)

// CreatePlane creates a plane given numTris number of triangles
// Requesting a Plane of 0 or less triangles will return a blank mesh.
// This plane has zero in the z direction with a y of either 0 or 1
// Starting with a cartesian plane, we extend out in the positive x direction
// forevever to create the number of triangles we need.
// _________
// | / | / |
// |/__|/__| etc
//
func CreatePlane(numTris int) *cloudmesh.IndexedMesh {

	if numTris <= 0 {
		fmt.Println("You can't have a plane with no triangles.")
		return cloudmesh.NewMesh()
	}

	// For a given number of triangles, it takes + 2 vertices to make that triangle
	numVertexes := (numTris + 2)
	theseVertices := make([]float32, 0, numVertexes)
	for i := 0; i < numVertexes; i++ {
		// The x coordinate goes 0,0,2,2,4,4,6,6 etc
		// We take the floor of the current index divided by 2 to get the x coordinates
		// that increment by 2, every 2
		x := float32(math.Floor(float64(i / 2)))
		// The y coordinate oscilates 0,1,0,1,0,1 etc
		// mod 2 provides us the perfect value for the y coordinate
		y := float32(i % 2)
		// The z coordinate is always 0.
		z := float32(0)
		theseVertices = append(theseVertices, x, y, z)
	}

	// The initial vertices are ordered starting at x,y,z = 0,0,0 from bottom to top, left to right
	// The ordering is 0(0,0,0), 1(0,1,0), 2(1,0,0)
	// This ordering doesn't work with an odd number of triangles so we reverse
	// two vertices at a time so that the ordering is 0(0,1,0), 1 (0,0,0), 2 (1,1,0)
	// This makes it possible so that the math for creating the triangles will
	// accomodate an odd number of triangles.
	for i := 0; i < len(theseVertices); i += 6 {

		// this is the case that we have an odd number of triangles like the 3 below
		// ________
		// | / | /
		// |/__|/
		//
		if len(theseVertices)-i <= 3 {
			// add one to the y component and stop to "swap" the index of the vertices
			// in the triangle above vertex 4 is 2,0,0 flipped around is 2,1,0
			theseVertices[i+1] = theseVertices[i+1] + 1
			break // we have flipped the last triangle around so exit
		}
		theseVertices[i], theseVertices[i+3] = theseVertices[i+3], theseVertices[i]
		theseVertices[i+1], theseVertices[i+1+3] = theseVertices[i+1+3], theseVertices[i+1]
		theseVertices[i+2], theseVertices[i+2+3] = theseVertices[i+2+3], theseVertices[i+2]
	}

	// The triangles here are very easy to create because of the flipping above.
	// They are of the form:
	// [0,1,2], [3,2,1], [2,3,4], [5,4,3] etc
	theseTriangles := make([]uint32, 0, numTris)
	for i := 0; i < numTris; i++ {
		var x uint32
		var y uint32
		var z uint32
		if i%2 == 0 {
			x = uint32(i)
			y = uint32(i + 1)
			z = uint32(i + 2)
		} else {
			x = uint32(i + 2)
			y = uint32(i + 1)
			z = uint32(i)
		}
		theseTriangles = append(theseTriangles, x, y, z)

	}
	// setup the new mesh and swap in the vertices and indices
	ret := cloudmesh.NewMesh()
	ret.Vertices = theseVertices
	ret.Indices = theseTriangles
	return ret
}
