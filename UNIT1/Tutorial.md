Go  is a compiled lanaguage means first compile and then its binary saves and then runs 
# Gotoolchain :)  Go toolchain is the official collection of bundled software tools, compilers, and libraries used to develop, build, and manage Go applications 

go run :) compiles the source code from one or more source file  end with the .go link with the library and then run the binary which is formed during the compilation 
go build :) create the executable binary file which can be run by ./filename without the furture  processing 

package :) is a directory containing one or more .go source files that are compiled together 

println belongs to the "fmt" package

fstonic in go lang for watching the directory 

# why main  package is important 
    A complete application that can run by itself. standalone executable program 

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
