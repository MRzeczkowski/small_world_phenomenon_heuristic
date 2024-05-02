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
			initialSolutions[i][j] = rand.Float64()*(Max-Min) + Min
		}
	}

	return initialSolutions
}

func smallWorldPhenomenon1(
	maxIterations int,
	numberOfCandidateSolutions int,
	dimensions int,
	local mutateResult,
	distant mutateResult) []float64 {

	candidates := drawInitialSolutions(numberOfCandidateSolutions, dimensions)

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
	numberOfCandidateSolutions int,
	dimensions int,
	localSearchProbability float64,
	local mutateResult,
	distant mutateResult) []float64 {

	candidates := drawInitialSolutions(numberOfCandidateSolutions, dimensions)

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

var Min = -5.12
var Max = 5.12

func clamp(value float64) float64 {
	if value < Min {
		return Min
	} else if value > Max {
		return Max
	}

	return value
}

func reflect(value float64) float64 {
	if value < Min {

		return 2*Min - value
	} else if value > Max {

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
		mutatedSolution[i] = mutation(solution[i])

		//clamp(mutation(solution[i]))
		//reflect(mutation(solution[i]))
		//mutateWithinBounds(solution[i], mutation)
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
	numberOfCandidateSolutions := 10
	numberOfTests := 100

	fmt.Println("Simulation Parameters:")
	fmt.Println("- Number of dimensions:", dimensions)
	fmt.Println("- Max iterations per test:", maxIterations)
	fmt.Println("- Number of candidate solutions:", numberOfCandidateSolutions)
	fmt.Println("- Number of tests per algorithm:", numberOfTests)
	fmt.Println()

	localMultiplierStart := 1.0
	localMultiplierEnd := 1.0
	localMultiplierStep := 0.01

	distantMultiplierStart := 0.05
	distantMultiplierEnd := 0.05
	distantMultiplierStep := 0.05

	localSearchProbabilityStart := 0.5
	localSearchProbabilityEnd := 0.5
	localSearchProbabilityStep := 0.05

	fmt.Println("| Algorithm | Local multiplier | Distant multiplier | Local Search Probability | Average Result | Average Time (ms) |")
	fmt.Println("|-|-|-|-|-|-|")

	// Initialize best results for each mutation package
	bestResults := make(map[string]float64)
	bestLocalMultipliers := make(map[string]float64)
	bestDistantMultipliers := make(map[string]float64)
	bestLocalSearchProbabilities := make(map[string]float64)

	for _, mutationName := range []string{"Normal+Cauchy", "Normal+Normal", "Uniform+Cauchy"} {
		bestResults[mutationName] = math.Inf(1)
	}

	for localSearchProbability := localSearchProbabilityStart; localSearchProbability <= localSearchProbabilityEnd; localSearchProbability += localSearchProbabilityStep {
		for localMultiplier := localMultiplierStart; localMultiplier <= localMultiplierEnd; localMultiplier += localMultiplierStep {
			for distantMultiplier := distantMultiplierStart; distantMultiplier <= distantMultiplierEnd; distantMultiplier += distantMultiplierStep {
				mutations := []MutationPackage{
					makeMutationPackage(
						"Normal+Cauchy",
						func(f float64) float64 { return f + rand.NormFloat64()*localMultiplier },
						func(f float64) float64 { return cauchyRandom(f, distantMultiplier) }),
					makeMutationPackage(
						"Normal+Normal",
						func(f float64) float64 { return f + rand.NormFloat64()*localMultiplier },
						func(f float64) float64 { return f + rand.NormFloat64()*distantMultiplier }),
					makeMutationPackage(
						"Uniform+Cauchy",
						func(f float64) float64 { return f + rand.Float64()*localMultiplier },
						func(f float64) float64 { return cauchyRandom(f, distantMultiplier) }),
				}

				for _, mutationPackage := range mutations {
					sumResults := 0.0
					sumTime := time.Duration(0)

					for i := 0; i < numberOfTests; i++ {
						start := time.Now()
						solution := smallWorldPhenomenon2(maxIterations, numberOfCandidateSolutions, dimensions, localSearchProbability, mutationPackage.Local, mutationPackage.Distant)
						result := rastrigin(solution)
						sumResults += result
						sumTime += time.Since(start)
					}

					averageResult := sumResults / float64(numberOfTests)
					fmt.Printf("| %s | %.2f | %.2f | %.2f | %.4f | %v |\n", mutationPackage.Name, localMultiplier, distantMultiplier, localSearchProbability, averageResult, sumTime.Milliseconds()/int64(numberOfTests))

					if averageResult < bestResults[mutationPackage.Name] {
						bestResults[mutationPackage.Name] = averageResult
						bestLocalMultipliers[mutationPackage.Name] = localMultiplier
						bestDistantMultipliers[mutationPackage.Name] = distantMultiplier
						bestLocalSearchProbabilities[mutationPackage.Name] = localSearchProbability
					}
				}
			}
		}
	}

	fmt.Println("\nBest Parameters Found for Each Mutation Package:")
	for name, result := range bestResults {
		fmt.Printf("- %s:\n", name)
		fmt.Printf("\t- Local Search Probability: %.2f\n", bestLocalSearchProbabilities[name])
		fmt.Printf("\t- Local Multiplier: %.2f\n", bestLocalMultipliers[name])
		fmt.Printf("\t- Distant Multiplier: %.2f\n", bestDistantMultipliers[name])
		fmt.Printf("\t- Average Result: %.4f\n", result)
	}
}
