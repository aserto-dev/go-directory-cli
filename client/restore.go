package client

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"io"
	"os"
	"path"

	"github.com/aserto-dev/go-directory-cli/js"
	dsc "github.com/aserto-dev/go-directory/aserto/directory/common/v2"
	dsw "github.com/aserto-dev/go-directory/aserto/directory/writer/v2"
)

func (c *Client) Restore(ctx context.Context, file string) error {

	c.UI.Normal().Msgf("Reading file %s", file)
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

		name := path.Clean(header.Name)
		switch name {
		case "object_types.json":
			if err := c.loadObjectTypes(ctx, tr); err != nil {
				return err
			}

		case "permissions.json":
			if err := c.loadPermissions(ctx, tr); err != nil {
				return err
			}

		case "relation_types.json":
			if err := c.loadRelationTypes(ctx, tr); err != nil {
				return err
			}

		case "objects.json":
			if err := c.loadObjects(ctx, tr); err != nil {
				return err
			}

		case "relations.json":
			if err := c.loadRelations(ctx, tr); err != nil {
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

func (c *Client) loadObjectTypes(ctx context.Context, r io.Reader) error {
	objTypes, _ := js.NewArrayReader(r)
	defer objTypes.Close()

	c.UI.Normal().Msg("Loading object types")
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

	}
	return nil
}

func (c *Client) loadPermissions(ctx context.Context, r io.Reader) error {
	permissions, _ := js.NewArrayReader(r)
	defer permissions.Close()

	c.UI.Normal().Msg("Loading pemissions")
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

	}
	return nil
}

func (c *Client) loadRelationTypes(ctx context.Context, r io.Reader) error {
	relTypes, _ := js.NewArrayReader(r)
	defer relTypes.Close()

	c.UI.Normal().Msg("Loading relation types")
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
	}
	return nil
}

func (c *Client) loadObjects(ctx context.Context, r io.Reader) error {
	objects, _ := js.NewArrayReader(r)
	defer objects.Close()

	c.UI.Normal().Msg("Loading objects")
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

	}
	return nil
}

func (c *Client) loadRelations(ctx context.Context, r io.Reader) error {
	relations, _ := js.NewArrayReader(r)
	defer relations.Close()

	c.UI.Normal().Msg("Loading relations")
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

		if err != nil {
			return err
		}
	}
	return nil
}
