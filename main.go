package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	flag "github.com/spf13/pflag"
)

func main() {
	var writeSize uint64 = 4 * 1024 * 1024 * 1024

	var flagFull bool
	flag.BoolVarP(&flagFull, "full", "f", false, "Fullfill a file until the drive is full")
	flag.Parse()
	if flagFull {
		writeSize = math.MaxUint64
	}

	f, err := os.CreateTemp("", "fill-drive")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(f.Name()) // clean up

	const BUF_SIZE = 4096
	var buf [BUF_SIZE]byte
	for i := 0; i < BUF_SIZE; i++ {
		buf[i] = 0
	}

	rand.Seed(time.Now().UnixNano())

	var byteWritten uint64 = 0
	for byteWritten < writeSize {
		rand.Read(buf[:])
		n, err := f.Write(buf[:])
		if err == nil {
			err = f.Sync()
		}
		if err != nil {
			log.Print(err)
			break
		}
		byteWritten += uint64(n)
		fmt.Printf("\r%d bytes written", byteWritten)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
