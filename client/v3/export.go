package v3

import (
	"context"
	"fmt"
	"io"

	"github.com/aserto-dev/go-directory-cli/client/x"
	"github.com/aserto-dev/go-directory-cli/counter"
	"github.com/aserto-dev/go-directory-cli/js"
	dse3 "github.com/aserto-dev/go-directory/aserto/directory/exporter/v3"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (c *Client) Export(ctx context.Context, objectsFile, relationsFile string) error {
	stream, err := c.Exporter.Export(ctx, &dse3.ExportRequest{
		Options:   uint32(dse3.Option_OPTION_DATA),
		StartFrom: &timestamppb.Timestamp{},
	})
	if err != nil {
		return err
	}

	objects, err := js.NewWriter(objectsFile, x.ObjectsStr)
	if err != nil {
		return err
	}
	defer objects.Close()

	relations, err := js.NewWriter(relationsFile, x.RelationsStr)
	if err != nil {
		return err
	}
	defer relations.Close()

	ctr := counter.New()
	objectsCounter := ctr.Objects()
	relationsCounter := ctr.Relations()

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		switch m := msg.Msg.(type) {
		case *dse3.ExportResponse_Object:
			err = objects.Write(m.Object)
			objectsCounter.Incr().Print(c.Out())

		case *dse3.ExportResponse_Relation:
			err = relations.Write(m.Relation)
			relationsCounter.Incr().Print(c.Out())

		default:
			fmt.Fprintf(c.Err(), "unknown message type\n")
		}

		if err != nil {
			fmt.Fprintf(c.Err(), "err: %v\n", err)
		}
	}

	ctr.Print(c.Out())

	return nil
}
