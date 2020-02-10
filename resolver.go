package cfgql

//go:generate go run github.com/99designs/gqlgen

import (
	"context"
	"io/ioutil"

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
	dat, err := ioutil.ReadFile("foundations.yml")
	if err != nil {
		return nil, err
	}

	var connections []*Connection
	err = yaml.Unmarshal(dat, &connections)
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

type foundationResolver struct{ *Resolver }

func (r *foundationResolver) Applications(ctx context.Context, f *models.Foundation) ([]*models.Application, error) {
	id := "abc123"
	name := "hello"
	application := &models.Application{
		ID:   &id,
		Name: &name,
	}
	return []*models.Application{application}, nil
}
