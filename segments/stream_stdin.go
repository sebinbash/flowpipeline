package segments

import (
	"bufio"
	"bytes"
	flow "github.com/bwNetFlow/protobuf/go"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	// "io"
	"os"
	"sync"
)

type StdIn struct {
	BaseSegment
}

func (segment StdIn) New(config map[string]string) Segment {
	return &StdIn{}
}

func (segment *StdIn) Run(wg *sync.WaitGroup) {
	defer func() {
		close(segment.out)
		wg.Done()
	}()
	fromStdin := make(chan []byte)
	go func() {
		for {
			rdr := bufio.NewReader(os.Stdin)
			line, err := rdr.ReadBytes('\n')
			if err != nil {
				log.Printf("[warning] StdIn: Skipping a flow, could not read line from stdin: %v", err)
				continue
			}
			if len(line) == 0 {
				continue
			}
			fromStdin <- bytes.TrimSuffix(line, []byte("\n"))
		}
	}()
	for {
		select {
		case msg, ok := <-segment.in:
			if !ok {
				return
			}
			segment.out <- msg
		case line := <-fromStdin:
			msg := &flow.FlowMessage{}
			err := protojson.Unmarshal(line, msg)
			if err != nil {
				log.Printf("[warning] StdIn: Skipping a flow, failed to recode stdin to protobuf: %v", err)
				continue
			}
			segment.out <- msg
		}
	}
}