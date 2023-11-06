package client

import (
	"github.com/aserto-dev/clui"
	v2 "github.com/aserto-dev/go-directory-cli/client/v2"
	v3 "github.com/aserto-dev/go-directory-cli/client/v3"

	"google.golang.org/grpc"
)

type Client struct {
	conn grpc.ClientConnInterface
	V2   *v2.Client
	V3   *v3.Client
	UI   *clui.UI
}

func New(conn grpc.ClientConnInterface, ui *clui.UI) (*Client, error) {
	dsc2, err := v2.New(conn, ui)
	if err != nil {
		return nil, err
	}
	dsc3, err := v3.New(conn, ui)
	if err != nil {
		return nil, err
	}

	c := Client{
		conn: conn,
		V2:   dsc2,
		V3:   dsc3,
		UI:   ui,
	}
	return &c, nil
}
