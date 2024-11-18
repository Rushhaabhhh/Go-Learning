package main 

import "fmt"

const Constant int = 10      // Exported constant (Public)
const pvtConstant int = 20   // Unexported constant (Private)

func main()  {
	fmt.Println("Hello World")

	// Variables
	var name string = "John Doe"
	var age int = 30
	var isCool bool = true

	fmt.Println(name)
	fmt.Println(age)
	fmt.Println(isCool)

	// Shorthand
	name2 := "Harry Potter"
	age2 := 17

	fmt.Println(name2)
	fmt.Println(age2)

	fmt.Println(Constant)
	fmt.Println(pvtConstant)
}