# Small-World phenomenon heuristic

This program is a simple implementation of the Small World Optimization Algorithm (SWOA) based on these two papers:
https://www.researchgate.net/publication/229039667_An_optimization_algorithm_based_on_the_small-world_phenomenon
https://www.researchgate.net/publication/221162052_Small-World_Optimization_Algorithm_for_Function_Optimization

This program is designed to optimize the Rastrigin function, a non-linear, multimodal function used commonly as a benchmark in optimization and testing of algorithms. The Rastrigin function is known for its large number of local minima, making it particularly challenging for optimization algorithms.

## Program Overview

The algorithm will generate random candidate solutions, perform a set number of iteration in which it will iterate over each candidate solution. For each of them it will randomly modify it using a local or distal mutation function and if the new solution is better it will overwrite the old candidate solution.

The main objective of this program is to experiment with different mutation strategies to find the minimum value of the Rastrigin function. The program explores combinations of parameters through a structured testing framework that systematically varies the following:
- Local and distant mutation multiplier i.e. mutation magnitude
- Probability of selecting local mutations over distant mutations

By default the program uses optimal parameters I've found in my experiments.

The algorithm can use different local + distal mutations combinations and 3 combinations of random distributions have been tested:
1. **Normal + Cauchy**: Normal distribution mutation + Cauchy distribution mutation.
2. **Normal + Norm**: Uses two normal distribution mutations with different standard deviations.
3. **Uniform + Cauchy**: Uses a uniform distribution mutation + Cauchy distribution mutation.

## Simulation results

Below is the output of the program. For brevity only the best found parameters are printed out.

Simulation Parameters:
- Number of dimensions: 3
- Max iterations per test: 1000
- Number of candidate solutions: 10
- Number of tests per algorithm: 100

| Algorithm | Local multiplier | Distant multiplier | Local Search Probability | Average Result | Average Time (ms) |
|-|-|-|-|-|-|
| Normal+Cauchy | 1.00 | 0.05 | 0.50 | 0.3327 | 3 |
| Normal+Normal | 1.00 | 0.05 | 0.50 | 0.8413 | 2 |
| Uniform+Cauchy | 1.00 | 0.05 | 0.50 | 0.5541 | 3 |

Best Parameters Found for Each Mutation Package:
- Uniform+Cauchy:
	- Local Search Probability: 0.50
	- Local Multiplier: 1.00
	- Distant Multiplier: 0.05
	- Average Result: 0.5541
- Normal+Cauchy:
	- Local Search Probability: 0.50
	- Local Multiplier: 1.00
	- Distant Multiplier: 0.05
	- Average Result: 0.3327
- Normal+Normal:
	- Local Search Probability: 0.50
	- Local Multiplier: 1.00
	- Distant Multiplier: 0.05
	- Average Result: 0.8413

Below are observations from experiments:
 - The combination of Normal and Cauchy distributions has always been the best. 
 - In most experiments 50% chance of local search was best - this is probably because of the characteristics of the Rastrigin function. 
 - Low gamma values for the Cauchy distribution have always given best results. This reduces the probability of medium range jumps but there is still the probability of long jumps.
 - Testing many combinations too a long time and adding some parallelization to tests would be beneficial to speedup the process.