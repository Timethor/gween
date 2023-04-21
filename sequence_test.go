package gween

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tanema/gween/ease"
)

func TestSequenceNew(t *testing.T) {
	seq := NewSequence(New(0, 1, 1, ease.Linear))

	current, finishedTween, seqFinished := seq.Update(0.0)
	assert.Equal(t, float32(0), current)
	assert.False(t, finishedTween)
	assert.False(t, seqFinished)
	assert.Equal(t, 0, seq.d.Index)
}

func TestSequence_Update(t *testing.T) {
	seq := NewSequence(
		New(0, 1, 1, ease.Linear),
		New(1, 2, 1, ease.Linear),
	)

	current, finishedTween, seqFinished := seq.Update(0.5)
	assert.Equal(t, float32(0.5), current)
	assert.False(t, finishedTween)
	assert.False(t, seqFinished)
	assert.Equal(t, 0, seq.d.Index)
}

func TestSequence_Reset(t *testing.T) {
	seq := NewSequence(
		New(0, 1, 1, ease.Linear),
		New(1, 2, 1, ease.Linear),
	)

	seq.Update(1.5)
	seq.Reset()
	assert.Equal(t, 0, seq.d.Index)
	assert.Equal(t, float32(0.0), seq.d.Tweens[0].d.Time)
	assert.Equal(t, float32(0.0), seq.d.Tweens[0].d.Overflow)
	assert.Equal(t, float32(0.0), seq.d.Tweens[1].d.Time)
	assert.Equal(t, float32(0.0), seq.d.Tweens[1].d.Overflow)
}

func TestSequence_CompleteFirst(t *testing.T) {
	seq := NewSequence(
		New(0, 1, 1, ease.Linear),
		New(1, 2, 1, ease.Linear),
	)

	current, finishedTween, seqFinished := seq.Update(1.0)
	assert.Equal(t, float32(1.0), current)
	assert.True(t, finishedTween)
	assert.False(t, seqFinished)
	assert.Equal(t, 1, seq.d.Index)
}

func TestSequence_OverflowSecond(t *testing.T) {
	seq := NewSequence(
		New(0, 1, 1, ease.Linear),
		New(1, 2, 1, ease.Linear),
	)

	current, finishedTween, seqFinished := seq.Update(1.5)
	assert.Equal(t, float32(1.5), current)
	assert.True(t, finishedTween)
	assert.False(t, seqFinished)
	assert.Equal(t, 1, seq.d.Index)
}

func TestSequence_OverflowAndComplete(t *testing.T) {
	seq := NewSequence(
		New(0, 1, 1, ease.Linear),
		New(1, 2, 1, ease.Linear),
		New(2, 3, 1, ease.Linear),
	)

	current, finishedTween, seqFinished := seq.Update(3.5)
	assert.Equal(t, float32(3.0), current)
	assert.True(t, finishedTween)
	assert.True(t, seqFinished)
	assert.Equal(t, 3, seq.d.Index)
}

func TestSequence_Loops(t *testing.T) {
	seq := NewSequence(
		New(0, 1, 1, ease.Linear),
		New(1, 2, 1, ease.Linear),
		New(2, 3, 1, ease.Linear),
	)
	seq.SetLoop(2)
	current, finishedTween, seqFinished := seq.Update(5.25)
	assert.Equal(t, float32(2.25), current)
	assert.True(t, finishedTween)
	assert.False(t, seqFinished)
	assert.Equal(t, 1, seq.d.LoopRemaining)
	assert.Equal(t, 2, seq.d.Index)

	current, finishedTween, seqFinished = seq.Update(0.75)
	assert.Equal(t, float32(3), current)
	assert.True(t, finishedTween)
	assert.True(t, seqFinished)
	assert.Equal(t, 0, seq.d.LoopRemaining)
	assert.Equal(t, 3, seq.d.Index)
}

func TestSequence_LoopsForever(t *testing.T) {
	seq := NewSequence(
		New(0, 1, 1, ease.Linear),
		New(1, 2, 1, ease.Linear),
		New(2, 3, 1, ease.Linear),
	)
	seq.SetLoop(-1)
	current, finishedTween, seqFinished := seq.Update(3*1_000_000 + 2.25)
	assert.Equal(t, float32(2.25), current)
	assert.True(t, finishedTween)
	assert.False(t, seqFinished)
	assert.Equal(t, -1, seq.d.LoopRemaining)
	assert.Equal(t, 2, seq.d.Index)
}

