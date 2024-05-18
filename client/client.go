package client

import (
	"github.com/aserto-dev/clui"
	v3 "github.com/aserto-dev/go-directory-cli/client/v3"

	"google.golang.org/grpc"
)

type Client struct {
	conn grpc.ClientConnInterface
	V3   *v3.Client
	UI   *clui.UI
}

func New(conn grpc.ClientConnInterface, ui *clui.UI) (*Client, error) {
	dsc3, err := v3.New(conn, ui)
	if err != nil {
		return nil, err
	}

	c := Client{
		conn: conn,
		V3:   dsc3,
		UI:   ui,
	}
	return &c, nil
}
