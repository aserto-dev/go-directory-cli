package client

import (
	"context"
	"os"
	"sort"

	v2 "github.com/aserto-dev/go-directory/aserto/directory/common/v2"
	"github.com/aserto-dev/go-directory/aserto/directory/writer/v2"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Relation map[string][]string

type ObjectRelation map[string]Relation

type Manifest map[string]ObjectRelation

func (c *Client) Load(ctx context.Context, file string) error {

	c.UI.Normal().Msgf("Reading manifest [%s]", file)

	yfile, err := os.ReadFile(file)

	if err != nil {
		return err
	}

	manifestData := make(Manifest, 0)

	err = yaml.Unmarshal(yfile, &manifestData)

	if err != nil {
		return err
	}

	manifestEntriesWithUnions := make(ObjectRelation, 0)
	permissions := make(map[string]bool, 0)

	c.UI.Normal().Msg("Loading from manifest")
	for objectType, manifestEntry := range manifestData {
		req := &writer.SetObjectTypeRequest{
			ObjectType: &v2.ObjectType{
				Name: objectType,
			}}
		_, err := c.Writer.SetObjectType(ctx, req)
		if err != nil {
			return errors.Wrapf(err, "failed to set object type %s", objectType)
		}
		relationsLoaded := []string{}
		for relationType, data := range manifestEntry {
			// at first we create relation types that don't have unions
			if len(data["union"]) > 0 {
				manifestEntriesWithUnions[relationType] = data
			} else {
				err := c.setPermissions(data["permissions"], permissions)
				if err != nil {
					return errors.Wrapf(err, "failed to set permissions for relation %s", relationType)
				}

				req := &writer.SetRelationTypeRequest{
					RelationType: &v2.RelationType{
						Name:        relationType,
						Permissions: data["permissions"],
						ObjectType:  objectType,
					}}
				_, err = c.Writer.SetRelationType(ctx, req)
				if err != nil {
					return errors.Wrapf(err, "failed to set relation type %s", relationType)
				}
				relationsLoaded = append(relationsLoaded, relationType)
			}
		}
		for len(manifestEntriesWithUnions) > 0 {
			manifestEntriesWithUnions, err = c.loadRels(ctx, manifestEntriesWithUnions, permissions, relationsLoaded, objectType)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Client) loadRels(ctx context.Context, manifestEntriesWithUnions ObjectRelation, permissions map[string]bool, relationsLoaded []string, objectType string) (ObjectRelation, error) {
	for relationType, data := range manifestEntriesWithUnions {
		err := c.setPermissions(data["permissions"], permissions)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to set permissions for relation %s", relationType)
		}

		if !allUnionRelWereLoaded(data["union"], relationsLoaded) {
			continue
		}

		req := &writer.SetRelationTypeRequest{
			RelationType: &v2.RelationType{
				Name:        relationType,
				Permissions: data["permissions"],
				ObjectType:  objectType,
				Unions:      data["union"],
			}}
		_, err = c.Writer.SetRelationType(ctx, req)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to set relation type %s", relationType)
		}
		relationsLoaded = append(relationsLoaded, relationType)
		delete(manifestEntriesWithUnions, relationType)
	}

	return manifestEntriesWithUnions, nil
}

func (c *Client) setPermissions(permissions []string, alreadyAddedPerms map[string]bool) error {
	for _, perm := range permissions {
		if !alreadyAddedPerms[perm] {
			req := &writer.SetPermissionRequest{Permission: &v2.Permission{Name: perm}}
			_, err := c.Writer.SetPermission(context.Background(), req)
			if err != nil {
				return err
			}
			alreadyAddedPerms[perm] = true
		}
	}
	return nil
}

func allUnionRelWereLoaded(unions, relationsLoaded []string) bool {
	if len(unions) > len(relationsLoaded) {
		return false
	}
	sort.Strings(unions)
	sort.Strings(relationsLoaded)

	for _, unionRel := range unions {
		var found bool
		for _, rel := range relationsLoaded {
			if unionRel == rel {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}
	return true
}
