package server

import (
	"errors"
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/karthiklsarma/cedar-logging/logging"
	"github.com/karthiklsarma/cedar-schema/gen"
	"github.com/karthiklsarma/cedar-server/stream"
)

func StartGraphQlServer() graphql.Schema {
	queryFields := graphql.Fields{
		"getUsers": &graphql.Field{
			Type:        graphql.NewList(UserType),
			Description: "List of users",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return UsersList, nil
			},
		},
		"getLocations": &graphql.Field{
			Type:        graphql.NewList(LocationType),
			Description: "Location of users",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"username": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"group": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				username, ok := p.Args["username"].(string)
				if !ok {
					return nil, errors.New("username not provided or invalid username")
				}

				logging.Debug(fmt.Sprintf("will query locations for user: %s", username))
				// TODO: Query last location of user
				return LocationList, nil
			},
		},
	}

	mutationFields := graphql.Fields{
		"setUsers": &graphql.Field{
			Type:        UserType,
			Description: "Add New User",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"username": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"password": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, _ := p.Args["id"].(string)
				username, _ := p.Args["username"].(string)

				newUser := User{
					Id:       id,
					Username: username,
				}

				UsersList = append(UsersList, newUser)
				return newUser, nil
			},
		},
		"setLocation": &graphql.Field{
			Type:        LocationType,
			Description: "location of the user",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"lat": &graphql.ArgumentConfig{
					Type: graphql.Float,
				},
				"lng": &graphql.ArgumentConfig{
					Type: graphql.Float,
				},
				"timestamp": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"device": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id := p.Args["id"].(string)
				lat := p.Args["lat"].(float64)
				lng := p.Args["lng"].(float64)
				timestamp := p.Args["timestamp"].(int)
				device := p.Args["device"].(string)
				logging.Info(fmt.Sprintf("Received id: %v, lat: %v, lng: %v, timestamp: %v, device: %v", id, lat, lng, timestamp, device))
				location := &gen.Location{Id: id, Lat: lat, Lng: lng, Timestamp: int64(timestamp), Device: device}
				stream.EmitLocation(location)
				logging.Info(fmt.Sprintf("sending location message : %v to eventqueue", location))
				return location, nil
			},
		},
	}

	rootQuery := graphql.ObjectConfig{
		Name:        "Query",
		Description: "Root Query",
		Fields:      queryFields,
	}

	rootMutation := graphql.ObjectConfig{
		Name:        "Mutation",
		Description: "Root Mutation",
		Fields:      mutationFields,
	}

	schemaConfig := graphql.SchemaConfig{
		Query:    graphql.NewObject(rootQuery),
		Mutation: graphql.NewObject(rootMutation),
	}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		logging.Fatal("Invalid graphQl schema")
	}

	return schema
}
