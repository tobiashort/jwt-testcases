package main

import "testing"

func TestSimilarity(t *testing.T) {
	type Example struct {
		textA         string
		textB         string
		minSimilarity float64
	}
	examples := []Example{
		{
			textA:         "machine learning",
			textB:         "machine learning",
			minSimilarity: 0.99,
		},
		{
			textA:         "machine learning",
			textB:         "machine learning",
			minSimilarity: 0.99,
		},
		{
			textA:         "machine learning",
			textB:         "machine and learning",
			minSimilarity: 0.81,
		},
	}
	for _, example := range examples {
		similarity := CosineSimilarity(example.textA, example.textB)
		if similarity < example.minSimilarity {
			t.Error(similarity)
		}
	}
}
