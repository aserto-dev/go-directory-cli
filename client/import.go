package client

import (
	"context"
	"encoding/json"
	"os"

	dsc "github.com/aserto-dev/go-directory/aserto/directory/common/v2"
	dsw "github.com/aserto-dev/go-directory/aserto/directory/writer/v2"
	"github.com/pkg/errors"
)

type Importer struct {
	Objects       []*dsc.Object       `json:"objects"`
	Relations     []*dsc.Relation     `json:"relations"`
	Permissions   []*dsc.Permission   `json:"permissions"`
	ObjectTypes   []*dsc.ObjectType   `json:"object_types"`
	RelationTypes []*dsc.RelationType `json:"relation_types"`
}

func (c *Client) Import(ctx context.Context, files []string) error {
	var data []Importer

	// read all files
	for _, file := range files {
		var loader Importer
		c.UI.Normal().Msgf("Reading file %s", file)
		b, err := os.ReadFile(file)
		if err != nil {
			return errors.Wrapf(err, "failed to read file: [%s]", file)
		}
		if err := json.Unmarshal(b, &loader); err != nil {
			return errors.Wrapf(err, "failed unmarshal file: [%s]", file)
		}

		data = append(data, loader)
	}

	// import all object types
	c.UI.Normal().Msg("Importing object types...")
	for _, d := range data {
		for _, ot := range d.ObjectTypes {
			resp, err := c.Writer.SetObjectType(ctx, &dsw.SetObjectTypeRequest{ObjectType: ot})
			if err != nil {
				return err
			}
			c.UI.Normal().Msgf("Imported %s", resp.Result.Name)
		}
	}

	// import all permissions
	c.UI.Normal().Msg("Importing permissions...")
	for _, d := range data {
		for _, p := range d.Permissions {
			resp, err := c.Writer.SetPermission(ctx, &dsw.SetPermissionRequest{Permission: p})
			if err != nil {
				return err
			}
			c.UI.Normal().Msgf("Imported %s:%s",
				resp.Result.Id,
				resp.Result.Name,
			)
		}
	}

	// import all relation types
	c.UI.Normal().Msg("Importing relation types...")
	for _, d := range data {
		for _, p := range d.RelationTypes {
			resp, err := c.Writer.SetRelationType(ctx, &dsw.SetRelationTypeRequest{RelationType: p})
			if err != nil {
				return err
			}
			c.UI.Normal().Msgf("Imported %s", resp.Result.Name)
		}
	}

	// import all objects
	c.UI.Normal().Msg("Importing objects...")
	for _, d := range data {
		for _, object := range d.Objects {
			object.Id = ""
			resp, err := c.Writer.SetObject(ctx, &dsw.SetObjectRequest{Object: object})
			if err != nil {
				return err
			}
			c.UI.Normal().Msgf("Imported %s:%s", resp.Result.Type, resp.Result.Id)
		}
	}

	// import all relations
	c.UI.Normal().Msg("Importing relations...")
	for _, d := range data {
		for _, relation := range d.Relations {
			resp, err := c.Writer.SetRelation(ctx, &dsw.SetRelationRequest{Relation: relation})
			if err != nil {
				return err
			}
			c.UI.Normal().Msgf("Imported %s:%s|%s:%s|%s|%s",
				resp.Result.Subject.GetType(),
				resp.Result.Subject.GetId(),
				resp.Result.Object.GetType(),
				resp.Result.Relation,
				resp.Result.Object.GetType(),
				resp.Result.Object.GetId(),
			)
		}
	}

	return nil
}
