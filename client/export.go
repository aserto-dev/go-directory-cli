package client

import (
	"context"
	"io"

	"github.com/aserto-dev/go-directory-cli/counter"
	"github.com/aserto-dev/go-directory-cli/js"
	dse "github.com/aserto-dev/go-directory/aserto/directory/exporter/v2"
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

	c.UI.Normal().Msgf("Exporting objects to %s", objectsFile)
	objects, _ := js.NewArrayWriter(objectsFile)
	defer objects.Close()

	c.UI.Normal().Msgf("Exporting relations to %s", relationsFile)
	relations, _ := js.NewArrayWriter(relationsFile)
	defer relations.Close()

	counter := counter.New()
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
			counter.Objects.Incr().Print(c.UI.Output())

		case *dse.ExportResponse_Relation:
			err = relations.Write(m.Relation)
			counter.Relations.Incr().Print(c.UI.Output())

		default:
			c.UI.Problem().Msg("unknown message type")
		}

		if err != nil {
			c.UI.Problem().Msgf("err: %v", err)
		}
	}

	c.UI.Normal().Msg("Finished export.")
	counter.Print(c.UI.Output())

	return nil
}
