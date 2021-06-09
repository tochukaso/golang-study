package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/tochukaso/golang-study/graph/generated"
	"github.com/tochukaso/golang-study/graph/model"
)

func (r *mutationResolver) CreateProduct(ctx context.Context, input model.NewProduct) (*model.Product, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *productResolver) ID(ctx context.Context, obj *model.Product) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *productResolver) CreatedAt(ctx context.Context, obj *model.Product) (*string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *productResolver) UpdatedAt(ctx context.Context, obj *model.Product) (*string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *productResolver) DeletedAt(ctx context.Context, obj *model.Product) (*string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ListProducts(ctx context.Context) ([]*model.Product, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetProduct(ctx context.Context, id string) (*model.Product, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Product returns generated.ProductResolver implementation.
func (r *Resolver) Product() generated.ProductResolver { return &productResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type productResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
