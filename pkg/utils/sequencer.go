package utils

type Sequencer struct {
    current uint8
    max uint8
}

func NewSequencer(max uint8) Sequencer {
    return Sequencer{
        current: max,
        max: max,
    }
}

func (s *Sequencer) SetCurrentValue(current uint8) {
	s.current = current
}

func (s *Sequencer) Next() uint8 {
    s.current++
    return s.current
}

func (s *Sequencer) Current() uint8 {
    return s.current
}

