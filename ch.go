package main


import (
        "fmt"
     _   "time"
        "runtime"
       )

func PrintMemUsage() {
        var m runtime.MemStats
        runtime.ReadMemStats(&m)
        // For info on each, see: https://golang.org/pkg/runtime/#MemStats
        //fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
        //fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
        fmt.Printf("\tSys memory = %v MiB", bToMb(m.Sys))
        fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
    return b / 1024 / 1024
}


func worker(ch chan string, done chan bool) {
    //time.Sleep(10*time.Second)

    for{
        s := <-ch

        if s == "done"{
            done <- true
            //break
        }

    }

}

func main() {
    ch := make(chan string, 100_000_000)
    done := make(chan bool, 1)

    go worker(ch, done)

    PrintMemUsage()
    for i := 0; i < 10_000_000 ; i++ {
        //fmt.Println(i)
        ch <- string(i) + "kkjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjk"
    }

    ch <- "done"

    close(ch)

    <-done

    fmt.Println("Done")
    PrintMemUsage()
}



