// Package stl provides core functionality of dealing with Stereo Lithography
// files.
package stl

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"

	"bitbucket.org/cloudcomputer/cloud-mesh-simplifier/cloudmesh"
	"bitbucket.org/cloudcomputer/cloud-mesh-simplifier/mesh"
)

// Read numBytes from file into a byte slice and return the byte slice
func readbytes(numBytes int, file *bufio.Reader) (retBuf []byte, err error) {
	buf := make([]byte, numBytes)
	if _, err := io.ReadAtLeast(file, buf, numBytes); err != nil {
		log.Printf("Failed at reading %d bytes.", numBytes)
		return nil, err
	}
	return buf, nil
}

// LoadSTLFile either a binary or ASCII STL file.
// TODO: check for text STL
// TODO: read into mesh instead of triangles
func LoadSTLFile(path string) (mesh cloudmesh.IndexedMesh, err error) {

	file, err := os.Open(path)
	if err != nil {
		log.Printf("Could not open file: %s", path)
		return emptyMesh(), err
	}
	defer file.Close()

	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Printf("Could not stat: %s", path)
		return emptyMesh(), err
	}
	fileSize := int(fileInfo.Size())

	// read the header
	bufReader := bufio.NewReader(file)
	header, err := readbytes(80, bufReader)
	if err != nil {
		log.Println("Couldn't read header.")
		return emptyMesh(), err
	}

	// read the number of triangles
	numTrianglesBytes, err := readbytes(4, bufReader)
	if err != nil {
		log.Println("Couldn't read the triangle number.")
		return emptyMesh(), err
	}
	numTriangles := binary.LittleEndian.Uint32(numTrianglesBytes)

	//If the file doesn't divide out by 50 bytes after the 84 byte header and
	//the number of triangles we parsed is not perfectly divisible then the file
	//likely has an error in it.
	if ((fileSize-84)%50) != 0 || (fileSize-84)/50 != int(numTriangles) {
		log.Println("Either the file size or number of triangles are wrong.")
		return emptyMesh(), fmt.Errorf("Either the file size or number of triangles are wrong.")
	}

	log.Printf("Opened File:\t%s\nSize: %d\n", fileInfo.Name(), fileSize)
	log.Printf("Header:\t%s\nTriangles:\t%d", header, numTriangles)

	// run the decode into the mesh
	return CreateMesh(bufReader, numTriangles)
}

// CreateMesh returns a mesh given a buffer reader
//
func CreateMesh(file *bufio.Reader, numTris uint32) (m cloudmesh.IndexedMesh, err error) {
	emptyIndices := make([]uint32, 0)
	emptyVertices := make([]float32, 0)

	// a hashmap for O(n) vertex adding
	speed := make(map[[3]float32]uint32)

	// our internal indexed mesh
	newMesh := cloudmesh.IndexedMesh{Indices: emptyIndices, Vertices: emptyVertices}
	//vertexLookup = make(map[*[]float32]bool) // hashmap for O(n) lookup of vertex

	// for every triangle in the file
	for i := 0; i < int(numTris); i++ {
		triBytes, err := readbytes(50, file)
		if err != nil {
			return emptyMesh(), err
		}
		var vfloat float32

		//throw away the facet normal
		triBytes = triBytes[12:]

		// extract the 9 floats of the triangles
		newTriangle := make([]float32, 9, 9)
		for index := 0; index < 9; index++ {
			// take the first 4 bytes and convert them to a float
			vfloat = math.Float32frombits(binary.LittleEndian.Uint32(triBytes[:4]))

			triBytes = triBytes[4:] // lop off the first 4 bytes
			newTriangle[index] = vfloat
		}

		var triIndexes [3]uint32

		for i := 0; i < 3; i++ {
			// extract the key as an array of 3 floats
			key := [3]float32{newTriangle[0], newTriangle[1], newTriangle[2]}
			//index into the hashmap to see if it exists
			val, ok := speed[key]
			if ok {
				triIndexes[i] = val
			} else {
				newMesh.Vertices = append(newMesh.Vertices, newTriangle[0], newTriangle[1], newTriangle[2])
				triIndexes[i] = uint32(len(newMesh.Vertices)-3) / 3
				speed[key] = triIndexes[i]
			}
			newTriangle = newTriangle[3:] // lop off the first 3 floats
		}
		// add a triangle with indices
		newMesh.AddTriangle(triIndexes[0], triIndexes[1], triIndexes[2])

	}
	fmt.Printf("speed: %v\n", len(speed))
	fmt.Printf("Tri: %v, Vert:%v\n", newMesh.GetNumFacets(), newMesh.GetNumVertices())
	return newMesh, nil
}

//WriteSTLMesh Mesh version of the above
func WriteSTLMeshName(m mesh.Mesh,  name string){

	var header [80]byte

	numTris := m.GetNumFacets()
	fmt.Printf("Writing STL -> numTris: %d \n", numTris)
	normalGarbage := make([]byte, 12)
	garbage := make([]byte, 2)

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, header)
	binary.Write(buf, binary.LittleEndian, numTris)
	for tri := uint32(0); tri < uint32(numTris); tri++ {
		binary.Write(buf, binary.LittleEndian, normalGarbage)
		triVertices, err := m.GetVertices(tri)
		if err != nil {
			log.Fatalf("WriteMesh GetVertices: %s", err)
		}
		for p := uint32(0); p < 3; p++ {
			physicalPoint, err := m.GetPoint(triVertices[p])
			if err != nil {
				log.Fatalf("WriteMesh: %s", err)
			}
			binary.Write(buf, binary.LittleEndian, physicalPoint)
		}
		binary.Write(buf, binary.LittleEndian, garbage)

	}

	ioutil.WriteFile(name, buf.Bytes(), 0644)

}

//WriteSTLMesh Mesh version of the above
func WriteSTLMesh(m mesh.Mesh) {
	WriteSTLMeshName(m, "defaultMesh.stl")
}

func emptyMesh() cloudmesh.IndexedMesh {
	vertex := make([]float32, 0, 0)
	index := make([]uint32, 0, 0)
	return cloudmesh.IndexedMesh{Indices: index, Vertices: vertex}
}
