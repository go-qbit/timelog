package timelog

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"time"
)

// A message for unknown actions.
const UnknownActionCaption = "Working"

// Contains calculated information about an action. See TlEntity.Analyze()
type Action struct {
	StartOffset      time.Duration
	Duration         time.Duration
	DurationPercents float64
	Message          string
	Children         []*Action
}

/*
Returns a string with timings as a tree:
	[100.0000%] Level 1		(1.000399975s: 0s ⟼ 1.000399975s)
		[0.0001%] Working		(824ns: 0s ⟼ 824ns)
		[79.9963%] Level 2.1		(800.282978ms: 824ns ⟼ 800.283802ms)
			[62.4931%] Working		(500.121338ms: 0s ⟼ 500.121338ms)
			[37.5065%] Level 3		(300.157855ms: 500.121338ms ⟼ 800.279193ms)
			[0.0005%] Working		(3.785µs: 800.279193ms ⟼ 800.282978ms)
		[0.0001%] Working		(765ns: 800.283802ms ⟼ 800.284567ms)
		[20.0032%] Level 2.2		(200.11168ms: 800.284567ms ⟼ 1.000396247s)
		[0.0004%] Working		(3.728µs: 1.000396247s ⟼ 1.000399975s)
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

func (a *Action) print(buf *bytes.Buffer, offsetStr string, offset int) {
	for i := 0; i < offset; i++ {
		buf.WriteString(offsetStr)
	}

	buf.WriteByte('[')
	buf.WriteString(strconv.FormatFloat(a.DurationPercents, 'f', 4, 64))
	buf.WriteString("%] ")
	buf.WriteString(a.Message)
	buf.WriteString("\t\t(")
	buf.WriteString(a.Duration.String())
	buf.WriteString(": ")
	buf.WriteString(a.StartOffset.String())
	buf.WriteString(" ⟼ ")
	buf.WriteString((a.StartOffset + a.Duration).String())
	buf.WriteString(")\n")

	for _, child := range a.Children {
		child.print(buf, offsetStr, offset+1)
	}
}

// Returns calculated information about actions.
func (e *TlEntity) Analyze() *Action {
	res := &Action{
		Duration:         e.finish.Sub(e.start),
		DurationPercents: 100.0,
		Message:          getMessage(e.message),
	}

	if len(e.children) > 0 {
		beforeDuration := e.children[0].start.Sub(e.start)
		if beforeDuration > 0 {
			res.Children = append(res.Children, &Action{
				Duration:         beforeDuration,
				StartOffset:      0,
				DurationPercents: float64(beforeDuration) / float64(res.Duration) * 100.0,
				Message:          UnknownActionCaption,
			})
		}

		for i, child := range e.children {
			analyzedChild := child.Analyze()
			analyzedChild.StartOffset = child.start.Sub(e.start)
			analyzedChild.DurationPercents = float64(analyzedChild.Duration) / float64(res.Duration) * 100.0

			res.Children = append(res.Children, analyzedChild)

			var afterDuration, startOffset time.Duration
			if i == len(e.children)-1 { // A last element
				afterDuration = e.finish.Sub(child.finish)
				startOffset = e.children[i].finish.Sub(e.start)
			} else {
				afterDuration = e.children[i+1].start.Sub(child.finish)
				startOffset = child.finish.Sub(e.start)
			}
			if afterDuration > 0 {
				res.Children = append(res.Children, &Action{
					Duration:         afterDuration,
					StartOffset:      startOffset,
					DurationPercents: float64(afterDuration) / float64(res.Duration) * 100.0,
					Message:          UnknownActionCaption,
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
