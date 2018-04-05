package faketerm

// diff.go provides functions that highlight differences between fake terminals.

import (
	"bytes"
	"fmt"
	"image"
	"reflect"

	"github.com/kylelemons/godebug/pretty"
	"github.com/mum4k/termdash/cell"
)

// optDiff is used to display differences in cell options.
type optDiff struct {
	// point indicates the cell with the differing options.
	point image.Point

	got  *cell.Options
	want *cell.Options
}

// Diff compares the two terminals, returning an empty string if there is not
// difference. If a difference is found, returns a human readable description
// of the differences.
func Diff(want, got *Terminal) string {
	if reflect.DeepEqual(want, got) {
		return ""
	}

	var b bytes.Buffer
	b.WriteString("found differences between the two fake terminals.\n")
	b.WriteString("   got:\n")
	b.WriteString(got.String())
	b.WriteString("  want:\n")
	b.WriteString(want.String())
	b.WriteString("  diff (unexpected cells highlighted with rune '࿃')\n")
	b.WriteString("  note - this excludes cell options:\n")

	size := got.Size()
	var optDiffs []*optDiff
	for row := 0; row < size.Y; row++ {
		for col := 0; col < size.X; col++ {
			gotCell := got.BackBuffer()[col][row]
			wantCell := want.BackBuffer()[col][row]
			r := gotCell.Rune
			if r != wantCell.Rune {
				r = '࿃'
			} else if r == 0 {
				r = ' '
			}
			b.WriteRune(r)

			if !reflect.DeepEqual(gotCell.Opts, wantCell.Opts) {
				optDiffs = append(optDiffs, &optDiff{
					point: image.Point{col, row},
					got:   gotCell.Opts,
					want:  wantCell.Opts,
				})
			}
		}
		b.WriteRune('\n')
	}

	if len(optDiffs) > 0 {
		b.WriteString("  Found differences in options on some of the cells:\n")
		for _, od := range optDiffs {
			if diff := pretty.Compare(od.want, od.got); diff != "" {
				b.WriteString(fmt.Sprintf("cell %v, diff (-want +got):\n%s\n", od.point, diff))
			}
		}
	}
	return b.String()
}
