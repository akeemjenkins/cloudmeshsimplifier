package shape

import (
	"testing"
	"bitbucket.org/cloudcomputer/cloud-mesh-simplifier/stl"
)


func TestBasicCube(t *testing.T) {
	cube := BasicCube()
	if( 12 != cube.GetNumFacets()){
		t.Error("Expected 12 triangles.")
	}
	if(8 != cube.GetNumVertices()) {
		t.Error("Expected 8 vertices")
	}
	stl.WriteSTLMeshName(cube, "cubeTest.stl")
}


