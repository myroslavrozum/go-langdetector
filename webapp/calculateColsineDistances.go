package webapp

import "math"

// Cosine Distance algorythm
// https://en.wikipedia.org/wiki/Cosine_similarity
func calculateColsineDistances(model map[string]float64, vector2compare map[string]float64) float64 {
	var dotProduct, normModel, norm2comapare float64

	for k, frequencyModel := range model {
		normModel += frequencyModel * frequencyModel
		if frequency2compare, ok := vector2compare[k]; ok {
			dotProduct += frequencyModel * frequency2compare
		}
	}

	for _, frequency2compare := range vector2compare {
		norm2comapare += frequency2compare * frequency2compare
	}

	if normModel == 0 || norm2comapare == 0 {
		return 1.0
	}

	cosineSimilarity := dotProduct / (math.Sqrt(normModel) * math.Sqrt(norm2comapare))
	return 1.0 - cosineSimilarity
}
