package freader

import (
	"bufio"
	"log"
	"os"
	"runtime/debug"
)

type ChannelBuf struct {
	Str     string
	channel chan string
	eof     bool
}

func (this *ChannelBuf) Next() bool {
	if !this.eof {
		if line, ok := <-this.channel; ok {
			this.Str = line
			return true
		} else {
			this.eof = true
		}
	}
	return false
}

func Open(fname string) *ChannelBuf {
	ch := make(chan string, 32) // Some arbitrary amount of read ahead

	file, err := os.Open(fname)
	if err != nil {
		debug.PrintStack()
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)

	go func() {
		defer func() {
			file.Close()
			close(ch)
		}()
		for scanner.Scan() {
			ch <- scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			debug.PrintStack()
			log.Fatal(err)
		}
	}()

	return &ChannelBuf{"", ch, false}
}
