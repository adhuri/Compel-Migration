package strategy

func calculateValue(array []float32) float32 {
	return average(array)
}

func average(array []float32) float32 {
	var total float32
	for _, value := range array {
		total += value
	}
	return (total / float32(len(array)))

}

func max(array []float32) (maxElement float32) {
	//Find max in an array
	maxElement = array[0]
	for _, element := range array {
		if maxElement > element {
			maxElement = element
		}
	}
	return
}
