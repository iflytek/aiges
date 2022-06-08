package nettrace

import (
	"errors"
	"io"
)

type Stream struct {
	dec *Decoder

	EventHandler              func(*Blob) error
	MetadataHandler           func(*Metadata) error
	StackBlockHandler         func(*StackBlock) error
	SequencePointBlockHandler func(*SequencePointBlock) error
}

func NewStream(r io.Reader) *Stream {
	return &Stream{dec: NewDecoder(r)}
}

func (s *Stream) Open() (*Trace, error) {
	return s.dec.OpenTrace()
}

func (s *Stream) Next() error {
	var o Object
	if err := s.dec.Decode(&o); err != nil {
		return err
	}

	switch o.Type {
	case ObjectTypeSPBlock:
		if s.SequencePointBlockHandler == nil {
			return nil
		}
		block, err := SequencePointBlockFromObject(o)
		if err != nil {
			return err
		}
		return s.SequencePointBlockHandler(block)

	case ObjectTypeStackBlock:
		if s.StackBlockHandler == nil {
			return nil
		}
		block, err := StackBlockFromObject(o)
		if err != nil {
			return err
		}
		return s.StackBlockHandler(block)

	case ObjectTypeEventBlock:
		if s.EventHandler == nil {
			return nil
		}
		block, err := BlobBlockFromObject(o)
		if err != nil {
			return err
		}
		var blob Blob
		for {
			err = block.Next(&blob)
			switch {
			case err == nil:
			case errors.Is(err, io.EOF):
				return nil
			default:
				return err
			}
			if err = s.EventHandler(&blob); err != nil {
				return err
			}
		}

	case ObjectTypeMetadataBlock:
		if s.MetadataHandler == nil {
			return nil
		}
		block, err := BlobBlockFromObject(o)
		if err != nil {
			return err
		}
		var blob Blob
		for {
			err = block.Next(&blob)
			switch {
			case err == nil:
			case errors.Is(err, io.EOF):
				return nil
			default:
				return err
			}
			md, err := MetadataFromBlob(blob)
			if err != nil {
				return err
			}
			if err = s.MetadataHandler(md); err != nil {
				return err
			}
		}

	case ObjectTypeTrace:
		return ErrUnexpectedObjectType

	default:
		return ErrInvalidObjectType
	}
}
