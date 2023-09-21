package client

import (
	"context"
	"os"

	"github.com/aserto-dev/go-directory-cli/counter"
	"github.com/aserto-dev/go-directory-cli/js"
	dsi3 "github.com/aserto-dev/go-directory/aserto/directory/importer/v3"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func (c *Client) Import(ctx context.Context, files []string) error {
	ctr := counter.New()
	defer ctr.Print(c.UI.Output())

	g, iCtx := errgroup.WithContext(context.Background())
	stream, err := c.Importer.Import(iCtx)
	if err != nil {
		return err
	}

	g.Go(c.receiver(stream))

	g.Go(c.importHandler(stream, files, ctr))

	return g.Wait()
}

func (c *Client) importHandler(stream dsi3.Importer_ImportClient, files []string, ctr *counter.Counter) func() error {
	return func() error {
		for _, file := range files {
			if err := c.importFile(stream, file, ctr); err != nil {
				return err
			}
		}

		if err := stream.CloseSend(); err != nil {
			return err
		}

		return nil
	}
}

func (c *Client) importFile(stream dsi3.Importer_ImportClient, file string, ctr *counter.Counter) error {
	r, err := os.Open(file)
	if err != nil {
		return errors.Wrapf(err, "failed to open file: [%s]", file)
	}
	defer r.Close()

	reader, err := js.NewReader(r, c.UI)
	if err != nil || reader == nil {
		c.UI.Problem().Msgf("Skipping file [%s]: [%s]", file, err.Error())
		return nil
	}
	defer reader.Close()

	objectType := reader.GetObjectType()
	switch objectType {
	case ObjectsStr:
		if err := c.loadObjects(stream, reader, ctr.Objects()); err != nil {
			return err
		}

	case RelationsStr:
		if err := c.loadRelations(stream, reader, ctr.Relations()); err != nil {
			return err
		}

	default:
		c.UI.Problem().Msgf("skipping file [%s] with object type [%s]", file, objectType)
	}

	return nil
}
