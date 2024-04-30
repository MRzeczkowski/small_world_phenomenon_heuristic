package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func findBestSolution(solutions [][]float64) []float64 {
	bestSolution := solutions[0]
	bestValue := rastrigin(bestSolution)
	for _, solution := range solutions {
		currentValue := rastrigin(solution)
		if currentValue < bestValue {
			bestValue = currentValue
			bestSolution = solution
		}
	}

	return bestSolution
}

func rastrigin(x []float64) float64 {
	n := len(x)
	sum := 10.0 * float64(n)
	for _, xi := range x {
		sum += (xi*xi - 10.0*math.Cos(2*math.Pi*xi))
	}

	return sum
}

func cauchyRandom(x0 float64, gamma float64) float64 {
	return x0 + gamma*math.Tan((rand.Float64()-0.5)*math.Pi)
}

func drawInitialSolutions(numberOfSolutions int, dimensions int) [][]float64 {
	initialSolutions := make([][]float64, numberOfSolutions)
	for i := range initialSolutions {
		initialSolutions[i] = make([]float64, dimensions)
		for j := range initialSolutions[i] {
			initialSolutions[i][j] = rand.Float64()*(Max-Min) - Max
		}
	}

	return initialSolutions
}

func smallWorldPhenomenon1(
	maxIterations int,
	numberOfSolutions int,
	dimensions int,
	localSearchProbability float64,
	local mutateResult,
	distant mutateResult) []float64 {

	candidates := drawInitialSolutions(numberOfSolutions, dimensions)

	for k := 0; k < maxIterations; k++ {
		for i := range candidates {
			newLocalCandidate := local(candidates[i])
			newDistantCandidate := distant(candidates[i])

			candidatesToCheck := [][]float64{candidates[i], newLocalCandidate, newDistantCandidate}

			candidates[i] = findBestSolution(candidatesToCheck)
		}
	}

	return findBestSolution(candidates)
}

func smallWorldPhenomenon2(
	maxIterations int,
	numberOfSolutions int,
	dimensions int,
	localSearchProbability float64,
	local mutateResult,
	distant mutateResult) []float64 {

	candidates := drawInitialSolutions(numberOfSolutions, dimensions)

	for k := 0; k < maxIterations; k++ {
		for i := range candidates {
			currentCandidate := candidates[i]
			var newCandidate []float64

			if rand.Float64() <= localSearchProbability {
				newCandidate = local(currentCandidate)

			} else {
				newCandidate = distant(currentCandidate)
			}

			if rastrigin(newCandidate) < rastrigin(currentCandidate) {
				candidates[i] = newCandidate
			}
		}
	}

	return findBestSolution(candidates)
}

type mutateResult func([]float64) []float64

type MutationPackage struct {
	Name    string
	Local   mutateResult
	Distant mutateResult
}

var tooLow = 0
var tooHigh = 0

var Min = -5.12
var Max = 5.12

func clamp(value float64) float64 {
	if value < Min {
		tooLow++
		return Min
	} else if value > Max {
		tooHigh++
		return Max
	}
	return value
}

func reflect(value float64) float64 {
	if value < Min {
		tooLow++
		return 2*Min - value
	} else if value > Max {
		tooHigh++
		return 2*Max - value
	}
	return value
}

type mutateOne func(float64) float64

func isValid(value float64) bool {
	return value >= Min && value <= Max
}

func mutateWithinBounds(solution float64, mutation mutateOne) float64 {
	for {
		mutatedValue := mutation(solution)
		if isValid(mutatedValue) {
			return mutatedValue
		}
	}
}

func mutationBase(solution []float64, mutation mutateOne) []float64 {
	n := len(solution)
	mutatedSolution := make([]float64, n)

	for i := range solution {
		mutatedSolution[i] = clamp(mutation(solution[i]))
	}

	return mutatedSolution
}

func makeMutationPackage(name string, local mutateOne, distant mutateOne) MutationPackage {
	return MutationPackage{
		name,
		func(solution []float64) []float64 { return mutationBase(solution, local) },
		func(solution []float64) []float64 { return mutationBase(solution, distant) }}
}

func main() {

	dimensions := 3
	maxIterations := 1000
	numberOfSolutions := 10
	numberOfTests := 10
	localSearchProbability := 0.75

	fmt.Println("Simulation Parameters:")
	fmt.Println("- Number of dimensions:", dimensions)
	fmt.Println("- Local search Probability:", localSearchProbability)
	fmt.Println("- Max iterations per test:", maxIterations)
	fmt.Println("- Number of solutions per iteration:", numberOfSolutions)
	fmt.Println("- Number of tests per algorithm:", numberOfTests)

	fmt.Println()

	fmt.Println("| Algorithm | Average Result | Average Time (ms) |")
	fmt.Println("|-|-|-|")

	localMultiplier := 1.0
	distantMultiplier := 1.0

	//localMultiplier := 0.01
	//distantMultiplier := 0.05

	mutations := []MutationPackage{
		makeMutationPackage(
			fmt.Sprintf("Norm(X, %.1f) + Cauchy(X, %.1f)", localMultiplier, distantMultiplier),
			func(f float64) float64 { return f + rand.NormFloat64()*localMultiplier },
			func(f float64) float64 { return cauchyRandom(f, distantMultiplier) }),
		makeMutationPackage(
			fmt.Sprintf("Norm(X, %.1f) + Norm(X, %.1f)", localMultiplier, distantMultiplier),
			func(f float64) float64 { return f + rand.NormFloat64()*localMultiplier },
			func(f float64) float64 { return f + rand.NormFloat64()*distantMultiplier }),
		makeMutationPackage(
			fmt.Sprintf("Uniform(X, %.1f) + Cauchy(X, %.1f)", localMultiplier, distantMultiplier),
			func(f float64) float64 { return f + rand.Float64()*localMultiplier },
			func(f float64) float64 { return cauchyRandom(f, distantMultiplier) }),
	}

	for _, mutationPackage := range mutations {
		sum1 := 0.0
		sum2 := time.Duration(0)

		for i := 0; i < numberOfTests; i++ {
			start := time.Now()
			solution1 := smallWorldPhenomenon2(maxIterations, numberOfSolutions, dimensions, localSearchProbability, mutationPackage.Local, mutationPackage.Distant)
			result1 := rastrigin(solution1)
			sum1 += result1
			sum2 += time.Since(start)
		}

		fmt.Printf("| %s | %.4f | %v |\n", mutationPackage.Name, sum1/float64(numberOfTests), sum2.Milliseconds()/int64(numberOfTests))
	}

	fmt.Printf("Too low count =  %d\n", tooLow)
	fmt.Printf("Too low high =  %d\n", tooHigh)
}
