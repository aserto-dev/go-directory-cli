package client

import (
	"context"
	"io"

	"github.com/aserto-dev/go-directory-cli/counter"
	"github.com/aserto-dev/go-directory-cli/js"
	dse "github.com/aserto-dev/go-directory/aserto/directory/exporter/v2"
	"github.com/fatih/color"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (c *Client) Export(ctx context.Context, objectsFile, relationsFile string) error {
	stream, err := c.Exporter.Export(ctx, &dse.ExportRequest{
		Options:   uint32(dse.Option_OPTION_DATA),
		StartFrom: &timestamppb.Timestamp{},
	})
	if err != nil {
		return err
	}

	color.Green(">>> exporting objects to %s", objectsFile)
	objects, err := js.NewWriter(objectsFile, ObjectsStr)
	if err != nil {
		return err
	}
	defer objects.Close()

	color.Green(">>> exporting relations to %s", relationsFile)
	relations, err := js.NewWriter(relationsFile, RelationsStr)
	if err != nil {
		return err
	}
	defer relations.Close()

	ctr := counter.New()

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		switch m := msg.Msg.(type) {
		case *dse.ExportResponse_Object:
			err = objects.Write(m.Object)
			ctr.Objects.Incr().Print(c.UI.Output())

		case *dse.ExportResponse_Relation:
			err = relations.Write(m.Relation)
			ctr.Relations.Incr().Print(c.UI.Output())

		default:
			c.UI.Problem().Msg("unknown message type")
		}

		if err != nil {
			c.UI.Problem().Msgf("err: %v", err)
		}
	}

	ctr.Print(c.UI.Output())
	color.Green(">>> finished export")

	return nil
}
