package server

import (
	"errors"
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/karthiklsarma/cedar-logging/logging"
	"github.com/karthiklsarma/cedar-schema/gen"
	"github.com/karthiklsarma/cedar-server/stream"
)

type IResolver interface {
	ConnectEventQueue() error
	GetSchema() graphql.Schema
}

type GraphQlResolver struct {
	eventQueue stream.IEventQueue
}

func (resolver *GraphQlResolver) GetSchema() graphql.Schema {
	if resolver.eventQueue == nil {
		logging.Fatal("Event Queue is nil for Resolver. Please connect to event queue before fetching schema.")
	}

	queryFields := graphql.Fields{
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
		"authenticate": &graphql.Field{
			Type:        graphql.String,
			Description: "Authenticate User",
			Args: graphql.FieldConfigArgument{
				"username": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"password": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				username, ok := p.Args["username"].(string)
				if !ok {
					return "no crenetials present.", errors.New("username not provided or invalid username")
				}

				password, ok := p.Args["password"].(string)
				if !ok {
					return "no credentials present.", errors.New("password not provided or invalid username")
				}

				status, err := AuthenticateUser(username, password)
				if !status {
					return "invalid Credentials", errors.New("Invalid credentials ! login failed.")
				}

				if err != nil {
					return "internal server error", errors.New("Unable to authenticate at the moment. Something went wrong !")
				}

				return GetNewToken(username), nil
			},
		},
	}

	mutationFields := graphql.Fields{
		"addUser": &graphql.Field{
			Type:        UserType,
			Description: "Add New User",
			Args: graphql.FieldConfigArgument{
				"username": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"firstname": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"lastname": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"password": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"email": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"phone": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				username, _ := p.Args["username"].(string)
				firstname, _ := p.Args["firstname"].(string)
				lastname, _ := p.Args["lastname"].(string)
				password, _ := p.Args["password"].(string)
				email, _ := p.Args["email"].(string)
				phone, _ := p.Args["phone"].(string)

				newUser := &gen.User{
					Firstname: firstname,
					Lastname:  lastname,
					Username:  username,
					Password:  password,
					Email:     email,
					Phone:     phone,
				}

				if status, err := RegisterUser(newUser); !status {
					logging.Error(fmt.Sprintf("failed to register user. Error: {%v}", err))
					return nil, err
				}

				return newUser, nil
			},
		},
		"sendLocation": &graphql.Field{
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
				location := &gen.Location{Id: id, Lat: lat, Lng: lng, Timestamp: uint64(timestamp), Device: device}

				if err := resolver.eventQueue.EmitLocation(location); err != nil {
					return nil, err
				}

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

func (resolver *GraphQlResolver) ConnectEventQueue() error {
	resolver.eventQueue = &stream.EventQueue{}
	if err := resolver.eventQueue.Connect(); err != nil {
		logging.Fatal("Failed to connect to event queue.")
		return err
	}
	return nil
}
