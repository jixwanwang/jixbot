package nlp

import (
	"math"
	"strings"

	"github.com/agonopol/go-stem"
)

func cleanWords(s string) []string {
	words := []string{}
	for _, w := range strings.Split(s, " ") {
		if len(w) == 0 {
			continue
		}
		word := stopwordRegex.ReplaceAll([]byte(w), []byte(""))
		if len(word) == 0 {
			continue
		}
		a := stemmer.Stem(word)
		words = append(words, string(a))
	}
	return words
}

func countWords(words []string) map[string]int {
	counts := map[string]int{}
	for _, word := range words {
		if _, ok := counts[word]; !ok {
			counts[word] = 0
		}
		counts[word]++
	}

	return counts
}

func cosineSimilarity(count1, count2 map[string]int) float64 {
	// Calculate cosine score using the formula:
	// cos(v1, v2) = ( v1 . v2 )/ ( ||v1|| * ||v2|| )

	// Calculate product of lengths
	sum1 := sumSquares(count1)
	sum2 := sumSquares(count2)
	denominator := math.Sqrt(float64(sum1 * sum2))

	// Dot product of vectors
	sharedWords := sharedKeys(count1, count2)
	numerator := 0
	for _, word := range sharedWords {
		numerator += count1[word] * count2[word]
	}

	if denominator == 0 {
		return 0
	}

	return float64(numerator) / denominator
}

func sumSquares(m map[string]int) int {
	sum := 0
	for _, v := range m {
		sum += v * v
	}
	return sum
}

func sharedKeys(m1, m2 map[string]int) []string {
	shared := []string{}
	for k := range m1 {
		if _, ok := m2[k]; ok {
			shared = append(shared, k)
		}
	}
	return shared
}
