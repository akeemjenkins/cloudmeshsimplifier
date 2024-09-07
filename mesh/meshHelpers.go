package mesh

import (
	"bitbucket.org/cloudcomputer/cloud-mesh-simplifier/auxmath"
	"errors"
	"log"
	"math"
)

type MeshNeighborhood interface {
	//Returns the set of 0,1,2, or 3 neighbor triangles of the given triangle.
	//If the input triangle index is out of range on the mesh, an error is returned.
	GetTriangleNeighborsOfTriangle(tri uint32) ([]uint32, error)

	//For future use:
	//GetTriangleNeighborsOfVertex(tri uint32) ([]uint32,error)
}

//myNeighborhood a precomputed table of neighboring triangles
type myNeighborhood struct {
	triNeighbors []uint32
	m            Mesh
}

//Note: this implementation assumes the triangles are internally labelled 0,...,n
func (neighb myNeighborhood) GetTriangleNeighborsOfTriangle(tri uint32) ([]uint32, error) {
	if tri >= neighb.m.GetNumFacets() || tri < 0 {
		return []uint32{}, errors.New("GetVertices:requested index is out of bounds")
	}

	//return []uint32{neighb.triNeighbors[tri+0],
	//			neighb.triNeighbors[tri+1],
	//			neighb.triNeighbors[tri+2]}, nil
	validRetVal := make([]uint32, 0, 3)
	//TODO: there is probably  performance for these checks.  Write a test with and without them (on a large solid with no boundary)
	if neighb.triNeighbors[3*tri+0] != math.MaxUint32 {
		validRetVal = append(validRetVal, neighb.triNeighbors[3*tri+0])
	}

	if neighb.triNeighbors[3*tri+1] != math.MaxUint32 {
		validRetVal = append(validRetVal, neighb.triNeighbors[3*tri+1])
	}

	if neighb.triNeighbors[3*tri+2] != math.MaxUint32 {
		validRetVal = append(validRetVal, neighb.triNeighbors[3*tri+2])
	}

	return validRetVal, nil

}

//findIntersection Given two ORDERED sets a and b, this finds intersection
//between them in O(n) time (where n = max(|a|,|b|)).
func findIntersection(a []uint32, b []uint32) []uint32 {
	sizeA := len(a)
	sizeB := len(b)
	iB := int(0)
	iA := int(0)
	intersection := make([]uint32, 0)

	for iA < sizeA && iB < sizeB {
		currA := a[iA]
		currB := b[iB]
		if currA < currB {
			iA++
		} else if currB < currA {
			iB++
		} else { //equality: match
			intersection = append(intersection, currA)
			iA++
			iB++
		}
	}
	return intersection
}

type triNeighborsOfVertex struct {
	triNeighbors []uint32
}

//CreateNeighborhood returns a neighborhood for a given mesh
//The algorithm is O(m.GetNumVertices() + m.GetNumTriangles())
func CreateNeighborhood(m Mesh) MeshNeighborhood {

	neighbArray := make([]uint32, 3*m.GetNumFacets(), 3*m.GetNumFacets())
	//initialize to an "invalid" index.  Note that 0 won't work.
	for i := range neighbArray {
		neighbArray[i] = math.MaxUint32
	}

	triNeighborsOfVertices := make([]triNeighborsOfVertex, m.GetNumVertices(), m.GetNumVertices())
	for i := range triNeighborsOfVertices {
		//this step sets the capacity of each vector
		triNeighborsOfVertices[i].triNeighbors = make([]uint32, 0, 6) //expect the average vertex degree to be 6
	}
	for triangle := uint32(0); triangle < m.GetNumFacets(); triangle++ {
		vertices, _ := m.GetVertices(triangle)
		for i := 0; i < 3; i++ {
			vertex := vertices[i]
			triNeighborsOfVertices[vertex].triNeighbors = append(triNeighborsOfVertices[vertex].triNeighbors, triangle) //this line is the magic
		}
	}

	for triangle := uint32(0); triangle < m.GetNumFacets(); triangle++ {
		vertices2, _ := m.GetVertices(triangle)
		//fmt.Printf("vertices: %d, %d, %d \n", vertices[0], vertices[1], vertices[2])

		neighborhoodV0 := triNeighborsOfVertices[vertices2[0]].triNeighbors
		neighborhoodV1 := triNeighborsOfVertices[vertices2[1]].triNeighbors
		neighborhoodV2 := triNeighborsOfVertices[vertices2[2]].triNeighbors

		trisBar0 := findIntersection(neighborhoodV0, neighborhoodV1)
		//if(len(trisBar0) > 2) //TODO: handle non-manifold bars

		trisBar1 := findIntersection(neighborhoodV1, neighborhoodV2)
		trisBar2 := findIntersection(neighborhoodV2, neighborhoodV0)

		for neighb := 0; neighb < len(trisBar0); neighb++ {
			if trisBar0[neighb] != triangle {
				neighbArray[3*triangle+0] = trisBar0[neighb]
			}
		}

		for neighb := 0; neighb < len(trisBar1); neighb++ {
			if trisBar1[neighb] != triangle {
				neighbArray[3*triangle+1] = trisBar1[neighb]
			}
		}

		for neighb := 0; neighb < len(trisBar2); neighb++ {
			if trisBar2[neighb] != triangle {
				neighbArray[3*triangle+2] = trisBar2[neighb]
			}
		}
	}
	return myNeighborhood{triNeighbors: neighbArray, m: m}
}

