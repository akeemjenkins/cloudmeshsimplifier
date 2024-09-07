package vsa

import (
	"testing"
	"bitbucket.org/cloudcomputer/cloud-mesh-simplifier/shape"
	"fmt"
)

func TestVSAVanilla(t *testing.T){
	myMesh := shape.BasicCube()
	proxies,_ := VSAVanilla(myMesh)
	for i := 0; i < len(proxies); i++ {
		fmt.Printf("Proxy: %v\n", proxies[i])
	}

	if len(proxies) != 6 {
		t.Fail()
	}
}


func TestVSAVanilla2(t *testing.T) {
	myMesh := shape.Octahedron(2000)
	proxies,_ := vsaVanillaError(myMesh,.01, 1)
	for i := 0; i < len(proxies); i++ {
		fmt.Printf("Proxy: %v\n", proxies[i])
	}

	if len(proxies) != 8 {
		t.Fail()
	}
}