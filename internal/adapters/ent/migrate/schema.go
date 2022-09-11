// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// UsersColumns holds the columns for the "users" table.
	UsersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "login", Type: field.TypeString, Unique: true, Size: 2147483647},
		{Name: "name", Type: field.TypeString, Unique: true, Size: 2147483647},
		{Name: "email", Type: field.TypeString, Unique: true, Size: 2147483647},
		{Name: "groups", Type: field.TypeJSON},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
	}
	// UsersTable holds the schema information for the "users" table.
	UsersTable = &schema.Table{
		Name:       "users",
		Columns:    UsersColumns,
		PrimaryKey: []*schema.Column{UsersColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "user_created_at_id",
				Unique:  false,
				Columns: []*schema.Column{UsersColumns[5], UsersColumns[0]},
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		UsersTable,
	}
)

func init() {
}
