package nettrace

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

var (
	// UTF8 string 'Nettrace'
	netTraceMagic = [8]byte{0x4E, 0x65, 0x74, 0x74, 0x72, 0x61, 0x63, 0x65}
	// UTF8 string '!FastSerialization.1'
	fastSerializationMagic = [20]byte{0x21, 0x46, 0x61, 0x73, 0x74, 0x53, 0x65, 0x72, 0x69, 0x61, 0x6C, 0x69, 0x7A, 0x61, 0x74, 0x69, 0x6F, 0x6E, 0x2E, 0x31}
)

var (
	ErrInvalidNetTraceHeader    = errors.New("invalid NetTrace header")
	ErrUnsupportedObjectVersion = errors.New("unsupported object version")
	ErrInvalidObjectType        = errors.New("invalid object type")
	ErrUnexpectedObjectType     = errors.New("unexpected object type")
	ErrUnexpectedTag            = errors.New("unexpected tag")
)

const Version int32 = 4

type Decoder struct{ r *netTraceReader }

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r: &netTraceReader{
			buf:   bytes.NewBuffer(make([]byte, 0, 4<<10)),
			inner: r,
		},
	}
}

type netTraceReader struct {
	offset uint64
	inner  io.Reader
	buf    *bytes.Buffer
}

func (c *netTraceReader) Read(b []byte) (int, error) {
	n, err := io.CopyN(c.buf, c.inner, int64(len(b)))
	copy(b, c.buf.Bytes()[:n])
	c.buf.Reset()
	c.offset += uint64(n)
	return int(n), err
}

type Object struct {
	Type                 ObjectType
	Version              int32
	MinimumReaderVersion int32
	Payload              *bytes.Buffer
}

type ObjectType string

const (
	ObjectTypeTrace         ObjectType = "Trace"
	ObjectTypeEventBlock    ObjectType = "EventBlock"
	ObjectTypeMetadataBlock ObjectType = "MetadataBlock"
	ObjectTypeStackBlock    ObjectType = "StackBlock"
	ObjectTypeSPBlock       ObjectType = "SPBlock"
)

var knownObjectTypes = []ObjectType{
	ObjectTypeTrace,
	ObjectTypeEventBlock,
	ObjectTypeMetadataBlock,
	ObjectTypeStackBlock,
	ObjectTypeSPBlock,
}

func (t ObjectType) IsValid() bool {
	for _, v := range knownObjectTypes {
		if v == t {
			return true
		}
	}
	return false
}

type Tag byte

const (
	NullReference      Tag = 0x1
	BeginPrivateObject Tag = 0x5
	EndObject          Tag = 0x6
)

var objectHeaderTags = [3]Tag{
	BeginPrivateObject,
	BeginPrivateObject,
	NullReference,
}

type Trace struct {
	Year                    int16
	Month                   int16
	DayOfWeek               int16
	Day                     int16
	Hour                    int16
	Minute                  int16
	Second                  int16
	Millisecond             int16
	SyncTimeQPC             int64
	QPCFrequency            int64
	PointerSize             int32
	ProcessID               int32
	NumberOfProcessors      int32
	ExpectedCPUSamplingRate int32
}

// unsafe.Sizeof(Trace{})
const traceLen = 48

type netTraceHeader struct {
	NetTraceMagic          [8]byte
	Len                    int32
	FastSerializationMagic [20]byte
}

func (h netTraceHeader) validate() error {
	if !(h.NetTraceMagic == netTraceMagic &&
		h.FastSerializationMagic == fastSerializationMagic) {
		return ErrInvalidNetTraceHeader
	}
	return nil
}

type objectHeader struct {
	Tags                 [3]Tag
	Version              int32
	MinimumReaderVersion int32
	TypeNameLen          int32
}

func (h objectHeader) validate() error {
	if h.Tags != objectHeaderTags {
		return fmt.Errorf("invalid object header: %w", ErrUnexpectedTag)
	}
	// Version check. Not sure if it should be done after we get the type name:
	// in this case we would have a bit more meaningful error message.
	if Version < h.MinimumReaderVersion {
		return fmt.Errorf("%w: %#x", ErrUnsupportedObjectVersion, h.Version)
	}
	return nil
}

// unsafe.Sizeof(objectHeader{}) aligns to 16.
const objectHeaderSize = 15

func (d *Decoder) OpenTrace() (*Trace, error) {
	var header netTraceHeader
	var err error
	if err = d.read(&header); err != nil {
		return nil, err
	}
	if err = header.validate(); err != nil {
		return nil, ErrInvalidNetTraceHeader
	}
	var o Object
	if err = d.readObject(&o); err != nil {
		return nil, err
	}
	if o.Type != ObjectTypeTrace {
		return nil, fmt.Errorf("%w: %s", ErrUnexpectedObjectType, o.Type)
	}
	var trace Trace
	if err = d.readFrom(o.Payload, &trace); err != nil {
		return nil, fmt.Errorf("invalid trace object: %w", err)
	}
	return &trace, nil
}

