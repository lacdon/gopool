# gopool
go routine pool


## Usage 
```go
package main
import "github.com/lacdon/gopool"
import "fmt"
import "time"

func main() {
    poolWorkerSize := uint(1000)
    p, err := gopool.New(poolWorkerSize)
    if err != nil {
    panic(err)
    }
    for i:= 0 ; i <= 99999; i++ {
        job := func() {
            fmt.Println(time.Now().String())
        }
        go p.Dispatch(job)
    }
}
```