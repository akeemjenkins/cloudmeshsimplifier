package main

import (
	//"fmt"
	"bitbucket.org/cloudcomputer/cloud-mesh-simplifier/shape"
	//"bitbucket.org/cloudcomputer/cloud-mesh-simplifier/mesh"
	"bitbucket.org/cloudcomputer/cloud-mesh-simplifier/vsa"
	"bitbucket.org/cloudcomputer/cloud-mesh-simplifier/stl"
)

func main() {

	//it TriangleIterator = FacetSphere(r double )

	octahedron := shape.Octahedron(2000)
	stl.WriteSTLMeshName(octahedron, "fancyOctahedron.stl")
	_,_ = vsa.VSAVanilla(octahedron)

	/*go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	*/
	//myMesh, _ := stl.LoadSTLFile("/Users/bluzytrix/Downloads/Charon 01.stl")
	//myMesh := shape.CreatePlane(600)
	//myMesh := shape.BasicCube()
	//e := vsa.VSAVanilla(myMesh)
	//for i := 0; i < len(e); i++ {
	//	fmt.Printf("Proxy: %v\n", e[i])
	//}

	//stl.WriteSTLMesh(myMesh)
}
