# Control flow :) The mechanism which that determine the order in which the individual code statement or the instruction  or the fucntions calls and executed 
    -> for 
    -> if 
    -> switch 
    syntax 
    switch expression {
    
    case value1:
        // code
    case value2:
        // code
    default:
       // code
    }

    eg 
    switch coinflip(){
    case "Head":
        head++
    case "tails":
        tail++
    default:
        fmt.println("Done in doing the things")

# Name types 
    a named type is simply giving name to a type and type tell go what kind of data a variable can hold 

"
# pointer is a variable which is used to store the memeory address of another variable 

# function :) is an independent block of code it doesnot belongs to the type 
syntax 
    func add (_){}
    
# method is a special type of function which is attached to the specific type 
syntax
    func (reciver *typename)name {}

# package 
    package in go is a collection of a go source files that are cmpilled together and provide the relevant functionality 
    all the go files are in the same package declare the same package naem 
    package help to organise code and prompt modularity and enabel the code resuse it 

