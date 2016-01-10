package main

import (
	"fmt"
	"math"
	"io/ioutil"
	"strings"
	"strconv"
)


type edge struct {
	length float64
	start int
	end int
}

type point struct {
	id int
	distance float64
	pred int
}

type graph struct {
	points map[int]point
	edges map[int][]edge
}


func setUp(G graph, start int) {
	G.points[start] = point{id: start, distance: 0, pred: -1}
	for index,_ := range G.points {
		if index != start {
			G.points[index] = point{id: index, distance: math.MaxFloat64, pred: -1}
		}
	}
}

func findMin(S map[int]point) int {
	localMin := math.MaxFloat64
	result := 0
	
	for _,p := range S {
		if (p.distance <= localMin)  {
			localMin = p.distance
			result = p.id
		}
	}
	return result
}


func Dijkstra(G graph, start int) {
	setUp(G, start)
	Q := make(map[int]point)
	// IMPORTANT REMARK : the data structure for Q should be an implementation of a priority queue, like a Fibonacci heap for example,
	// in order not to have to go throught all the structure when calling findMin... This will be corrected in a later commit to enhance 
	// the performance...
	for k, v := range G.points {
		Q[k] = v
	}
	var s int
	
	for len(Q) > 0 {
		s = findMin(Q)
		
		if G.points[s].distance < math.MaxFloat64 {
			delete(Q, s)
			for _,e := range G.edges[s] {
				endIndex := e.end
				if G.points[endIndex].distance > G.points[s].distance + e.length {
					G.points[endIndex] = point{id: endIndex, distance: G.points[s].distance + e.length, pred: s}
					Q[endIndex] = point{id: endIndex, distance: G.points[s].distance + e.length, pred: s}
				}
			}
		} else {
			break
		}
	} 
}

func interpretLine(line string) (int, int, float64) {
	components := strings.Split(line, "|")
	
	if len(components) == 2 {
		i, err1 := strconv.ParseInt(components[0], 10, 0)
		j, err2 := strconv.ParseInt(components[1], 10, 0)
		
		if err1 == nil && err2 == nil {
			return int(i), int(j), -20
		} else {
			fmt.Println("error reading file")
			return 1, 1, 0
		}
	} else if len(components) == 3 {
		i, err1 := strconv.ParseInt(components[0], 10, 0)
		j, err2 := strconv.ParseInt(components[1], 10, 0)
		l, err3 := strconv.ParseFloat(components[2], 64)
	
		if err1 == nil && err2 == nil && err3 == nil {
			return int(i), int(j), l
		} else {
			fmt.Println("error reading file")
			return 1, 1, 0
		}
	} else {
		return 1, 1, 0
	}
}

func getGraphFromFile (fileName string) (graph, int, int) {
	myEdges := make(map[int][]edge)
	myPoints := make(map[int]point)
	start := 1
	end := 1
	
	content, _ := ioutil.ReadFile(fileName)
	lines := strings.Split(string(content), "\n")
	
	for _, line := range lines {
		i, j, l := interpretLine(line)
		
		if l == -20 {
			start = i
			end = j
		} else {
			newEdge := edge{length: l, start: i, end: j}
			if myEdges[i] != nil {
				myEdges[i] = append(myEdges[i], newEdge )
			} else {
				newSlice := make([]edge, 1)
				newSlice[0] = newEdge
				myEdges[i] = newSlice
			}
			myPoints[i] = point{id: i, distance: 0, pred: -1}
			myPoints[j] = point{id: j, distance: 0, pred: -1}
		}
	}
	return graph{points: myPoints, edges: myEdges}, start, end
}


func getPath (G graph, start, end int) (float64, []int) {
	Dijkstra(G, start)
	
	result := make([]int, 0)
	pointToTreat := end
	
	if G.points[end].distance == math.MaxFloat64 {
		return G.points[end].distance, result
	} else {
		for pointToTreat != start {
			result = append(result, pointToTreat)
			pointToTreat = G.points[pointToTreat].pred
		}
		result = append(result, start)
	
		reversed := make([]int, len(result))
		for index, point := range result {
			reversed[len(result) - index - 1] = point
		}
		return G.points[end].distance, reversed
	}
}


func main () {
	//A FAIRE ICI : prendre les données ds le fichier source et remplir graphpaths avec
	
	G, start, end := getGraphFromFile("/tmp/graph.txt")
	
	dist, path := getPath(G, start, end)
	
	if dist == math.MaxFloat64 {
		fmt.Println("pathFrom ", start, " to ", end, ": no path")
		fmt.Println("length of the path : infinite")
	} else {
		fmt.Println("pathFrom ", start, " to ", end, ": ", path)
		fmt.Println("length of the path :", dist)
	}
}
