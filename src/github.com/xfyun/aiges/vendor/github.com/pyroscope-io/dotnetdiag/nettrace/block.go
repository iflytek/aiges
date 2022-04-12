package nettrace

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"unsafe"
)

// BlobBlock contains a set of Blobs.
type BlobBlock struct {
	Header BlobBlockHeader
	// Payload contains Blobs - serialized fragments which represent
	// actual events and metadata records.
	Payload *bytes.Buffer

	compressed bool
	// lastHeader holds the most recent Blob header: this is used if the
	// block uses compressed headers format.
	lastHeader    BlobHeader
	extractHeader func(*Blob) error

	p *Parser
}

type BlobBlockHeader struct {
	// Size of the header including this field.
	Size  int16
	Flags int16
	// MinTimestamp specifies the minimum timestamp of any event in this block.
	MinTimestamp int64
	// MinTimestamp specifies the maximum timestamp of any event in this block.
	MaxTimestamp int64
	// Padding: optional reserved space to reach Size bytes.
}

const blockHeaderSize = int32(unsafe.Sizeof(BlobBlockHeader{}))

type Blob struct {
	Header BlobHeader
	// Payload contains Event and Metadata records.
	Payload *bytes.Buffer
	sorted  bool
}

// BlobHeader used for both compressed and uncompressed blobs.
type BlobHeader struct {
	_                 int32 // Unused: Size specifies record size not counting this field.
	MetadataID        int32
	SequenceNumber    int32
	ThreadID          int64
	CaptureThreadID   int64
	CaptureProcNumber int32
	StackID           int32
	TimeStamp         int64
	ActivityID        [16]byte
	RelatedActivityID [16]byte
	PayloadSize       int32
}

type StackBlock struct {
	Stacks []Stack
}

type Stack struct {
	ID   int32
	Data []byte
}

type SequencePointBlock struct {
	TimeStamp int64
	Threads   []Thread
}

type Thread struct {
	ThreadID       int64
	SequenceNumber int32
}

type compressedHeaderFlag byte

// The section lists compressed header flags: if otherwise specified, every
// flag is set when an appropriate header member should be read from stream.
const (
	flagMetadataID compressedHeaderFlag = 1 << iota
	// Specifies that CaptureThreadID, CaptureProcNumber and SequenceNumber
	// values are to be read from the stream: Incremented SequenceNumber should
	// be added to the previous value. If the flag is not set, and MetadataID
	// field is zero, incremented previous SequenceNumber should be used.
	flagCaptureThreadAndSequence
	flagThreadID
	flagStackID
	flagActivityID
	flagRelatedActivityID
	// Specifies that the events are sorted.
	flagIsSorted
	flagPayloadSize
)

func (b *BlobBlock) IsCompressed() bool { return b.compressed }

func (b *Blob) IsSorted() bool { return b.sorted }

func BlobBlockFromObject(o Object) (*BlobBlock, error) {
	p := Parser{Buffer: o.Payload}
	var b BlobBlock
	p.Read(&b.Header)
	b.compressed = b.Header.Flags&0x0001 != 0
	if b.compressed {
		b.extractHeader = b.readHeaderCompressed
	} else {
		b.extractHeader = b.readHeader
	}
	// Skip header padding.
	p.Skip(int(b.Header.Size) - int(blockHeaderSize))
	b.Payload = o.Payload
	b.p = &p
	return &b, p.Err()
}

func (b *BlobBlock) Next(blob *Blob) error {
	err := b.extractHeader(blob)
	switch {
	default:
	case errors.Is(err, io.EOF):
		return io.EOF
	case err != nil:
		return err
	}
	pb := b.Payload.Next(int(blob.Header.PayloadSize))
	blob.Payload = bytes.NewBuffer(pb)
	return nil
}

func (b *BlobBlock) readHeader(blob *Blob) error {
	b.p.Read(blob)
	// In the context of an EventBlock the low 31 bits are a foreign key to the
	// event's metadata. In the context of a metadata block the low 31 bits are
	// always zeroed. The high bit is the IsSorted flag.
	blob.Header.MetadataID &= 0x7FFF
	blob.sorted = uint32(blob.Header.MetadataID)&0x8000 == 0
	return b.p.Err()
}

func (b *BlobBlock) readHeaderCompressed(blob *Blob) error {
	blob.Header = b.lastHeader
	var flags compressedHeaderFlag
	b.p.Read(&flags)
	blob.sorted = flags&flagIsSorted != 0
	if flags&flagMetadataID != 0 {
		blob.Header.MetadataID = int32(b.p.Uvarint())
	}
	if flags&flagCaptureThreadAndSequence != 0 {
		blob.Header.SequenceNumber = int32(b.p.Uvarint()) + 1
		blob.Header.CaptureThreadID = int64(b.p.Uvarint())
		blob.Header.CaptureProcNumber = int32(b.p.Uvarint())
	} else if blob.Header.MetadataID != 0 {
		blob.Header.SequenceNumber++
	}
	if flags&flagThreadID != 0 {
		blob.Header.ThreadID = int64(b.p.Uvarint())
	}
	if flags&flagStackID != 0 {
		blob.Header.StackID = int32(b.p.Uvarint())
	}
	blob.Header.TimeStamp += int64(b.p.Uvarint())
	if flags&flagActivityID != 0 {
		b.p.Read(&blob.Header.ActivityID)
	}
	if flags&flagRelatedActivityID != 0 {
		b.p.Read(&blob.Header.RelatedActivityID)
	}
	if flags&flagPayloadSize != 0 {
		blob.Header.PayloadSize = int32(b.p.Uvarint())
	}
	b.lastHeader = blob.Header
	return b.p.Err()
}

func StackBlockFromObject(o Object) (*StackBlock, error) {
	var b StackBlock
	var count, size, id int32
	p := Parser{Buffer: o.Payload}
	p.Read(&id)
	p.Read(&count)
	b.Stacks = make([]Stack, count)
	for i := int32(0); i < count; i++ {
		p.Read(&size)
		b.Stacks[i] = Stack{
			ID:   id + i,
			Data: o.Payload.Next(int(size)),
		}
	}
	return &b, p.Err()
}

func SequencePointBlockFromObject(o Object) (*SequencePointBlock, error) {
	var b SequencePointBlock
	var count int32
	p := Parser{Buffer: o.Payload}
	p.Read(&b.TimeStamp)
	p.Read(&count)
	b.Threads = make([]Thread, count)
	for i := int32(0); i < count; i++ {
		var t Thread
		p.Read(&t)
		b.Threads[i] = t
	}
	return &b, p.Err()
}

func (s Stack) InstructionPointers(size int32) []uint64 {
	var address func([]byte, int) uint64
	if size == 8 {
		address = address64
	} else {
		address = address32
	}
	n := make([]uint64, len(s.Data)/int(size))
	for i := 0; i < len(n); i++ {
		n[len(n)-1-i] = address(s.Data, i)
	}
	return n
}

func address64(b []byte, i int) uint64 {
	return binary.LittleEndian.Uint64(b[i*8 : (i+1)*8])
}

func address32(b []byte, i int) uint64 {
	return uint64(binary.LittleEndian.Uint32(b[i*4 : (i+1)*4]))
}
