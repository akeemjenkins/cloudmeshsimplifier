package vsa

import (
	"fmt"
	"log"
	"math"

	"bitbucket.org/cloudcomputer/cloud-mesh-simplifier/auxmath"
	"bitbucket.org/cloudcomputer/cloud-mesh-simplifier/mesh"
)

type plane struct {
	point  []float32
	normal []float32
}
type pError struct {
	p       *plane
	perror  float32
	trindex uint32
}

func vanillaGeometricPartition(m mesh.Mesh, proxies [] plane, pErrors [] pError) {
	for i := range proxies {
		// compute the error for all triangles in the mesh
		seedError := ComputePlaneError(m, proxies[i])
		//fmt.Printf("looking at seed: %d\n", i)
		for x := range pErrors {

			// if the error for this proxy is less than what we have on record in pErrors
			// set the proxy for this plan as well as the new error
			if seedError[x].perror < pErrors[x].perror {
				pErrors[x].p = seedError[x].p
				pErrors[x].perror = seedError[x].perror
			}
		}
	}
}

// for every plane that is consumed by x triangles, change the definition of the plane to be the average
// of the barycenters and the average of all of the normals of the triangles.
//returns the triangle index who had the worst error along with the error value
func vanillaProxyFit(m mesh.Mesh, pErrors [] pError) (uint32, float32){

	proxies := make(map[*plane]int)
	totalErr := float32(0)
	worstTri := uint32(0)
	maxError := float32(0)

	for i := range pErrors {

		tri := pErrors[i].trindex
		triError := pErrors[i].perror
		totalErr += triError
		if triError > maxError {
			worstTri = tri
			maxError = triError
		}

		triNorm, err := mesh.ComputeNormal(m, tri)
		triBarycenter := mesh.ComputeCentroid(m,tri)
		if err != nil {
			fmt.Println("Ouch that hurt.")
		}

		_, ok := proxies[pErrors[i].p]
		// if we have already seen this plane
		if ok {
			// update how many triangles have this plane as its proxy
			proxies[pErrors[i].p] += 1

			// extract the current normal(s)
			normal := pErrors[i].p.normal
			// Compute the new normal by adding the x, y and z components
			// We will renormalize at the end, so no need for division
			updatedNormal := []float32{normal[0] + triNorm[0],
						normal[1] + triNorm[1],
						normal[2] + triNorm[2]}
			// assign the new normal
			pErrors[i].p.normal = updatedNormal

			//do the same for the proxy's centroid
			//again, we will average once at the end
			center := pErrors[i].p.point
			//TODO need an add/subtract/negate service for float slices
			updatedPoint := []float32{center[0] + triBarycenter[0],
						center[1] + triBarycenter[1],
						center[2] + triBarycenter[2]}
			pErrors[i].p.point = updatedPoint

		} else {
			// we have never seen this triangle so add one
			proxies[pErrors[i].p] += 1
			// set the normal of the plane to the normal of the current triangle
			// next time we will add in the normal above
			pErrors[i].p.normal = triNorm
		}
	}
	fmt.Printf("The total error is %v\n", totalErr/float32(m.GetNumFacets()))
	fmt.Printf("The max error is %v\n", maxError)
	// iterate over the speed hashmap to compute the actual average normal for each plane
	// k = *plane, value = int( number of triangles with this plane as their proxy)
	for k, v := range proxies {
		k.normal = auxmath.Normalize(k.normal)
		//average the barycenter
		currCenter := k.point
		avgCenter := make([]float32,3)
		//TODO: auxMath.Scale
		for i := 0; i < 3; i++ {
			// compute the average of the normals
			avgCenter[i] = currCenter[i] / float32(v)
		}
		k.point = avgCenter
	}
	return worstTri,maxError
}

func initialize(numTris int) (p []pError){
	pErrors := make([]pError, int(numTris), int(numTris))
	// setup all triangles to have infinite error
	for i := 0; i < int(numTris); i++ {
		// set the error 2^32 -1 , very large
		pErr := pError{p: nil, perror: math.MaxFloat32, trindex: uint32(i)}
		pErrors[i] = pErr
	}
	return pErrors
}

