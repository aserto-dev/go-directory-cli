package client

import (
	"io"

	v3 "github.com/aserto-dev/go-directory-cli/client/v3"

	"google.golang.org/grpc"
)

type Client struct {
	conn grpc.ClientConnInterface
	V3   *v3.Client
}

func New(conn grpc.ClientConnInterface, stdout, stderr io.Writer) (*Client, error) {
	dsc3, err := v3.New(conn, stdout, stderr)
	if err != nil {
		return nil, err
	}

	c := Client{
		conn: conn,
		V3:   dsc3,
	}
	return &c, nil
}
