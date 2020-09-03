
package main

import (
        "encoding/csv"
        "fmt"
        "io"
        "log"
        "os"
        "time"
        "encoding/json"
        _ "unsafe"
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

type TsInfo struct {
    Sent []float64
    Recv []float64
}

var m map[string][]float64
var datetime_template string

func ParseRecord(ch chan []string, done chan bool) {

    for{
        record := <-ch

        if record[0] == "done"{
            done <- true
            break
        }

        parse_time, err := time.Parse(datetime_template, record[0])
        if err != nil {
            fmt.Println(err)
        }

        num := float64(parse_time.UnixNano() / 1000)/1000000.0
        msg := record[1]

        m[msg] = append(m[msg], num)
    }
}

func ParseFile(fname string, isRecv bool) {

    // Open the file
    csvfile, err := os.Open(fname)
    if err != nil {
        log.Fatalln("Couldn't open the csv file", err)
    }
    defer csvfile.Close()

    // Parse the file
    r := csv.NewReader(csvfile)
    r.Comma = ';'

    ch := make(chan []string, 1_000)
    done := make(chan bool, 1)

    go ParseRecord(ch, done)

    // Iterate through the records
    for {
        // Read each record from csv
        record, err := r.Read()

        if err == io.EOF {
            t := make([]string, 0)
            t = append(t, "done")
            ch <- t

            close(ch)
            break
        }
        if err != nil {
            log.Fatal(err)
        }

        ch <- record[0:2]
    }

    <-done

    PrintMemUsage()
}

func main() {
    PrintMemUsage()
    m = make(map[string][]float64)

    datetime_template = "2006/01/02-15:04:05.000000"

    ParseFile("ts_sorted_10M_lines_sent.csv", false)

    for msg, val := range m {
        m[msg] = append(val, 0)
    }

    ParseFile("ts_sorted_10M_lines_recv.csv", true)

    if false {
        b, err := json.MarshalIndent(m, "", "  ")
        if err != nil {
            fmt.Println("error:", err)
        }
        fmt.Print(string(b))

    }

    //fmt.Printf("m: %T, %d\n", m, unsafe.Sizeof(m))
    PrintMemUsage()
}