func TestSequence_Yoyos(t *testing.T) {
	seq := NewSequence(
		New(0, 1, 1, ease.Linear),
		New(1, 2, 1, ease.Linear),
		New(2, 3, 1, ease.Linear),
	)

	seq.SetYoyo(true)
	current, finishedTween, seqFinished := seq.Update(5.75)
	assert.Equal(t, float32(0.25), current)
	assert.True(t, finishedTween)
	assert.False(t, seqFinished)
	assert.Equal(t, 1, seq.d.LoopRemaining)
	assert.Equal(t, 0, seq.d.Index)

	current, finishedTween, seqFinished = seq.Update(0.25)
	assert.Equal(t, float32(0), current)
	assert.True(t, finishedTween)
	assert.True(t, seqFinished)
	assert.Equal(t, 0, seq.d.LoopRemaining)
	assert.Equal(t, 0, seq.d.Index)
}

func TestSequence_YoyosAndLoops(t *testing.T) {
	seq := NewSequence(
		New(0, 1, 1, ease.Linear),
		New(1, 2, 1, ease.Linear),
		New(2, 3, 1, ease.Linear),
	)
	seq.SetYoyo(true)
	seq.SetLoop(2)
	current, finishedTween, seqFinished := seq.Update(7.25)
	assert.Equal(t, float32(1.25), current)
	assert.True(t, finishedTween)
	assert.False(t, seqFinished)
	assert.Equal(t, 1, seq.d.LoopRemaining)
	assert.Equal(t, 1, seq.d.Index)

	current, finishedTween, seqFinished = seq.Update(4.75)
	assert.Equal(t, float32(0), current)
	assert.True(t, finishedTween)
	assert.True(t, seqFinished)
	assert.Equal(t, 0, seq.d.LoopRemaining)
	assert.Equal(t, 0, seq.d.Index)
}

func TestSequence_SetReverse(t *testing.T) {
	seq := NewSequence(
		New(0, 1, 1, ease.Linear),
		New(1, 2, 1, ease.Linear),
		New(2, 3, 1, ease.Linear),
	)
	seq.SetLoop(2)

	// Normal operation
	current, finishedTween, seqFinished := seq.Update(2.25)
	assert.Equal(t, float32(2.25), current)
	assert.True(t, finishedTween)
	assert.False(t, seqFinished)
	assert.Equal(t, 2, seq.d.LoopRemaining)
	assert.Equal(t, 2, seq.d.Index)

	seq.SetReverse(true)

	// Goes in reverse
	current, finishedTween, seqFinished = seq.Update(2.0)
	assert.Equal(t, float32(0.25), current)
	assert.True(t, finishedTween)
	assert.False(t, seqFinished)
	assert.Equal(t, 2, seq.d.LoopRemaining)
	assert.Equal(t, 0, seq.d.Index)
	assert.True(t, seq.Reverse())

	// Consumes a loop at the start!, resets to the end, continues in reverse
	current, finishedTween, seqFinished = seq.Update(2.0)
	assert.Equal(t, float32(1.25), current)
	assert.True(t, finishedTween)
	assert.False(t, seqFinished)
	assert.Equal(t, 1, seq.d.LoopRemaining)
	assert.Equal(t, 1, seq.d.Index)
	assert.True(t, seq.Reverse())

	// Hits the beginning, no more loops, ends
	current, finishedTween, seqFinished = seq.Update(2.0)
	assert.Equal(t, float32(0.0), current)
	assert.True(t, finishedTween)
	assert.True(t, seqFinished)
	assert.Equal(t, 0, seq.d.LoopRemaining)
	assert.Equal(t, -1, seq.d.Index)
	assert.True(t, seq.Reverse())
}

