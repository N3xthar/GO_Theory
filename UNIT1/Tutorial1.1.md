Go  is a compiled lanaguage means first compile and then its binary saves and then runs 
# Gotoolchain :)  Go toolchain is the official collection of bundled software tools, compilers, and libraries used to develop, build, and manage Go applications 

go run :) compiles the source code from one or more source file  end with the .go link with the library and then run the binary which is formed during the compilation 
go build :) create the executable binary file which can be run by ./filename without the furture  processing 

package in go is a collection of a go source files that are cmpilled together and provide the relevant functionality 
    all the go files are in the same package declare the same package naem 
    package help to organise code and prompt modularity and enabel the code resuse it 

println belongs to the "fmt" package

fstonic in go lang for watching the directory 

# why main  package is important 
    A complete application that can run by itself. standalone executable program 

    Answer :) package main is a special package in Go used to create executable programs. The Go compiler treats it as the application's entry package. When combined with func main(), it defines where program execution begins. All other packages are libraries that provide reusable functionality and cannot be executed directly 


use fmt inside so how it become standalone 
    answer package main is the application. Imported packages like fmt are just helpers. The application is still standalone because it starts and runs by itself 

# why main function is important 
    it is the starting point of the execution of the program without it go  does not know where to start executation 

# gofmt :) it is a tool that automatically formats your Go code.

gofmt -w main.go 

The go tool's fmt subcommand applies gofmt to all the files... 
    go fmt  :) find all the go files recursively and format it 



# 1.2 Command line argument 
    some programs generate their own data but in many scenario they get the data from files output of the other fiels and also from the other networks or form the user  via terminal


# package 
    os :) this package help you  directly talk to the operating system 
    eg :) read delete and environment vairable etc 

The command lines variable are availables in the os.Args package help to read the extra information to it 
        eg :) go run main.go  Aman 20 360 
        then go creates 
        os.Args = [
    "main.go",
    "Aman",
    "22",
    "India",
]
    and then this happes 
    | Index | Value   |
| ----- | ------- |
| 0     | main.go |
| 1     | Aman    |
| 2     | 22      |
| 3     | India   |


The  os.Args[0] is the name of the program or command that was executed  

# Variable 
Var :) If a variable is declared using var but is not explicitly initialized, Go automatically initializes it with its zero value.

:= always requires a value, so Go does not need to assign a zero value. and it is the most compact and use within the methods , not a package level variables

# Zero Value of go data type 
| Data Type            | Zero Value          | Example                        |
| -------------------- | ------------------- | ------------------------------ |
| `int`                | `0`                 | `var age int` → `0`            |
| `float32`, `float64` | `0.0`               | `var price float64` → `0.0`    |
| `bool`               | `false`             | `var isAdmin bool` → `false`   |
| `string`             | `""` (empty string) | `var name string` → `""`       |
| Pointer              | `nil`               | `var p *int` → `nil`           |
| Slice                | `nil`               | `var nums []int` → `nil`       |
| Map                  | `nil`               | `var m map[string]int` → `nil` |
| Channel              | `nil`               | `var ch chan int` → `nil`      |
| Function             | `nil`               | `var f func()` → `nil`         |
| Interface            | `nil`               | `var i interface{}` → `nil`    |


# LOOPS 
     way 1 
    for initlization ; condition ; increment {
    }

    way 2 while loop 

    for condition {
    }
    
    way 3 infinite loop 
    
    for {
    } and can be terminated using the return and break keywords 

    
    way 4 range loop 
    for index , value := range iterating elements{
        work done 
    }
    
    we have to handle the both index and value and use it go dones not allow to left the variable unused to we have to use teh 
        blank identifier _ (underscore)  
