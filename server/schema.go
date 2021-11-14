package server

import (
	"github.com/graphql-go/graphql"
)

var UserType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"username": &graphql.Field{
			Type: graphql.String,
		},
		"password": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var LocationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Location",
	Fields: graphql.Fields{
		"lat": &graphql.Field{
			Type: graphql.Float,
		},
		"lng": &graphql.Field{
			Type: graphql.Float,
		},
	},
})
