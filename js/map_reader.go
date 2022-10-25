package js

import (
	"encoding/json"
	"io"
	"os"

	"github.com/aserto-dev/go-directory/pkg/pb"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

type Reader interface {
	Read(proto.Message) error
	Close() error
	GetObjectType() string
}

type MapReader struct {
	dec     *json.Decoder
	first   bool
	rootKey string
}

func NewReader(file string) (Reader, error) {
	r, err := os.Open(file)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file: [%s]", file)
	}
	dec := json.NewDecoder(r)

	// advance reader to start token
	tok, err := dec.Token()
	if err != nil {
		r.Close()
		return nil, err
	}

	keyStr := ""
	if del, ok := tok.(json.Delim); ok {
		// get key value if not array
		if del == '{' {
			t, err := dec.Token()
			if err != nil {
				r.Close()
				return nil, err
			}

			if key, ok := t.(string); ok {
				keyStr = key
			}

			tok, err := dec.Token()
			if err != nil {
				r.Close()
				return nil, err
			}
			if delim, ok := tok.(json.Delim); !ok && delim.String() != "[" {
				r.Close()
				return nil, errors.Errorf("file does not contain a JSON array")
			}

			return &MapReader{
				dec:     dec,
				first:   false,
				rootKey: keyStr,
			}, nil
		}
	}

	r.Close()
	return nil, nil
}

func (r *MapReader) GetObjectType() string {
	return r.rootKey
}

func (r *MapReader) Close() error {
	return nil
}

func (r *MapReader) Read(m proto.Message) error {
	if !r.dec.More() {
		return io.EOF
	}

	if err := pb.UnmarshalNext(r.dec, m); err != nil {
		return err
	}
	return nil
}