func TestSequence_SetReverseWithYoyo(t *testing.T) {
	seq := NewSequence(
		New(0, 1, 1, ease.Linear),
		New(1, 2, 1, ease.Linear),
		New(2, 3, 1, ease.Linear),
	)
	seq.SetYoyo(true)
	seq.SetLoop(2)

	// Standard operation
	current, finishedTween, seqFinished := seq.Update(2.25)
	assert.Equal(t, float32(2.25), current)
	assert.True(t, finishedTween)
	assert.False(t, seqFinished)
	assert.Equal(t, 2, seq.d.LoopRemaining)
	assert.Equal(t, 2, seq.d.Index)

	seq.SetReverse(true)

	// Goes in reverse
	current, finishedTween, seqFinished = seq.Update(2.0)
	assert.Equal(t, float32(0.25), current)
	assert.True(t, finishedTween)
	assert.False(t, seqFinished)
	assert.Equal(t, 2, seq.d.LoopRemaining)
	assert.Equal(t, 0, seq.d.Index)

	// Consumes a loop at the start, despite not reaching the end yet, and continues
	current, finishedTween, seqFinished = seq.Update(2.0)
	assert.Equal(t, float32(1.75), current)
	assert.True(t, finishedTween)
	assert.False(t, seqFinished)
	assert.Equal(t, 1, seq.d.LoopRemaining)
	assert.Equal(t, 1, seq.d.Index)

	// Hits the end, yoyos
	current, finishedTween, seqFinished = seq.Update(2.0)
	assert.Equal(t, float32(2.25), current)
	assert.True(t, finishedTween)
	assert.False(t, seqFinished)
	assert.Equal(t, 1, seq.d.LoopRemaining)
	assert.Equal(t, 2, seq.d.Index)
	assert.True(t, seq.Reverse()) // Is in reverse

	seq.SetReverse(false) // Go forward instead

	// Hits the end again, yoyos the same
	current, finishedTween, seqFinished = seq.Update(1.5)
	assert.Equal(t, float32(2.25), current)
	assert.True(t, finishedTween)
	assert.False(t, seqFinished)
	assert.Equal(t, 1, seq.d.LoopRemaining)
	assert.Equal(t, 2, seq.d.Index)

	// Consumes a loop at the start like normal, no more loops, end
	current, finishedTween, seqFinished = seq.Update(2.5)
	assert.Equal(t, float32(0.0), current)
	assert.True(t, finishedTween)
	assert.True(t, seqFinished)
	assert.Equal(t, 0, seq.d.LoopRemaining)
	assert.Equal(t, 0, seq.d.Index)
}

func TestSequence_SetReverseAfterComplete(t *testing.T) {
	seq := NewSequence(
		New(0, 1, 1, ease.Linear),
		New(1, 2, 1, ease.Linear),
		New(2, 3, 1, ease.Linear),
	)
	seq.SetLoop(1)

	// Normal operation
	current, finishedTween, seqFinished := seq.Update(3.0)
	assert.Equal(t, float32(3.0), current)
	assert.True(t, finishedTween)
	assert.True(t, seqFinished)
	assert.Equal(t, 0, seq.d.LoopRemaining)
	assert.Equal(t, 3, seq.d.Index)

	seq.SetReverse(true)
	seq.SetLoop(1)

	// Goes in reverse
	current, finishedTween, seqFinished = seq.Update(2.0)
	assert.Equal(t, float32(1.0), current)
	assert.True(t, finishedTween)
	assert.False(t, seqFinished)
	assert.Equal(t, 1, seq.d.LoopRemaining)
	assert.Equal(t, 0, seq.d.Index)
	assert.True(t, seq.Reverse())
}

func TestSequence_Remove(t *testing.T) {
	seq := NewSequence(
		New(0, 1, 1, ease.Linear),
		New(1, 2, 1, ease.Linear),
		New(2, 3, 1, ease.Linear),
		New(3, 4, 1, ease.Linear),
		New(4, 5, 1, ease.Linear),
	)
	assert.Equal(t, 5, len(seq.d.Tweens))
	seq.Remove(2)
	assert.Equal(t, 4, len(seq.d.Tweens))
	current, finishedTween, seqFinished := seq.Update(2.5)
	assert.Equal(t, float32(3.5), current)
	assert.True(t, finishedTween)
	assert.False(t, seqFinished)
	assert.Equal(t, 2, seq.d.Index)
	seq.Remove(0)
	assert.Equal(t, 3, len(seq.d.Tweens))
	seq.Remove(0)
	assert.Equal(t, 2, len(seq.d.Tweens))
	seq.Remove(0)
	assert.Equal(t, 1, len(seq.d.Tweens))
	// Out of bound checking
	seq.Remove(0)
	assert.Equal(t, 0, len(seq.d.Tweens))
	seq.Remove(2)
	assert.Equal(t, 0, len(seq.d.Tweens))
}

