package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		// field.UUID("id", uuid.New()).
		// 	Default(uuid.New).
		// 	Immutable().
		// 	Unique(),
		field.Text("login").
			Unique(),
		field.Text("name").
			Unique(),
		field.Text("email").
			Unique(),
		field.Strings("groups"),
		field.Time("created_at").
			Immutable().
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Indexes of the User.
func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("created_at", "id"),
	}
}
