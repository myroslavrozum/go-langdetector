package algos

import "math"

// Implement Cosine Distance algorythm
func CalculateColsineDistances(model map[string]float64, vector2investigate map[string]float64) float64 {
	var dotProduct, normModel, normInvestigate float64

	for k, vModel := range model {
		normModel += vModel * vModel
		if vInvestigate, ok := vector2investigate[k]; ok {
			dotProduct += vModel * vInvestigate
		}
	}

	for _, vInvestigate := range vector2investigate {
		normInvestigate += vInvestigate * vInvestigate
	}

	if normModel == 0 || normInvestigate == 0 {
		return 1.0
	}

	cosineSimilarity := dotProduct / (math.Sqrt(normModel) * math.Sqrt(normInvestigate))
	return 1.0 - cosineSimilarity
}
