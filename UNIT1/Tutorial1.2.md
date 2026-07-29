# Finding the Duplicate lines 
    
    map holds the key vale pair and provies the contant time for updating retriving and storin the data 
    
    initilizing the map 
        const := make(map[string]int)
    the key map be any type but has one constrains that is it  must be compared with the == 

    bufio.scanner  used to read the data line by line words by words fom input source like file or standar inputs 

#  format verbs 
# need ?? 
fmt.Printf("Hello %s, your age is %d\n", name, age)


    format verbs are the placeholder that specifies/ tells how the value should be displayed in the progrma and it starts witht he % [percent]
    Some example 
    %d decimal 
    %x, %o, %b integer in hexadecimal, octal, binary
    %f, %g, %e floating-point number: 3.141593 3.141592653589793 3.141593e+00
    %t boolean: true or false
    %c rune (Unicode code point)
    %s string
    %q quoted string "abc" or rune 'c'
    %v any value in a natural format
    %T type of any value
    %% literal percent sign (no operand)
  
  # ways 1 to open the file 
    another way is f , errr := os.open(os.Args[0]) and then close the resouce 
    f.close():
    f = *os.file  adn other is error 

# way 2 
    using the io/ioutiil 
    
    need the file name and file name can be taken by os.Args[0] 
