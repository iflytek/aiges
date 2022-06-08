package nettrace

import (
	"errors"

	"github.com/pyroscope-io/dotnetdiag/nettrace/typecode"
)

var ErrNotImplemented = errors.New("not implemented")

type Metadata struct {
	Header  MetadataHeader
	Payload MetadataPayload
	p       *Parser
}

type MetadataHeader struct {
	MetaDataID   int32
	ProviderName string
	EventID      int32
	EventName    string
	Keywords     int64
	Version      MetadataVersion
	Level        int32
}

type MetadataVersion int32

const (
	_ MetadataVersion = iota

	MetadataLegacyV1 // Used by NetPerf version 1
	MetadataLegacyV2 // Used by NetPerf version 2
	MetadataNetTrace // Used by NetPerf (version 3) and NetTrace (version 4+)
)

type MetadataPayload struct {
	Fields []MetadataField
}

type MetadataField struct {
	TypeCode typecode.TypeCode
	// ArrayTypeCode is an optional field only appears when TypeCode is Array.
	ArrayTypeCode typecode.TypeCode
	// For primitive types and strings Payload is not present, however if TypeCode is Object (1)
	// then Payload is another payload description (that is a field count, followed by a list of
	// field definitions). These can be nested to arbitrary depth.
	Payload MetadataPayload
	Name    string
}

func MetadataFromBlob(blob Blob) (*Metadata, error) {
	md := Metadata{p: &Parser{Buffer: blob.Payload}}
	md.p.Read(&md.Header.MetaDataID)
	md.Header.ProviderName = md.p.UTF16NTS()
	md.p.Read(&md.Header.EventID)
	md.Header.EventName = md.p.UTF16NTS()
	md.p.Read(&md.Header.Keywords)
	md.p.Read(&md.Header.Version)
	md.p.Read(&md.Header.Level)

	if err := md.readPayload(&md.Payload); err != nil {
		return nil, err
	}

	// Version 5 specifies that following the FieldCount number of fields
	// there are an optional set of metadata tags, and if the metadata event
	// specifies a TagKindV2Params tag, the event must have an empty V1
	// parameter FieldCount and no field definitions. Version 5 is not
	// supported yet, therefore we expect the buffer is read to the end.
	// 	if _, err := md.p.ReadByte(); err != io.EOF {
	// 		return nil, ErrNotImplemented
	// 	}

	return &md, md.p.Err()
}

func (md *Metadata) readPayload(mp *MetadataPayload) error {
	var count int32
	md.p.Read(&count)
	for i := int32(0); i < count; i++ {
		var f MetadataField
		if err := md.readField(&f); err != nil {
			return err
		}
		mp.Fields = append(mp.Fields, f)
	}
	return md.p.Err()
}

func (md *Metadata) readField(f *MetadataField) error {
	md.p.Read(&f.TypeCode)
	switch f.TypeCode {
	default:
		// Built-in types do not have payload.
	case typecode.Array:
		return ErrNotImplemented
	case typecode.Object:
		var p MetadataPayload
		if err := md.readPayload(&p); err != nil {
			return err
		}
		f.Payload = p
	}
	f.Name = md.p.UTF16NTS()
	return md.p.Err()
}
