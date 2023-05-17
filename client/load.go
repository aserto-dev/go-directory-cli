package client

import (
	"context"
	"fmt"
	"os"

	v2 "github.com/aserto-dev/go-directory/aserto/directory/common/v2"
	"github.com/aserto-dev/go-directory/aserto/directory/reader/v2"
	"github.com/aserto-dev/go-directory/aserto/directory/writer/v2"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

const (
	union string = "union"
	perms string = "permissions"
)

type Relation map[string][]string

type ObjectRelation map[string]Relation

type Manifest map[string]ObjectRelation

func (c *Client) Load(ctx context.Context, file string) error {
	buf, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	manifest := make(Manifest, 0)

	if err := yaml.Unmarshal(buf, &manifest); err != nil {
		return err
	}

	if err := c.createObjectTypes(ctx, manifest); err != nil {
		return err
	}

	if err := c.createPermissions(ctx, manifest); err != nil {
		return err
	}

	if err := c.createRelationTypes(ctx, manifest); err != nil {
		return err
	}

	if err := c.updateRelationTypes(ctx, manifest); err != nil {
		return err
	}

	return nil
}

// create object types.
func (c *Client) createObjectTypes(ctx context.Context, manifest Manifest) error {
	for objectType := range manifest {
		fmt.Printf("set-object %s\n", objectType)
		if err := c.getSetObjectType(ctx, objectType); err != nil {
			return errors.Wrapf(err, "failed to set object type %s", objectType)
		}
	}
	return nil
}

// create permissions.
func (c *Client) createPermissions(ctx context.Context, manifest Manifest) error {
	permissions := map[string]bool{}
	for _, objectRelation := range manifest {
		for _, v := range objectRelation {
			for _, permission := range v[perms] {
				if _, ok := permissions[permission]; !ok {
					fmt.Printf("set-permission %s\n", permission)
					if err := c.getSetPermission(ctx, permission); err != nil {
						return errors.Wrapf(err, "failed to set permission %s", permission)
					}
					permissions[permission] = true
				}
			}
		}
	}
	return nil
}

// create relation types, without unions or permissions.
func (c *Client) createRelationTypes(ctx context.Context, manifest Manifest) error {
	for objectType, objectRelation := range manifest {
		for relationType := range objectRelation {
			fmt.Printf("set-relation %s#%s\n", objectType, relationType)
			if err := c.getSetRelationType(ctx, objectType, relationType, []string{}, []string{}); err != nil {
				return errors.Wrapf(err, "failed to set relation type %s#%s", objectType, relationType)
			}
		}
	}
	return nil
}

// update relation types with unions and permission
func (c *Client) updateRelationTypes(ctx context.Context, manifest Manifest) error {
	for objectType, objectRelation := range manifest {
		for relationType, v := range objectRelation {
			if len(v[union]) == 0 {
				continue
			}

			fmt.Printf("update-relation %s#%s\n", objectType, relationType)
			if err := c.getSetRelationType(ctx, objectType, relationType, v[union], []string{}); err != nil {
				return errors.Wrapf(err, "failed to set relation type %s#%s", objectType, relationType)
			}
		}
	}

	for objectType, objectRelation := range manifest {
		for relationType, v := range objectRelation {
			if len(v[perms]) == 0 {
				continue
			}

			fmt.Printf("update-relation %s#%s\n", objectType, relationType)
			if err := c.getSetRelationType(ctx, objectType, relationType, v[union], v[perms]); err != nil {
				return errors.Wrapf(err, "failed to set relation type %s#%s", objectType, relationType)
			}
		}
	}

	return nil
}

// get  object type.
func (c *Client) getObjectType(ctx context.Context, objectType string) (*v2.ObjectType, error) {
	resp, err := c.Reader.GetObjectType(ctx, &reader.GetObjectTypeRequest{
		Param: &v2.ObjectTypeIdentifier{
			Name: proto.String(objectType),
		}})

	if err != nil {
		return &v2.ObjectType{}, err
	}

	return resp.Result, nil
}

// get and set object type.
func (c *Client) getSetObjectType(ctx context.Context, objectType string) error {
	req := &writer.SetObjectTypeRequest{
		ObjectType: &v2.ObjectType{
			Name:        objectType,
			DisplayName: objectType,
		}}

	resp, err := c.getObjectType(ctx, objectType)
	switch {
	case status.Code(err) == codes.NotFound:
		req.ObjectType.Hash = ""
	case err != nil:
		return err
	default:
		req.ObjectType.Hash = resp.Hash
	}

	_, err = c.Writer.SetObjectType(ctx, req)

	return err
}

// get relation type.
func (c *Client) getRelationType(ctx context.Context, objectType, relationType string) (*v2.RelationType, error) {
	resp, err := c.Reader.GetRelationType(ctx, &reader.GetRelationTypeRequest{
		Param: &v2.RelationTypeIdentifier{
			ObjectType: proto.String(objectType),
			Name:       proto.String(relationType),
		}})

	if err != nil {
		return &v2.RelationType{}, err
	}

	return resp.Result, nil
}

// get set relation type.
func (c *Client) getSetRelationType(ctx context.Context, objectType, relationType string, unions, permissions []string) error {
	req := &writer.SetRelationTypeRequest{
		RelationType: &v2.RelationType{
			ObjectType:  objectType,
			Name:        relationType,
			DisplayName: objectType + ":" + relationType,
		}}

	if len(unions) > 0 {
		req.RelationType.Unions = append(req.RelationType.Unions, unions...)
	}

	if len(permissions) > 0 {
		req.RelationType.Permissions = append(req.RelationType.Permissions, permissions...)
	}

	resp, err := c.getRelationType(ctx, objectType, relationType)
	switch {
	case status.Code(err) == codes.NotFound:
		resp.Hash = ""
	case err != nil:
		return err
	default:
		req.RelationType.Hash = resp.Hash
	}

	_, err = c.Writer.SetRelationType(ctx, req)

	return err
}

// get permission.
func (c *Client) getPermission(ctx context.Context, permission string) (*v2.Permission, error) {
	resp, err := c.Reader.GetPermission(ctx, &reader.GetPermissionRequest{
		Param: &v2.PermissionIdentifier{
			Name: proto.String(permission),
		}})

	if err != nil {
		return &v2.Permission{}, err
	}
	return resp.Result, nil
}

// get set permission.
func (c *Client) getSetPermission(ctx context.Context, permission string) error {
	req := &writer.SetPermissionRequest{
		Permission: &v2.Permission{
			Name:        permission,
			DisplayName: permission,
		}}

	resp, err := c.getPermission(ctx, permission)
	switch {
	case status.Code(err) == codes.NotFound:
		req.Permission.Hash = ""
	case err != nil:
		return err
	default:
		req.Permission.Hash = resp.Hash
	}

	_, err = c.Writer.SetPermission(ctx, req)

	return err
}
