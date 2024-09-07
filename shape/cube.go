package shape

import (
	"bitbucket.org/cloudcomputer/cloud-mesh-simplifier/cloudmesh"
)

//BasicCube creates a cube with corners (0,0,0) and (1,1,1) with 12 triangles and outward-pointing normals
func BasicCube() cloudmesh.IndexedMesh {

	theseVertices := []float32{
		0, 0, 0,
		100, 0, 0,
		100, 0, 100,
		0, 0, 100,
		0, 100, 0,
		100, 100, 0,
		100, 100, 100,
		0, 100, 100}

	theseTriangles := []uint32{
		0, 1, 2,  //left 1
		0, 2, 3,  //left 2
		0, 5, 1,  //bottom 1
		0, 4, 5,  //bottom 2
		1, 5, 2,  //front 1
		5, 6, 2,  //front 2
		4, 6, 5,  //right 1
		4, 7, 6,  //right 2
		6, 3, 2,  //top 1
		7, 3, 6,  //top 2
		0, 3, 4,  //back 1
		4, 3, 7}  //back 2

	fullTriangles := make([]uint32, 12*3)
	copy(fullTriangles, theseTriangles)

	fullVertices := make([]float32, 8*3)
	copy(fullVertices, theseVertices)

	return cloudmesh.IndexedMesh{Indices: fullTriangles, Vertices: fullVertices}
}

////CreateCube creates a cube with approximately the input number of triangles.
//// This number must be greater than or equal to 12.
////TODO: CreateBadCube(numTriangles int) a function that creates a mesh like
////TODO: the one below, but with anomalies, e.g., non-watertight geometry.
//func CreateCube(numTriangles int){
//	numSubdivisionsPerFace := int(math.Sqrt(float64(numTriangles) / 12.0)) - 1
//	//start with the back face, then make 5 more transformed copies of it
//
//	stepSize := 1.0 / (float64(numSubdivisionsPerFace) + 1.0)
//}


