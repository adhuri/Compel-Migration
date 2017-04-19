package strategy

func calculateValue(array []float64) float64 {
	return average(array)
}

func average(array []float64) float64 {
	var total float64
	for _, value := range array {
		total += value
	}
	return (total / float64(len(array)))

}

func max(array []float64) (maxElement float64) {
	//Find max in an array
	maxElement = array[0]
	for _, element := range array {
		if maxElement > element {
			maxElement = element
		}
	}
	return
}
