package main

import (
	"math"
	"strings"
)

func CosineSimilarity(textA, textB string) float64 {
	text2vec := func(text string) map[string]float64 {
		vec := make(map[string]float64)
		for _, field := range strings.Fields(text) {
			count := vec[field]
			count++
			vec[field] = count
		}
		return vec
	}
	vec1 := text2vec(textA)
	vec2 := text2vec(textB)
	biggerVec := vec1
	if len(vec2) > len(vec1) {
		biggerVec = vec2
	}
	var divident float64
	for key := range biggerVec {
		divident += vec1[key] * vec2[key]
	}
	var divisorPart1 float64
	for _, value := range vec1 {
		divisorPart1 += value * value
	}
	divisorPart1 = math.Sqrt(divisorPart1)
	var divisorPart2 float64
	for _, value := range vec2 {
		divisorPart2 += value * value
	}
	divisorPart2 = math.Sqrt(divisorPart2)
	divisor := divisorPart1 * divisorPart2
	return divident / divisor
}
