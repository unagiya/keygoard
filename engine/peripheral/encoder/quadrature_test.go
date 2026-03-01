package encoder

import "testing"

func TestDecoder_CW(t *testing.T) {
	// CW: 00→01→11→10→00
	d := Decoder{}
	steps := []struct {
		a, b bool
		want Direction
	}{
		{false, true, DirCW},  // 00→01
		{true, true, DirCW},   // 01→11
		{true, false, DirCW},  // 11→10
		{false, false, DirCW}, // 10→00
	}
	for i, s := range steps {
		got := d.Update(s.a, s.b)
		if got != s.want {
			t.Errorf("step %d: Update(%v, %v) = %d, want %d", i, s.a, s.b, got, s.want)
		}
	}
}

func TestDecoder_CCW(t *testing.T) {
	// CCW: 00→10→11→01→00
	d := Decoder{}
	steps := []struct {
		a, b bool
		want Direction
	}{
		{true, false, DirCCW},  // 00→10
		{true, true, DirCCW},   // 10→11
		{false, true, DirCCW},  // 11→01
		{false, false, DirCCW}, // 01→00
	}
	for i, s := range steps {
		got := d.Update(s.a, s.b)
		if got != s.want {
			t.Errorf("step %d: Update(%v, %v) = %d, want %d", i, s.a, s.b, got, s.want)
		}
	}
}

func TestDecoder_NoChange(t *testing.T) {
	d := Decoder{}
	// 同じ状態を連続で送ると DirNone
	if got := d.Update(false, false); got != DirNone {
		t.Errorf("same state: got %d, want DirNone", got)
	}
	if got := d.Update(false, false); got != DirNone {
		t.Errorf("same state again: got %d, want DirNone", got)
	}
}

func TestDecoder_InvalidTransition(t *testing.T) {
	d := Decoder{}
	// 00→11 はスキップ（不正遷移）→ DirNone
	if got := d.Update(true, true); got != DirNone {
		t.Errorf("skip 00→11: got %d, want DirNone", got)
	}
}

func TestDecoder_MultipleCycles(t *testing.T) {
	d := Decoder{}
	// CW を 2 周
	cw := []struct{ a, b bool }{
		{false, true}, {true, true}, {true, false}, {false, false},
	}
	for cycle := range 2 {
		for _, s := range cw {
			if got := d.Update(s.a, s.b); got != DirCW {
				t.Errorf("cycle %d: got %d, want DirCW", cycle, got)
			}
		}
	}
}
