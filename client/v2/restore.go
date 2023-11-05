package v2

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"io"
	"os"
	"path"
	"strings"

	"github.com/aserto-dev/go-directory-cli/client/x"
	"github.com/aserto-dev/go-directory-cli/counter"
	"github.com/aserto-dev/go-directory-cli/js"
	dsc2 "github.com/aserto-dev/go-directory/aserto/directory/common/v2"
	dsi2 "github.com/aserto-dev/go-directory/aserto/directory/importer/v2"
	"golang.org/x/sync/errgroup"
)

func (c *Client) Restore(ctx context.Context, file string) error {
	tf, err := os.Open(file)
	if err != nil {
		return err
	}
	defer tf.Close()

	gz, err := gzip.NewReader(tf)
	if err != nil {
		return err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)

	ctr := counter.New()
	defer ctr.Print(c.UI.Output())

	g, iCtx := errgroup.WithContext(context.Background())
	stream, err := c.Importer.Import(iCtx)
	if err != nil {
		return err
	}

	g.Go(c.receiver(stream))

	g.Go(c.restoreHandler(stream, tr, ctr))

	return g.Wait()
}

func (c *Client) receiver(stream dsi2.Importer_ImportClient) func() error {
	return func() error {
		for {
			_, err := stream.Recv()
			if err == io.EOF {
				return nil
			}

			if err != nil {
				return err
			}
		}
	}
}

func (c *Client) restoreHandler(stream dsi2.Importer_ImportClient, tr *tar.Reader, ctr *counter.Counter) func() error {
	objectsCounter := ctr.Objects()
	relationsCounter := ctr.Relations()

	return func() error {
		for {
			header, err := tr.Next()
			if errors.Is(err, io.EOF) {
				break
			}

			if err != nil {
				return err
			}

			if header == nil || header.Typeflag != tar.TypeReg {
				continue
			}

			r, err := js.NewReader(tr, c.UI)
			if err != nil {
				return err
			}

			name := path.Clean(header.Name)
			switch name {
			case x.ObjectsFileName:
				if err := c.loadObjects(stream, r, objectsCounter); err != nil {
					return err
				}

			case x.RelationsFileName:
				if err := c.loadRelations(stream, r, relationsCounter); err != nil {
					return err
				}
			}
		}

		if err := stream.CloseSend(); err != nil {
			return err
		}

		return nil
	}
}

func (c *Client) loadObjects(stream dsi2.Importer_ImportClient, objects *js.Reader, ctr *counter.Item) error {
	defer objects.Close()

	var m dsc2.Object

	for {
		err := objects.Read(&m)
		if err == io.EOF {
			break
		}

		if err != nil {
			if strings.Contains(err.Error(), "unknown field") {
				ctr.Skip()
				continue
			}
			return err
		}

		if err := stream.Send(&dsi2.ImportRequest{
			OpCode: dsi2.Opcode_OPCODE_SET,
			Msg: &dsi2.ImportRequest_Object{
				Object: &m,
			},
		}); err != nil {
			return err
		}
		ctr.Incr().Print(c.UI.Output())
	}

	return nil
}

func (c *Client) loadRelations(stream dsi2.Importer_ImportClient, relations *js.Reader, ctr *counter.Item) error {
	defer relations.Close()

	var m dsc2.Relation

	for {
		err := relations.Read(&m)
		if err == io.EOF {
			break
		}
		if err != nil {
			if strings.Contains(err.Error(), "unknown field") {
				ctr.Skip()
				continue
			}
			return err
		}

		if err := stream.Send(&dsi2.ImportRequest{
			OpCode: dsi2.Opcode_OPCODE_SET,
			Msg: &dsi2.ImportRequest_Relation{
				Relation: &m,
			},
		}); err != nil {
			return err
		}

		ctr.Incr().Print(c.UI.Output())
	}

	return nil
}
