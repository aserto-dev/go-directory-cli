package client

import (
	"context"
	"io"

	"github.com/aserto-dev/go-directory-cli/counter"
	"github.com/aserto-dev/go-directory-cli/js"
	dse2 "github.com/aserto-dev/go-directory/aserto/directory/exporter/v2"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (c *Client) Export(ctx context.Context, objectsFile, relationsFile string) error {
	stream, err := c.Exporter.Export(ctx, &dse2.ExportRequest{
		Options:   uint32(dse2.Option_OPTION_DATA),
		StartFrom: &timestamppb.Timestamp{},
	})
	if err != nil {
		return err
	}

	objects, err := js.NewWriter(objectsFile, ObjectsStr)
	if err != nil {
		return err
	}
	defer objects.Close()

	relations, err := js.NewWriter(relationsFile, RelationsStr)
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
		case *dse2.ExportResponse_Object:
			err = objects.Write(m.Object)
			objectsCounter.Incr().Print(c.UI.Output())

		case *dse2.ExportResponse_Relation:
			err = relations.Write(m.Relation)
			relationsCounter.Incr().Print(c.UI.Output())

		default:
			c.UI.Problem().Msg("unknown message type")
		}

		if err != nil {
			c.UI.Problem().Msgf("err: %v", err)
		}
	}

	ctr.Print(c.UI.Output())

	return nil
}
