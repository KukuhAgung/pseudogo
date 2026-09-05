package main

import "fmt"

func SelectionSort(arr []int, N int) {
	var i, j, minIdx, temp int

	for i = 1; i <= (N - 1); i++ {
		minIdx = i
		for j = (i + 1); j <= N; j++ {
			if arr[(j)-1] < arr[(minIdx)-1] {
				minIdx = j
			}
		}
		temp = arr[(i)-1]
		arr[(i)-1] = arr[(minIdx)-1]
		arr[(minIdx)-1] = temp
	}

}

// Program: TestSelectionSort
func main() {
	data := make([]int, 6)
	var i int

	data[(1)-1] = 64
	data[(2)-1] = 25
	data[(3)-1] = 12
	data[(4)-1] = 22
	data[(5)-1] = 11
	data[(6)-1] = 90
	SelectionSort(data, 6)
	for i = 1; i <= 6; i++ {
		fmt.Print("data[", i, "] = ", data[(i)-1], "\n")
	}

}
