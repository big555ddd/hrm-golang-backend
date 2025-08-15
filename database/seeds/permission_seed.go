package seeds

import (
	"app/app/model"
	"context"

	"github.com/uptrace/bun"
)

func permissionSeed(db *bun.DB) error {
	permissions := []model.Permission{
		{
			Name:        "user:create",
			Description: "Permission to create a user",
		},
		{
			Name:        "user:read",
			Description: "Permission to read user data",
		},
		{
			Name:        "user:update",
			Description: "Permission to update user data",
		},
		{
			Name:        "user:delete",
			Description: "Permission to delete a user",
		},
		{
			Name:        "role:create",
			Description: "Permission to create a role",
		},
		{
			Name:        "role:read",
			Description: "Permission to read role data",
		},
		{
			Name:        "role:update",
			Description: "Permission to update role data",
		},
		{
			Name:        "role:delete",
			Description: "Permission to delete a role",
		},
		{
			Name:        "permission:create",
			Description: "Permission to create a permission",
		},
		{
			Name:        "permission:read",
			Description: "Permission to read permission data",
		},
		{
			Name:        "permission:update",
			Description: "Permission to update permission data",
		},
		{
			Name:        "permission:delete",
			Description: "Permission to delete a permission",
		},
		{
			Name:        "department:create",
			Description: "Permission to create a department",
		},
		{
			Name:        "department:read",
			Description: "Permission to read department data",
		},
		{
			Name:        "department:update",
			Description: "Permission to update department data",
		},
		{
			Name:        "department:delete",
			Description: "Permission to delete a department",
		},
		{
			Name:        "branch:create",
			Description: "Permission to create a branch",
		},
		{
			Name:        "branch:read",
			Description: "Permission to read branch data",
		},
		{
			Name:        "branch:update",
			Description: "Permission to update branch data",
		},
		{
			Name:        "branch:delete",
			Description: "Permission to delete a branch",
		},
		{
			Name:        "organization:create",
			Description: "Permission to create an organization",
		},
		{
			Name:        "organization:read",
			Description: "Permission to read organization data",
		},
		{
			Name:        "organization:update",
			Description: "Permission to update organization data",
		},
		{
			Name:        "organization:delete",
			Description: "Permission to delete an organization",
		},
		{
			Name:        "holiday:create",
			Description: "Permission to create a holiday",
		},
		{
			Name:        "holiday:read",
			Description: "Permission to read holiday data",
		},
		{
			Name:        "holiday:update",
			Description: "Permission to update holiday data",
		},
		{
			Name:        "holiday:delete",
			Description: "Permission to delete a holiday",
		},
		{
			Name:        "workshift:create",
			Description: "Permission to create a workshift",
		},
		{
			Name:        "workshift:read",
			Description: "Permission to read workshift data",
		},
		{
			Name:        "workshift:update",
			Description: "Permission to update workshift data",
		},
		{
			Name:        "workshift:delete",
			Description: "Permission to delete a workshift",
		},
		{
			Name:        "leave:create",
			Description: "Permission to create a leave",
		},
		{
			Name:        "leave:read",
			Description: "Permission to read leave data",
		},
		{
			Name:        "leave:update",
			Description: "Permission to update leave data",
		},
		{
			Name:        "leave:delete",
			Description: "Permission to delete a leave",
		},
		{
			Name:        "document:create",
			Description: "Permission to create a document",
		},
		{
			Name:        "document:read",
			Description: "Permission to read document data",
		},
		{
			Name:        "document:update",
			Description: "Permission to update document data",
		},
		{
			Name:        "document:delete",
			Description: "Permission to delete a document",
		},
		{
			Name:        "attendance:create",
			Description: "Permission to create an attendance",
		},
		{
			Name:        "attendance:read",
			Description: "Permission to read attendance data",
		},
		{
			Name:        "attendance:update",
			Description: "Permission to update attendance data",
		},
		{
			Name:        "attendance:delete",
			Description: "Permission to delete an attendance",
		},
	}

	_, err := db.NewInsert().Model(&permissions).Exec(context.Background())
	if err != nil {
		return err
	}
	//set for Admin
	role := &model.Role{}
	err = db.NewSelect().Model(role).Where("name = ?", "Admin").Scan(context.Background())
	if err != nil {
		return err
	}
	rolePermissions := make([]*model.RolePermission, len(permissions))
	for i, permission := range permissions {
		rolePermissions[i] = &model.RolePermission{
			RoleID:       role.ID,
			PermissionID: permission.ID,
		}
	}
	_, err = db.NewInsert().Model(&rolePermissions).Exec(context.Background())
	if err != nil {
		return err
	}

	//user can only read
	userRole := &model.Role{}
	err = db.NewSelect().Model(userRole).Where("name = ?", "User").Scan(context.Background())
	if err != nil {
		return err
	}

	// Get all read permissions for user role
	var readPermissions []model.Permission
	for _, permission := range permissions {
		if permission.Name == "user:read" ||
			permission.Name == "role:read" ||
			permission.Name == "permission:read" ||
			permission.Name == "department:read" ||
			permission.Name == "branch:read" ||
			permission.Name == "organization:read" ||
			permission.Name == "holiday:read" ||
			permission.Name == "workshift:read" ||
			permission.Name == "leave:read" ||
			permission.Name == "document:read" ||
			permission.Name == "attendance:read" {
			readPermissions = append(readPermissions, permission)
		}
	}

	userRolePermissions := make([]*model.RolePermission, len(readPermissions))
	for i, permission := range readPermissions {
		userRolePermissions[i] = &model.RolePermission{
			RoleID:       userRole.ID,
			PermissionID: permission.ID,
		}
	}
	_, err = db.NewInsert().Model(&userRolePermissions).Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}
