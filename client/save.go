package client

import (
	"context"
	"fmt"
	"os"

	v2 "github.com/aserto-dev/go-directory/aserto/directory/common/v2"
	"github.com/aserto-dev/go-directory/aserto/directory/reader/v2"
	"gopkg.in/yaml.v2"
)

func (c *Client) Save(ctx context.Context, file string) error {

	manifestData := make(Manifest, 0)

	// read object types
	pageToken := ""
	for {
		req := &reader.GetObjectTypesRequest{
			Page: &v2.PaginationRequest{Token: pageToken}}
		resp, err := c.Reader.GetObjectTypes(ctx, req)
		if err != nil {
			return err
		}

		for _, objType := range resp.Results {
			err := c.getRelationTypes(ctx, manifestData, objType)
			if err != nil {
				return err
			}
		}

		pageToken = resp.Page.NextToken
		if pageToken == "" {
			break
		}
	}

	fmt.Fprintf(c.UI.Output(), ">>> writing manifest to file [%s]", file)
	yamlData, err := yaml.Marshal(&manifestData)
	if err != nil {
		return err
	}

	err = os.WriteFile(file, yamlData, 0600)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) getRelationTypes(ctx context.Context, data Manifest, object *v2.ObjectType) error {
	token := ""
	objRel := make(ObjectRelation, 0)

	for {
		relReq := &reader.GetRelationTypesRequest{
			Param: &v2.ObjectTypeIdentifier{Id: &object.Id},
			Page:  &v2.PaginationRequest{Token: token},
		}
		resp, err := c.Reader.GetRelationTypes(ctx, relReq)
		if err != nil {
			return err
		}

		for _, relationType := range resp.Results {
			rels := make(Relation, 0)
			rels["union"] = relationType.Unions
			rels["permissions"] = relationType.Permissions
			objRel[relationType.Name] = rels
		}

		token = resp.Page.NextToken
		if token == "" {
			break
		}
	}

	data[object.Name] = objRel

	return nil
}
