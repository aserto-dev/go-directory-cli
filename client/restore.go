package client

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"io"
	"os"
	"path"

	"github.com/aserto-dev/go-directory-cli/counter"
	"github.com/aserto-dev/go-directory-cli/js"
	dsc "github.com/aserto-dev/go-directory/aserto/directory/common/v2"
	dsw "github.com/aserto-dev/go-directory/aserto/directory/writer/v2"
)

// nolint: gocyclo // to be refactored
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

	var stop bool
	for {
		header, err := tr.Next()

		switch {
		case err == io.EOF:
			return nil
		case err != nil:
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
		case ObjectTypesFileName:
			if err := c.loadObjectTypes(ctx, r, ctr.ObjectTypes); err != nil {
				return err
			}

		case PermissionsFileName:
			if err := c.loadPermissions(ctx, r, ctr.Permissions); err != nil {
				return err
			}

		case RelationTypesFileName:
			if err := c.loadRelationTypes(ctx, r, ctr.RelationTypes); err != nil {
				return err
			}

		case ObjectsFileName:
			if err := c.loadObjects(ctx, r, ctr.Objects); err != nil {
				return err
			}

		case RelationsFileName:
			if err := c.loadRelations(ctx, r, ctr.Relations); err != nil {
				return err
			}

		default:
			stop = true
		}

		if stop {
			break
		}
	}

	return nil
}

func (c *Client) loadObjectTypes(ctx context.Context, objTypes *js.Reader, ctr *counter.Item) error {
	defer objTypes.Close()

	var m dsc.ObjectType

	for {
		err := objTypes.Read(&m)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		_, err = c.Writer.SetObjectType(ctx, &dsw.SetObjectTypeRequest{
			ObjectType: &m,
		})
		if err != nil {
			return err
		}
		ctr.Incr().Print(c.UI.Output())
	}

	return nil
}

func (c *Client) loadPermissions(ctx context.Context, permissions *js.Reader, ctr *counter.Item) error {
	defer permissions.Close()

	var m dsc.Permission

	for {
		err := permissions.Read(&m)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		_, err = c.Writer.SetPermission(ctx, &dsw.SetPermissionRequest{
			Permission: &m,
		})
		if err != nil {
			return err
		}
		ctr.Incr().Print(c.UI.Output())
	}

	return nil
}

func (c *Client) loadRelationTypes(ctx context.Context, relTypes *js.Reader, ctr *counter.Item) error {
	defer relTypes.Close()

	var m dsc.RelationType

	for {
		err := relTypes.Read(&m)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		_, err = c.Writer.SetRelationType(ctx, &dsw.SetRelationTypeRequest{
			RelationType: &m,
		})
		if err != nil {
			return err
		}
		ctr.Incr().Print(c.UI.Output())
	}

	return nil
}

func (c *Client) loadObjects(ctx context.Context, objects *js.Reader, ctr *counter.Item) error {
	defer objects.Close()

	var m dsc.Object

	for {
		err := objects.Read(&m)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		_, err = c.Writer.SetObject(ctx, &dsw.SetObjectRequest{
			Object: &m,
		})
		if err != nil {
			return err
		}
		ctr.Incr().Print(c.UI.Output())
	}

	return nil
}

func (c *Client) loadRelations(ctx context.Context, relations *js.Reader, ctr *counter.Item) error {
	defer relations.Close()

	var m dsc.Relation

	for {
		err := relations.Read(&m)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		_, err = c.Writer.SetRelation(ctx, &dsw.SetRelationRequest{
			Relation: &m,
		})
		if err != nil {
			return err
		}
		ctr.Incr().Print(c.UI.Output())
	}

	return nil
}
