package httputil

import (
	"bytes"
	"encoding/base64"
	"errors"
	"strconv"
	"time"
)

var (
	specialSeparator = "--"
	lineEndSeparator = "\r\n"
	headSeparator    = ":"
)

type HttpEntity struct {
	Header map[string]string
	Body   []byte
}

type MultipartBuilder struct {
	boundary string
	entites  []HttpEntity
	buffer   *bytes.Buffer
}

func NewMultipartBuilder() *MultipartBuilder {
	builder := &MultipartBuilder{
		boundary: getDefaultBoundary(),
		buffer:   &bytes.Buffer{},
	}

	return builder
}

func (r *MultipartBuilder) SetBoundary(b string) *MultipartBuilder {
	if len(b) > 0 {
		r.boundary = b
	}

	return r
}

func (r *MultipartBuilder) AppendEntity(entity *HttpEntity) error {
	if entity != nil {
		if len(entity.Header) == 0 {
			return errors.New("empty header is not valid")
		}
		if len(entity.Body) == 0 {
			return errors.New("empty body is not valid")
		}

		writeEntity(r.buffer, r.boundary, entity)
	} else {
		return errors.New("nil entity is not valid")
	}

	return nil
}

func (r *MultipartBuilder) GetRequestData() (*bytes.Buffer, error) {
	if r.buffer.Len() == 0 {
		return nil, errors.New("empty data is not invalid")
	}

	writeEnd(r.buffer, r.boundary)

	return r.buffer, nil
}

func getDefaultBoundary() string {
	return base64.RawStdEncoding.EncodeToString([]byte(strconv.Itoa(time.Now().Nanosecond())))
}

func writeEntity(b *bytes.Buffer, boundary string, entity *HttpEntity) {
	b.WriteString(lineEndSeparator)
	b.WriteString(specialSeparator)
	b.WriteString(boundary)
	b.WriteString(lineEndSeparator)
	for k, v := range entity.Header {
		b.WriteString(k)
		b.WriteString(headSeparator)
		b.WriteString(v)
		b.WriteString(lineEndSeparator)
	}
	b.WriteString(lineEndSeparator)
	b.Write(entity.Body)
}

func writeEnd(b *bytes.Buffer, boundary string) {
	b.WriteString(lineEndSeparator)
	b.WriteString(specialSeparator)
	b.WriteString(boundary)
	b.WriteString(specialSeparator)
}