// SerializeIndexedMesh takes a mesh, computes all points(Indices + Vertices)
// and return them in a slice []float32
func SerializeIndexedMesh(m Mesh) []float32 {

	returnVal := make([]float32, 0, m.GetNumFacets()*9)

	for t := uint32(0); t < m.GetNumFacets(); t++ {
		indices, err := m.GetVertices(t)
		if err != nil {
			log.Fatal(err)
		}
		for v := 0; v < 3; v++ {
			point, err := m.GetPoint(indices[v])
			if err != nil {
				log.Fatal(err)
			}
			returnVal = append(returnVal, point[0], point[1], point[2])
		}
	}
	return returnVal
}

//ComputeArea computes the area of a triangle in a Mesh
func ComputeArea(m Mesh, triIndex uint32) float32 {
	vertices, _ := m.GetVertices(triIndex)

	p1, _ := m.GetPoint(vertices[0])
	p2, _ := m.GetPoint(vertices[1])
	p3, _ := m.GetPoint(vertices[2])

	return ComputeTriangleArea(p1, p2, p3)
}

//ComputeTriangleArea computes the area of a triangle defined by 3 3D points.
//TODO: consider adding a ComputeAreaAndNormal func.  The implementation
//of each of them uses a cross product calculation, which ought to be recycled for performance reasons.
func ComputeTriangleArea(p1 []float32, p2 []float32, p3 []float32) float32 {
	//half the magnitude of the cross product
	u := []float32{p2[0] - p1[0], p2[1] - p1[1], p2[2] - p1[2]}
	v := []float32{p3[0] - p1[0], p3[1] - p1[1], p3[2] - p1[2]}
	if auxmath.Parallel(u, v, 1e-6) {
		return 0
	}
	cross := auxmath.Cross(u, v)
	return .5 * auxmath.Magnitude(cross)
}

//ComputeNormal computes the normal of a triangle in a Mesh
func ComputeNormal(m Mesh, triIndex uint32) ([]float32, error) {
	vertices, _ := m.GetVertices(triIndex)

	p1, _ := m.GetPoint(vertices[0])
	p2, _ := m.GetPoint(vertices[1])
	p3, _ := m.GetPoint(vertices[2])

	return ComputeTriangleNormal(p1, p2, p3)
}

//ComputeTriangleNormal computes the normal of a triangle defined by 3 3D points
func ComputeTriangleNormal(p1 []float32, p2 []float32, p3 []float32) ([]float32, error) {
	u := []float32{p2[0] - p1[0], p2[1] - p1[1], p2[2] - p1[2]}
	v := []float32{p3[0] - p1[0], p3[1] - p1[1], p3[2] - p1[2]}
	if auxmath.Parallel(u, v, 1e-6) {
		return []float32{0, 0, 0}, errors.New("degenerate triangle detected")

	}
	return auxmath.Normalize(auxmath.Cross(u, v)), nil
}

// Compute the geometric center of a triangle given an indexed mesh and the
// number of the triangle in the mesh indexed from 0
func ComputeCentroid(m Mesh, triangle uint32) []float32 {

	vertices, err := m.GetVertices(triangle)
	if err != nil {
		log.Fatalf("Failed getting Vertices: %s", err)
	}
	center := make([]float32, 3, 3) //assumes this is initialized as {0,0,0}

	for i := 0; i < 3; i++ {
		vertexCoords, err := m.GetPoint(vertices[i])
		if err != nil {
			log.Fatalf("Failed to get a point: %s", err)
		}
		for coord := 0; coord < 3; coord++ {
			center[coord] = center[coord] + vertexCoords[coord]/3.0
		}
	}
	return center
}
