package nettrace

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unicode/utf16"
)

type Parser struct {
	*bytes.Buffer
	errs []error
}

func (p *Parser) Err() error {
	if len(p.errs) != 0 {
		return fmt.Errorf("parser: %w", p.errs[0])
	}
	return nil
}

func (p *Parser) Skip(n int) {
	if n > 0 {
		p.Next(n)
	}
}

func (p *Parser) Read(v interface{}) {
	if err := binary.Read(p.Buffer, binary.LittleEndian, v); err != nil {
		p.errs = append(p.errs, err)
	}
}

func (p *Parser) Uvarint() uint64 {
	n, err := binary.ReadUvarint(p)
	if err != nil {
		p.errs = append(p.errs, err)
	}
	return n
}

// UTF16NTS returns UTF8 string read from 2-bytes UTF16 null terminated string.
// Bytes are speculatively interpreted as 2-byte ASCII chars, and decoding
// takes place only if a code unit is not representable in ASCII.
func (p *Parser) UTF16NTS() string {
	b := p.Bytes()
	s := make([]byte, 0, 64) // The capacity has been chosen empirically.
	for i := 0; i < len(b)-1; i += 2 {
		if b[i] == 0x0 {
			p.Skip(i + 2)
			break
		}
		if b[i+1] != 0x0 {
			return p.decodeUTF16NTS()
		}
		s = append(s, b[i])
	}
	return string(s)
}

func (p *Parser) decodeUTF16NTS() string {
	b := p.Bytes()
	s := make([]uint16, 0, 64)
	for i := 0; i < len(b)-1; i += 2 {
		if b[i] == 0x0 {
			p.Skip(i + 2)
			break
		}
		s = append(s, binary.LittleEndian.Uint16(b[i:(i+2)]))
	}
	return string(utf16.Decode(s))
}
