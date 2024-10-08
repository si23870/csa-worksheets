package main

import "fmt"

func addOne(a int) int {
	return a + 1
}

func square(a int) int {
	return a * a
}

func double(slice []int) {
	slice2 := slice
	slice = append(slice2, slice2...)
}

func mapSlice(f func(a int) int, slice []int) {
	for i := range slice {
		slice[i] = f(slice[i])
	}
}

func mapArray(f func(a int) int, array [5]int) [5]int {
	for i := range array {
		array[i] = f(array[i])
	}
	return array
}

func main() {

	intSlice := []int{1, 2, 3, 4, 5}
	mapSlice(addOne, intSlice)
	fmt.Println(intSlice)
	intArray := [5]int{1, 2, 3, 4, 5}
	intArray = mapArray(addOne, intArray)
	fmt.Println(intArray)

	double(intSlice)
	fmt.Println(intSlice)
}
