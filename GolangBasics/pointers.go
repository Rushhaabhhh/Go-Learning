// Pointers

package main 

func pointer() {

	increment := func(num int) {
		num++
		println("Inside : ", num, &num)
	}
	count := 42

	// Increment declares count as  pointer variable whose value is always an address and points to an integer value
	incrementAddr := func(num *int) {
		*num++
		println("Inside Addr: ", num, &num)
	}

	// Pass by Value 

	// Displays value of count and its memory address
	println("Before : ", count, &count)
	// Pass the value of count to the function
	increment(count)
	println("After : ", count, &count)

	incrementAddr(&count) // Pass the address of count to the function
	println("After Addr: ", count, &count)

	// Pass by Reference

}