// Decode deserializes next NetTrace Object from stream to o.
// The call returns io.EOF when the stream is properly terminated,
// any further attempts to decode will return io.ErrUnexpectedEOF.
func (d *Decoder) Decode(o *Object) error {
	return d.readObject(o)
}

func (d *Decoder) Offset() uint64 {
	return d.r.offset
}

func (d *Decoder) read(v interface{}) error {
	return d.readFrom(d.r, v)
}

func (d *Decoder) readFrom(r io.Reader, v interface{}) error {
	return binary.Read(r, binary.LittleEndian, v)
}

// readObjectHeader reads NetTrace Object header to the given header.
// The call returns io.EOF, only when the stream is properly terminated.
func (d *Decoder) readObjectHeader(header *objectHeader) error {
	b := make([]byte, objectHeaderSize)
	n, err := d.r.Read(b)
	switch {
	default:
	case n != objectHeaderSize:
		// After the last object is emitted, the stream is ended by emitting
		// a NullReference Tag which indicates that there are no more objects
		// in the stream to read.
		if n == 1 && NullReference == Tag(b[0]) {
			return io.EOF
		}
		return io.ErrUnexpectedEOF
	case errors.Is(err, io.EOF):
		return io.ErrUnexpectedEOF
	case err != nil:
		return fmt.Errorf("reading object header: %w", err)
	}

	if err = d.readFrom(bytes.NewBuffer(b), header); err != nil {
		return fmt.Errorf("unmarshaling object header: %w", err)
	}
	return header.validate()
}

// readObject reads next NetTrace Object from stream to o.
func (d *Decoder) readObject(o *Object) error {
	var header objectHeader
	var err error
	if err = d.readObjectHeader(&header); err != nil {
		return err
	}

	typeName := make([]byte, header.TypeNameLen)
	if _, err = io.ReadFull(d.r, typeName); err != nil {
		return fmt.Errorf("reading type name: %w", err)
	}
	objectType := ObjectType(typeName)
	if !objectType.IsValid() {
		return fmt.Errorf("%w: %s", ErrInvalidObjectType, objectType)
	}

	// Type end.
	if err = d.expectEndTag(); err != nil {
		return fmt.Errorf("reading type object: %w", err)
	}

	o.Payload, err = d.objectPayload(objectType)
	if err != nil {
		return fmt.Errorf("reading object payload: %w", err)
	}

	// Object end.
	if err = d.expectEndTag(); err != nil {
		return fmt.Errorf("completing object read: %w", err)
	}

	o.Version = header.Version
	o.MinimumReaderVersion = header.MinimumReaderVersion
	o.Type = objectType

	return nil
}

func (d *Decoder) expectEndTag() error {
	var endTag Tag
	if err := d.read(&endTag); err != nil {
		return err
	}
	if endTag != EndObject {
		return fmt.Errorf("%w: %#x", ErrUnexpectedTag, endTag)
	}
	return nil
}

func (d *Decoder) objectPayload(t ObjectType) (*bytes.Buffer, error) {
	switch t {
	case ObjectTypeTrace:
		s := make([]byte, traceLen)
		if _, err := io.ReadFull(d.r, s); err != nil {
			return nil, err
		}
		return bytes.NewBuffer(s), nil

	case ObjectTypeMetadataBlock, ObjectTypeEventBlock, ObjectTypeSPBlock, ObjectTypeStackBlock:
		return d.blockPayload()

	default:
		// Should never happen as we perform the type check in advance.
		return nil, fmt.Errorf("%w: %s", ErrInvalidObjectType, t)
	}
}

func (d *Decoder) blockPayload() (*bytes.Buffer, error) {
	// Block Layout:
	//   1. BlockSize int32 - Size of the block in bytes starting after the alignment padding
	//   2. 0 padding to reach 4 byte alignment.
	//   3. Object-specific header (optional?).
	//   4. Object payload (optional?).
	var size int32
	if err := d.read(&size); err != nil {
		return nil, err
	}
	// Ensure 4-byte alignment.
	padLen := (d.r.offset) % 4
	if padLen != 0 {
		if _, err := io.ReadFull(d.r, make([]byte, 4-padLen)); err != nil {
			return nil, err
		}
	}
	blockData := make([]byte, int(size))
	if _, err := io.ReadFull(d.r, blockData); err != nil {
		return nil, err
	}
	return bytes.NewBuffer(blockData), nil
}
