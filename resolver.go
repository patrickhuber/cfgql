package cfgql

//go:generate go run github.com/99designs/gqlgen

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/patrickhuber/cfgql/exec"
	"github.com/patrickhuber/cfgql/models"
	"gopkg.in/yaml.v2"
)

// Resolver defines a graphql resolver for this particular schema
type Resolver struct {
	foundations []*models.Foundation
}

// Query is the entry point for graphql queries
func (r *Resolver) Query() exec.QueryResolver {

	return &queryResolver{r}
}

func (r *Resolver) Foundation() exec.FoundationResolver {
	return &foundationResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Foundations(ctx context.Context) ([]*models.Foundation, error) {
	connections, err := connections(ctx)
	if err != nil {
		return nil, err
	}
	foundations := []*models.Foundation{}
	for _, connection := range connections {
		foundation := &models.Foundation{
			ID: connection.ID,
		}
		foundations = append(foundations, foundation)
	}
	return foundations, nil
}

func connections(ctx context.Context) ([]*Connection, error) {
	dat, err := ioutil.ReadFile("foundations.yml")
	if err != nil {
		return nil, err
	}

	var connections []*Connection
	err = yaml.Unmarshal(dat, &connections)
	return connections, err
}

type foundationResolver struct{ *Resolver }

func (r *foundationResolver) Applications(ctx context.Context, f *models.Foundation) ([]*models.Application, error) {

	connections, err := connections(ctx)
	if err != nil {
		return nil, err
	}

	var connection *Connection
	for _, conn := range connections {
		if strings.Compare(*conn.ID, *f.ID) == 0 {
			connection = conn
		}
	}
	if connection == nil {
		return nil, fmt.Errorf("unable to find connection for foundation %s", *f.ID)
	}
	if connection.API == nil {
		return nil, fmt.Errorf("missing connection api for foundation %s", *f.ID)
	}

	username := os.Getenv("CF_USERNAME")
	password := os.Getenv("CF_PASSWORD")

	c := &cfclient.Config{
		ApiAddress:        *connection.API,
		Username:          username,
		Password:          password,
		SkipSslValidation: true,
	}

	client, err := cfclient.NewClient(c)
	if err != nil {
		return nil, err
	}

	apps, err := client.ListApps()
	if err != nil {
		return nil, err
	}

	applications := []*models.Application{}
	for _, app := range apps {
		id := app.Guid
		name := app.Name
		application := &models.Application{
			ID:   &id,
			Name: &name,
		}
		applications = append(applications, application)
	}
	return applications, nil
}
