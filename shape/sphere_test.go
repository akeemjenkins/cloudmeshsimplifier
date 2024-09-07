package shape

import "testing"

// This is a light test that just checks the number of triangles
// and vetices because we hand coded the mesh.
func TestCreateOctahedronMesh(t *testing.T) {
	retVal := createOctahedronMesh(1)
	if retVal.GetNumVertices() != 6 {
		t.Errorf("Expected 6 vertices and got %v\n", len(retVal.Vertices))
	}
	if retVal.GetNumFacets() != 8 {
		t.Errorf("Expected 8 facets and got %v\n", len(retVal.Indices))
	}
}

func TestFacetSphere(t *testing.T) {
	// light test for only 1 triangle
	retVal := FacetSphere(1)
	if retVal.GetNumVertices() != 6 {
		t.Errorf("Expected 6 vertices and got %v\n", retVal.GetNumVertices())
	}
	if retVal.GetNumFacets() != 8 {
		t.Errorf("Expected 8 facets and got %v\n", retVal.GetNumFacets())
	}
	// This is also a simple check for the number of triangles and vertices.
	retVal = FacetSphere(10)
	if retVal.GetNumVertices() != 7 {
		t.Errorf("Expected 21 vertices and got %v\n", retVal.GetNumVertices())
	}
	if retVal.GetNumFacets() != 10 {
		t.Errorf("Expected 8 facets and got %v\n", retVal.GetNumFacets())
	}
}
