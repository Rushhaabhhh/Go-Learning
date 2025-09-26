package main 

import "fmt"

const Constant int = 10      // Exported constant (Public)
const pvtConstant int = 20   // Unexported constant (Private)

func main()  {
	fmt.Println("Hello World")

	// Variables

	// Declare variables that are set to their zero values
	var a int 
	var b string // zero value is an empty string
	var c float64 
	var d bool
	fmt.Println(a, b, c, d) // 0
	
	// Declare and assign values
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

	// Convesrion 
	var x int = 100
	var y float64 = float64(x) // Convert int to float64
	fmt.Println(y)

	// Structs
	type example struct {
		pi float32
		radius int16
		length int16
		breadth int16
		isValid bool
	}
	var ex example 
	fmt.Println(ex) // Print zero value of struct

	// Assign values to struct fields 
	ex2 := example{
		pi: 3.14,
		radius: 5,
		length: 10,
		breadth: 15,
		isValid: true,
	}
	fmt.Println(ex2)

	// Anonymous Struct : Struct without a name using a struct literal, useful for one-time use
	ex3 := struct {
		name string
		age int
	}{
		name: "Alice",
		age: 25,
	}
	fmt.Println(ex3)

	type Alice struct {
		name string
		age int
	}
	type Bob struct {
		name string
		age int
	}
	var person1 Alice
	var person2 Bob

	//person1 = person2 // Error: cannot use person2 (type Bob) as type Alice in assignment due to integerity
	person1 = Alice(person2)
	fmt.Println(person1)
}