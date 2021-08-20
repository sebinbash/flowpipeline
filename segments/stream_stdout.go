package segments

import (
	"fmt"
	"log"
	"os"
	"sync"

	"google.golang.org/protobuf/encoding/protojson"
)

type StdOut struct {
	BaseSegment
}

func (segment StdOut) New(config map[string]string) Segment {
	return &StdOut{}
}

func (segment *StdOut) Run(wg *sync.WaitGroup) {
	defer func() {
		close(segment.Out)
		wg.Done()
	}()
	for msg := range segment.In {
		data, err := protojson.Marshal(msg)
		if err != nil {
			log.Printf("[warning] StdOut: Skipping a flow, failed to recode protobuf as JSON: %v", err)
			continue
		}
		fmt.Fprintln(os.Stdout, string(data))
		segment.Out <- msg
	}
}

func init() {
	segment := &StdOut{}
	RegisterSegment("stdout", segment)
}
