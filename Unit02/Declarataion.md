# 1  Declaration 
    introducing the entity and defining some of the properties such as name , type , scope , return type 

eg :) 
  var an  int =  20 

# 2  Four major declaration 
  var const type and func 
| Declaration | Purpose            | Example              |
| ----------- | ------------------ | -------------------- |
| `var`       | Declares variables | `var age int`        |
| `const`     | Declares constants | `const PI = 3.14`    |
| `type`      | Declares new types | `type User struct{}` |
| `func`      | Declares functions | `func main(){}`      |

# 3 Structure of the go file 

package declaration

import declaration

package-level declaration 

executable codes 

# 4 package declaration

# Panic 
  A panic is a state where the current goroutine has encountered an unrecoverable condition, so Go stops normal execution, begins unwinding the stack, runs deferred functions, and unless the panic is recovered, the program terminates.
  
every go source files exactly belongs to the exact one package 

it organise the code

# Short hand  Variable 
            -> always wihtin the functions  that is the local variables 
            -> go automatically detect its type or infer the type ;
          -> (syntax)    name := expresssion 
        declare and intilize the variable 
    Differnce 
  := 
  Declares and intilize
  only inside the function 
  automatically type inference
  
var 
  u want an explicit type 
  var applie int 

  when u dont have the value 

# 5 Pointers 
  a variables which is used to store the addresss 
  
  &x :) get the address of x 
  
  A pointer stores the memory address of a variable.
  It tells where a value is stored, not the value itself.
  Every variable has an address, but not every value has an address
  
  
#  Address Operators (&)
  ->    return the address of a variables
  
  example :) 
    var x int 
    p := &x 
  
  x -> variable
  &x -> featch the memeory address 
  p -> pointer to the memeory address

# Dereference Operators

  *p = go and fetch the value which is stored in the p 

  *p means the value stored at the address contained in p.
  Since *p represents the actual variable, it can be both read and modified. 
  
  go and fetch the value which is stroed in the p variable 

# Addressable Values 
  only variables have addresses 

# nil Pointers

  -> a nil pointer is the zero value of any pointer type in go , it does nto point to any valid memory address , before Dereferencing  a pointer u should check whether it is nill or nor 

# pointer comparison 
  two pointer can be equal only if 

->   they pointing the same variable
->   both are nil 

# Returning pointer to the funciton 
  the function return the address of a local valriable like this 
  func fun()*int{
    heelo:="My name is amandeep"
  return &heelo;
  }

# Passing Pointers to Functions
 passing a pointer allows a function to modify the original variables
  instead of passing a copy of the value , the function receives  the address and modify it and return it 

# pointer aliasing 
  when the multiple pointer refers to the same variable , they are called alias 

eg :) 
          x -- p1 
          |
          |----  p2 

  Aliasing is useful but makes code harder to track because a variable can be modified through different pointers.



# Pointers in the flag Package 

What is the flag package?

The flag package is a standard Go package used to read command-line arguments passed to a program.
It allows users to configure program behavior without modifying the source code.

# Why do flag.Bool(), flag.String(), etc. return pointers?

Functions like flag.String(), flag.Bool(), and flag.Int() return pointers to variables instead of returning the values directly.

This is because, at the time these functions(flag.string,flag.int) are called, the command-line arguments have not yet been parsed.

Returning a pointer allows the flag package to update the underlying variable later when the arguments are parsed.

# Why do we use the dereference operator (*)?

Since these functions return pointers (*string, *bool, *int, etc.), the actual value is obtained by dereferencing the pointer using the * operator.
Dereferencing accesses the value stored at the memory location the pointer refers to.

# What does flag.Parse() do?

flag.Parse() reads the command-line arguments supplied when the program starts.
It matches the provided arguments with the flags that were defined.
It then updates the values stored in the corresponding flag variables.
Therefore, flag.Parse() must be called before using any flag values.

# What happens if flag.Parse() is not called?

The command-line arguments are never processed.
Every flag variable keeps its default value.
Any values passed by the user on the command line are ignored.

# Why is this design used?

The flag package separates flag declaration from flag parsing.
Returning pointers allows the package to modify the same variables after parsing instead of creating copies.
This design makes command-line configuration flexible and efficient.

Interview Answer 

"The flag package in Go is used to read command-line arguments. Functions such as flag.String(), flag.Bool(), and flag.Int() return pointers to the flag variables rather than the values themselves. This is because the command-line arguments have not yet been parsed when the flags are declared. After defining the flags, we call flag.Parse(), which reads the command-line arguments and updates the values of those variables. Since the functions return pointers, we use the dereference operator (*) to access the actual values. If flag.Parse() is not called, the flags retain their default values, and any command-line arguments provided by the user are ignored."


