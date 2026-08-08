package main

import "fmt"

func main() {
    fmt.Println("Pipeline Example")

    // Stage 1: Producer
    number := make(chan int)

    go func() {
        for i := 0; i < 100; i++ {
            number <- i
        }
        close(number)
    }()

    // Stage 2: Square
    square := make(chan int)

    go func() {
        for x := range number {
            square <- x * x
        }
        close(square)
    }()

    // Stage 3: Consumer
    for y := range square {
        fmt.Println(y)
    }
}