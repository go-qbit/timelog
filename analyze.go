package timelog

import (
	"bytes"
	"fmt"
	"io"
	"time"
)

// A message for unknown actions.
const UnknownActionCaption = "Working"

// Contains calculated information about an action. See TlEntity.Analyze()
type Action struct {
	Duration time.Duration
	Message  string
	Children []*Action
}

/*
Returns a string with timings as a tree:
	1.000453077s	Level 1
		1.277µs	Working
		800.303067ms	Level 2.1
			500.166711ms	Working
			300.122218ms	Level 3
			14.138µs	Working
		843ns	Working
		200.143999ms	Level 2.2
		3.891µs	Working
*/
func (a *Action) String() string {
	buf := &bytes.Buffer{}
	a.Print(buf, "\t")

	return buf.String()
}

// Prints hierarchic tree of timings to w
func (a *Action) Print(w io.Writer, offsetStr string) error {
	buf := &bytes.Buffer{}
	a.print(buf, offsetStr, 0)
	_, err := w.Write(buf.Bytes())

	return err
}

func (a *Action) print(w io.Writer, offsetStr string, offset int) error {
	for i := 0; i < offset; i++ {
		if _, err := w.Write([]byte(offsetStr)); err != nil {
			return err
		}
	}

	if _, err := w.Write([]byte(a.Duration.String() + "\t" + a.Message + "\n")); err != nil {
		return err
	}

	for _, child := range a.Children {
		if err := child.print(w, offsetStr, offset+1); err != nil {
			return err
		}
	}

	return nil
}

// Returns calculated information about actions.
func (e *TlEntity) Analyze() *Action {
	res := &Action{
		Duration: e.finish.Sub(e.start),
		Message:  getMessage(e.message),
	}

	if len(e.children) > 0 {
		beforeDuration := e.children[0].start.Sub(e.start)
		if beforeDuration > 0 {
			res.Children = append(res.Children, &Action{
				Duration: beforeDuration,
				Message:  UnknownActionCaption,
			})
		}

		for i, child := range e.children {
			res.Children = append(res.Children, child.Analyze())

			var afterDuration time.Duration
			if i == len(e.children)-1 { // A last element
				afterDuration = e.finish.Sub(child.finish)
			} else {
				afterDuration = e.children[i+1].start.Sub(child.finish)
			}
			if afterDuration > 0 {
				res.Children = append(res.Children, &Action{
					Duration: afterDuration,
					Message:  UnknownActionCaption,
				})
			}
		}
	}

	return res
}

func getMessage(s interface{}) string {
	switch m := s.(type) {
	case string:
		return m
	case fmt.Stringer:
		return m.String()
	case fmt.GoStringer:
		return m.GoString()
	default:
		return fmt.Sprintf("%v", m)
	}
}
