package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Printf("Part 1: %d\n", countCubes("input.txt", true))
	fmt.Printf("Part 2: %d\n", countCubes("input.txt", false))
}

type Cube struct {
	on         bool
	xMin, xMax int
	yMin, yMax int
	zMin, zMax int
}

type Reactor struct {
	cs  []Cube
	lit int
}

// countCubes returns the number of lit cubes after going through the steps
// in order. If trim is true, the cubes outside of [-50, 50] are not counted.
func countCubes(path string, trim bool) int {
	var r Reactor

	cubes := parseSteps(path, trim)
	for _, c := range cubes {
		r.addCube(c)
	}

	return r.lit
}

// Adds a cube and the negation for any overlap.
func (r *Reactor) addCube(c Cube) {
	// "Subtract" the overlap
	for _, cube := range r.cs {
		if cube.overlaps(c) {
			overlap := cube.intersection(c)
			r.lit += overlap.lit()
			r.cs = append(r.cs, overlap)
		}
	}

	// Cubes that are off are only considered for their overlaps
	if c.on {
		r.cs = append(r.cs, c)
		r.lit += c.lit()
	}
}

func (c Cube) overlaps(other Cube) bool {
	return (other.xMax >= c.xMin && other.xMin <= c.xMax) &&
		(other.yMax >= c.yMin && other.yMin <= c.yMax) &&
		(other.zMax >= c.zMin && other.zMin <= c.zMax)
}

func (c Cube) intersection(other Cube) Cube {
	if !c.overlaps(other) {
		return Cube{}
	}

	newXMin, newYMin, newZMin := max(c.xMin, other.xMin), max(c.yMin, other.yMin), max(c.zMin, other.zMin)
	newXMax, newYMax, newZMax := min(c.xMax, other.xMax), min(c.yMax, other.yMax), min(c.zMax, other.zMax)

	// If two on cubes overlap, the common intersection needs to be subtracted
	// If an on and off cube overlap, the intersection needs to be subtracted
	// If two off cubes overlap, the common intersection needs to be added,
	// crucially the only source of off cubes are overlaps generated by the previous two cases.
	return Cube{
		on:   !c.on,
		xMin: newXMin,
		xMax: newXMax,
		yMin: newYMin,
		yMax: newYMax,
		zMin: newZMin,
		zMax: newZMax,
	}
}

func (c Cube) lit() int {
	volume := (c.xMax - c.xMin + 1) * (c.yMax - c.yMin + 1) * (c.zMax - c.zMin + 1)
	if c.on {
		return volume
	}
	return -volume
}

func parseSteps(path string, trim bool) []Cube {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	steps := []Cube{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var (
			on         string
			xMin, xMax int
			yMin, yMax int
			zMin, zMax int
		)

		_, err := fmt.Sscanf(scanner.Text(), "%s x=%d..%d,y=%d..%d,z=%d..%d", &on, &xMin, &xMax, &yMin, &yMax, &zMin, &zMax)
		if err != nil {
			panic(err)
		}

		if trim {
			// Cutout the cubes outside of [-50, 50]
			if xMin < -50 || yMin < -50 || zMin < -50 ||
				xMax > 50 || yMax > 50 || zMax > 50 {
				continue
			}
		}

		steps = append(steps, Cube{
			on:   on == "on",
			xMin: xMin,
			xMax: xMax,
			yMin: yMin,
			yMax: yMax,
			zMin: zMin,
			zMax: zMax,
		})
	}

	return steps
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
