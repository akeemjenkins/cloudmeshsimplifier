package vsa

import (
	"testing"
	"bitbucket.org/cloudcomputer/cloud-mesh-simplifier/shape"
	"fmt"
	"bitbucket.org/cloudcomputer/cloud-mesh-simplifier/stl"
)

func TestCreateAnchorVertices(t *testing.T) {
	myMesh := shape.Octahedron(8)
	proxies, pErrors := vsaVanillaError(myMesh,.01, 1)
	for i := 0; i < len(proxies); i++ {
		fmt.Printf("Proxy: %v\n", proxies[i])
	}

	if len(proxies) != 8 {
		t.Fail()
	}

	anchorVertices,e := vsaGetAnchorVertices(pErrors,myMesh)

	if e != nil{
		t.Fail()
	}

	if len(anchorVertices) != 6{ //one for each octahedron vertex
		t.Fail()
	}

	for _,a := range anchorVertices{
		if len(a.proxyPlanes) != 4 { //the degree of each vertex on the octahedron
			t.Fail()
		}
	}
	//TODO: lastly, check that the positions referenced from each plane all match
}

func TestCreateAnchorVertices3(t *testing.T) {
	myMesh := shape.Octahedron(10)
	stl.WriteSTLMeshName(myMesh, "octahedron_10.stl")

	proxies, pErrors := vsaVanillaError(myMesh,.01, 1)
	for i := 0; i < len(proxies); i++ {
		fmt.Printf("Proxy: %v\n", proxies[i])
	}

	if len(proxies) != 8 {
		t.Fail()
	}

	anchorVertices,e := vsaGetAnchorVertices(pErrors,myMesh)

	if e != nil{
		t.Fail()
	}

	if len(anchorVertices) != 6{ //one for each octahedron vertex
		t.Fail()
	}

	for _,a := range anchorVertices{
		if len(a.proxyPlanes) != 4 { //the degree of each vertex on the octahedron
			t.Fail()
		}
	}
	//TODO: lastly, check that the positions referenced from each plane all match
}



func TestCreateAnchorVertices2(t *testing.T) {
	myMesh := shape.Octahedron(2000)
	proxies, pErrors := vsaVanillaError(myMesh,.01, 1)
	for i := 0; i < len(proxies); i++ {
		fmt.Printf("Proxy: %v\n", proxies[i])
	}

	if len(proxies) != 8 {
		t.Fail()
	}

	anchorVertices,e := vsaGetAnchorVertices(pErrors,myMesh)

	if e != nil{
		t.Fail()
	}

	if len(anchorVertices) != 6{ //one for each octahedron vertex
		t.Fail()
	}

	for _,a := range anchorVertices{
		if len(a.proxyPlanes) != 4 { //the degree of each vertex on the octahedron
			t.Fail()
		}
	}
	//TODO: lastly, check that the positions referenced from each plane all match
}


func TestCreateAnchorVertices4(t *testing.T) {
	myMesh := shape.BasicCube()
	proxies, pErrors := vsaVanillaError(myMesh,.01, 1)
	for i := 0; i < len(proxies); i++ {
		fmt.Printf("Proxy: %v\n", proxies[i])
	}

	if len(proxies) != 6 {
		t.Fail()
	}

	anchorVertices,e := vsaGetAnchorVertices(pErrors,myMesh)

	if e != nil{
		t.Fail()
	}

	for _,a := range anchorVertices {
		fmt.Printf("anchorVertex: %v\n", a)
		//for p := range a.proxyPlanes{

		//}
	}

	if len(anchorVertices) != 8{
		t.Fail()
	}

	for _,a := range anchorVertices{
		if len(a.proxyPlanes) != 3 { //the degree of each vertex on the cube
			t.Fail()
		}
	}
	//TODO: lastly, check that the positions referenced from each plane all match
}


