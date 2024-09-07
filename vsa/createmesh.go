package vsa

import (
	"bitbucket.org/cloudcomputer/cloud-mesh-simplifier/mesh"
	"fmt"
	"bitbucket.org/cloudcomputer/cloud-mesh-simplifier/auxmath"
	"math"
	"errors"
)

//TODO: add "edge extraction" (from the paper), or determine if it's unnecessary

type proxyVertex struct{
	proxies [] plane
	meshIndex uint32
}

type vsaPolygon struct{
	proxy* plane
	vertices []proxyVertex
}

//TODO: consider passing the plane by ref (pointer)
func projectPointOntoPlane(point []float32, pl plane) []float32 {
	//assumes the normal is a unit vector

	v,_ := auxmath.Subtract(point, pl.point)
	pointDist,_ := auxmath.Dot(v,pl.normal)
	pointDist = float32(math.Abs(float64(pointDist)))
	nScale := auxmath.Scale(pl.normal, -pointDist)
	retVal,_:= auxmath.Subtract(point,nScale)
	return retVal
}

/*Returns the position of a "proxy vertex," which is a vertex of the mesh
 that belonging to 3 or more proxies (or technically 2 or more in regions with a mesh boundary).
 */
func proxyVertexPosition(m mesh.Mesh, p proxyVertex) (retVal []float32, err error){
	//Project this mesh vertex onto each of its proxy planes and take the average
	meshPos,err := m.GetPoint(p.meshIndex)
	if err != nil{
		return make([]float32, 0, 3), fmt.Errorf("proxyVertexPosition: bad input mesh index")
	}

	retVal = make([]float32, 0, 3)
	for i:= range p.proxies{
		proxy := p.proxies[i]
		proj := projectPointOntoPlane(meshPos,proxy)
		retVal,_ = auxmath.Add(retVal, proj)
	}

	retVal = auxmath.Scale(retVal,float32(len(p.proxies)) )

	return retVal,nil
}

func segmentOneFace() {
	//for
}

type anchorVertex struct{
	index uint32 //index in the mesh
	proxyPlanes map[*plane]bool
}

/*Computes the anchor vertices via a coloring algorithm.*/
func vsaGetAnchorVertices(pErrors [] pError, m mesh.Mesh) ([]anchorVertex, error) {

	coloredTris := make([]bool,m.GetNumFacets())

	neighborhood := mesh.CreateNeighborhood(m)

	//borderTris := make([]pError,0)
	borderTris := make(map[pError]bool)
	allColored := false
	debugNumProxyGroups := 0
	for !allColored{
		//seed the next proxyGroup with the first uncolored tri
		allColored = true;
		proxyGroup := make([]pError,0)
		for i := range coloredTris{
			if !coloredTris[i]{
				coloredTris[i] = true
				proxyGroup = append(proxyGroup, pErrors[i])
				allColored = false
				break;
			}
		}

		debugNumProxyGroups++
		//next proxy group
		for len(proxyGroup) > 0{

			//pop
			currTri := proxyGroup[0]
			proxyGroup = proxyGroup[1:]

			neighbors,_ := neighborhood.GetTriangleNeighborsOfTriangle(currTri.trindex)
			for n := range neighbors{
				neighb := neighbors[n]
				//get the proxy associated with this triangle
				//can we assume that the mesh is indexed in such a way that "neighb"
				// can used as an index into "pErrors"?  For now it seems that we can assume this since
				// vsa.initialize has this assumption (which it probably shouldn't; a one-time scan of the
				//mesh needs to be done for correctness)
				if neighb > uint32(len(pErrors)){
					return make([]anchorVertex,0),errors.New("mesh needs to be reindexed before it can be used in vsa")
				}
				neighbProxy := pErrors[neighb]
				if neighbProxy.trindex != neighb{
					return make([]anchorVertex,0),errors.New("mesh needs to be reindexed before it can be used in vsa")
				}

				if neighbProxy.p != currTri.p {
					//borderTris = append(borderTris, currTri)
					borderTris[currTri] = true
					continue
					//don't worry about neighbProxy; this will be added later
				}else{ //triangles are part of the same proxy group
					if !coloredTris[neighb]{
						coloredTris[neighb] = true;
						proxyGroup = append(proxyGroup,neighbProxy)
					}

				}
			}
		}
	}

	//Take borderTris and create a hashMap (vertex,NumProxies).  Might need to recompute a separate borderTris
	//for each proxy group

	borderVertexDegrees := make(map[uint32]anchorVertex)
	for b,_ := range borderTris{

		bIndex := b.trindex;
		vertices,err := m.GetVertices(bIndex)
		if err != nil{
			return make([]anchorVertex,0), errors.New("Could not get the vertices of the given border triangle")
		}
		for j := range vertices{
			vertexIndex := vertices[j]
			v, exists := borderVertexDegrees[vertexIndex]
			if exists{
				v.proxyPlanes[b.p] = true
				borderVertexDegrees[vertexIndex] = v;
			}else{
				a := anchorVertex{index: bIndex, proxyPlanes:make(map[*plane]bool)}
				a.proxyPlanes[b.p] = true
				borderVertexDegrees[vertexIndex] = a
			}
		}
	}
	//prune borderVertexDegrees to only include those with degree 3 or greater
	anchorVertices := make([]anchorVertex,0)
	for _,a := range borderVertexDegrees{
		valid := len(a.proxyPlanes) >= 3
		if !valid {
			continue
		}
		anchorVertices = append(anchorVertices, a)
	}
	//TODO: unlike the original VSA paper, we probably also want to insert the vertices that are on the mesh border
	//(for non-solid meshes we want anchor vertices on the mesh boundary as well)

	//anchor vertices are defined as those that map to 3 or greater
	return anchorVertices,nil
}
