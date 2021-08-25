// Runs flows through a filter and forwards only matching flows.
package flowfilter

import (
	"log"
	"sync"

	"github.com/bwNetFlow/flowfilter/parser"
	"github.com/bwNetFlow/flowfilter/visitors"
	"github.com/bwNetFlow/flowpipeline/segments"
)

type FlowFilter struct {
	segments.BaseSegment
	Filter string

	expression *parser.Expression
}

func (segment FlowFilter) New(config map[string]string) segments.Segment {
	var err error

	newSegment := &FlowFilter{
		Filter: config["filter"],
	}

	newSegment.expression, err = parser.Parse(config["filter"])
	if err != nil {
		log.Printf("[error] FlowFilter: Syntax error in filter expression: %v", err)
		return nil
	}
	return newSegment
}

func (segment *FlowFilter) Run(wg *sync.WaitGroup) {
	defer func() {
		close(segment.Out)
		wg.Done()
	}()

	log.Printf("[info] FlowFilter: Using filter expression: %s", segment.Filter)

	filter := &visitors.Filter{}
	for msg := range segment.In {
		if match, err := filter.CheckFlow(segment.expression, msg); match {
			if err != nil {
				log.Printf("[error] FlowFilter: Semantic error in filter expression: %v", err)
				continue // TODO: introduce option on-error action, current state equals 'drop'
			}
			segment.Out <- msg
		}
	}
}

func init() {
	segment := &FlowFilter{}
	segments.RegisterSegment("flowfilter", segment)
}