func TestSequence_Has(t *testing.T) {
	seq := NewSequence()
	assert.False(t, seq.HasTweens())
	seq.Add(New(0, 5, 1, ease.Linear))
	assert.True(t, seq.HasTweens())
	seq.Remove(0)
	assert.False(t, seq.HasTweens())
	current, finishedTween, seqFinished := seq.Update(1)
	assert.Equal(t, float32(0), current)
	assert.False(t, finishedTween)
	assert.True(t, seqFinished)
}

func TestSequence_SetIndex(t *testing.T) {
	seq := NewSequence(
		New(0, 1, 1, ease.Linear),
		New(1, 2, 1, ease.Linear),
	)
	seq.SetIndex(1)
	current, finishedTween, seqFinished := seq.Update(1.5)
	assert.Equal(t, float32(2), current)
	assert.True(t, finishedTween)
	assert.True(t, seqFinished)
	assert.Equal(t, 2, seq.d.Index)
}

func TestSequence_RealWorld(t *testing.T) {

	seq := NewSequence(
		New(0, 5, 1, ease.Linear),
		New(5, 0, 1, ease.Linear),
		New(0, 2, 2, ease.Linear),
		New(2, 0, 2, ease.Linear),
		New(0, 1, 100, ease.Linear),
	)

	assert.True(t, len(seq.d.Tweens) == 5)
	seq.Remove(0)
	seq.Remove(0)
	assert.True(t, len(seq.d.Tweens) == 3)

	current, finishedTween, sequenceFinished := seq.Update(1)
	// Half-way through first tween
	assert.Equal(t, float32(1), current)
	assert.False(t, finishedTween)
	assert.False(t, sequenceFinished)

	current, finishedTween, sequenceFinished = seq.Update(1)
	// Now at the start of the second tween
	assert.Equal(t, float32(2), current)
	assert.Equal(t, seq.Index(), 1)
	assert.True(t, finishedTween)
	assert.False(t, sequenceFinished)

	_, _, sequenceFinished = seq.Update(2)
	// Now at the start of the third Tween
	assert.Equal(t, seq.Index(), 2)
	assert.False(t, sequenceFinished)

	seq.Remove(2)
	_, finishedTween, sequenceFinished = seq.Update(1)
	// Now finished because we removed the third tween and then called Sequence.Update()
	assert.False(t, finishedTween)
	assert.True(t, sequenceFinished)
}

func TestSequence_Serializes(t *testing.T) {
	sControl := NewSequence(
		New(0, 1, 1, ease.Linear),
		New(1, 2, 1, ease.InCubic),
		New(2, 3, 1, ease.InQuad),
		New(3, 4, 1, ease.InQuart),
		New(4, 5, 1, ease.InSine),
	)
	sControl.Update(2.5)

	sb, err := json.Marshal(sControl)
	assert.NoError(t, err)

	sUnmarshalled := NewSequence()

	err = json.Unmarshal(sb, sUnmarshalled)
	assert.NoError(t, err)

	assert.True(t, sControl.Equal(sUnmarshalled))
}

func TestSequence_SerializesCustomEasing(t *testing.T) {
	// Map key must match function name
	ease.EasingFunctions["MyTestEasingFunc"] = MyTestEasingFunc
	sControl := NewSequence(
		New(0, 1, 1, MyTestEasingFunc),
		New(1, 2, 1, ease.InCubic),
		New(2, 3, 1, ease.InQuad),
		New(3, 4, 1, ease.InQuart),
		New(4, 5, 1, ease.InSine),
	)
	sControl.Update(2.5)

	sb, err := json.Marshal(sControl)
	assert.NoError(t, err)

	sUnmarshalled := NewSequence()

	err = json.Unmarshal(sb, sUnmarshalled)
	assert.NoError(t, err)

	assert.True(t, sControl.Equal(sUnmarshalled))
}