func printProxies( proxies [] plane){
	numProxies := len(proxies)
	for i:= 0; i < numProxies; i++ {
		fmt.Printf("Proxy: %v\n", i)
		n := proxies[i].normal
		c := proxies[i].point
		fmt.Printf("Normal: (%v,%v,%v)\n",n[0], n[1], n[2])
		fmt.Printf("Center: (%v,%v,%v)\n\n", c[0], c[1], c[2])
	}
}

func vsaVanillaError(m mesh.Mesh, errorThreshold float32, numSeeds int) ([]plane, []pError){
	numTris := m.GetNumFacets()
	// check to make sure that we have some triangles
	if numTris < 1 {
		log.Printf("There weren't any triangles in the mesh\n")
		return nil,nil
	}

	pErrors := initialize(int(numTris))


	// The case where we have less than 10 triangles
	if numSeeds < 1 {
		numSeeds = 1
	}

	rseeds := make([]uint32, int(numTris), int(numTris))
	for i := 0; i < int(numTris); i++ {
		rseeds[i] = uint32(i)
	}

	seeds := make([]uint32, 0, numSeeds)
	// for the 10% of seeds, generate the 10% of random seed triangles
	for i := 0; i < numSeeds; i++ {
		seed := auxmath.RandomUint32(0, uint32(int(numTris)-i))
		seeds = append(seeds, rseeds[seed])
		copy(rseeds[:i], rseeds[i:])
	}

	//Convert the seed triangles to proxies
	proxies := make([]plane,numSeeds)
	for i:= 0; i < len(seeds); i++{
		normal,err := mesh.ComputeNormal(m,seeds[i])
		if err != nil {
			//continue
			log.Printf("Couldn't compute normal: %v", err)
			//return nil
		}
		center := mesh.ComputeCentroid(m,seeds[i])
		proxies[i].point = center
		proxies[i].normal = normal
	}

	numIterations := 0
	maxNumIterations := 100
	maxError := float32(math.MaxFloat32)


	for maxError > errorThreshold && numIterations < maxNumIterations{
		vanillaGeometricPartition(m, proxies, pErrors)
		worstTri, thisIterationError := vanillaProxyFit(m, pErrors)
		center := mesh.ComputeCentroid(m, worstTri)
		normal, _ := mesh.ComputeNormal(m, worstTri)

		if( thisIterationError > errorThreshold) {
			proxies = append(proxies, plane{point: center, normal: normal})
		}

		if(thisIterationError < maxError){
			maxError = thisIterationError
		}

		numIterations++
	}
	fmt.Printf("The final number of proxies is %v\n", len(proxies))

	printProxies(proxies)
	return proxies,pErrors

}

func vsaVanillaNumSeeds(m mesh.Mesh, numSeeds int) ([]plane, []pError){
	return vsaVanillaError(m, .1, numSeeds)
}

func VSAVanilla(m mesh.Mesh) ([]plane, []pError) {
	return vsaVanillaError(m,.1, 1)
}

// ComputePlaneError - given a seed triangle compute the error for every triangle in the mesh
func ComputePlaneError(m mesh.Mesh, proxy plane) (p []pError) {
	numTris := m.GetNumFacets()
	pErrors := make([]pError, 0, int(numTris))

	// for every triangle
	for i := 0; i < int(numTris); i++ {
		// compute the normal of the current triangle
		triNormal, err := mesh.ComputeNormal(m, uint32(i))
		if err != nil {
			//continue
			log.Printf("Couldn't compute normal: %v", err)
			//return nil
		}

		// Compute the difference between the normals.
		diff := make([]float32, 0, 3)
		for x := 0; x < 3; x++ {
			//fmt.Printf("Seed Normal: %d\t TriNormal: %d\n", seedPlane.normal[x], triNormal[x])
			diff = append(diff, proxy.normal[x]-triNormal[x])
		}
		mag := auxmath.Magnitude(diff) * mesh.ComputeArea(m, uint32(i))

		pErr := pError{p: &proxy, perror: mag, trindex: uint32(i)}

		// append our new error to the slice of errors to be returned
		pErrors = append(pErrors, pErr)
	}
	return pErrors
}
