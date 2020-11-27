package db

import (
	"context"
	"github.com/avanticaTest/maze/pkg/models"
	"github.com/stretchr/testify/mock"
)

type Mock struct{
	mock.Mock
}

func (m Mock) FindOne(ctx context.Context, db, collection string, filter, objective interface{}) error{
	args := m.Called(ctx, db, collection, filter, objective)
	return args.Error(0)
}


func (m Mock) InsertOne(ctx context.Context, db, collection string, doc interface{}) (string, error){
	args := m.Called(ctx, db, collection, doc)
	return args.String(0), args.Error(1)
}

func (m Mock) UpdateOne(ctx context.Context, filter, update interface{}, db, collection string) (int, error){
	args := m.Called(ctx,filter, update,  db, collection)
	return args.Int(0), args.Error(1)
}


func (m Mock) DeleteOne(_ context.Context, _ interface{}, db, collection string) (int, error){
	args := m.Called( db, collection)
	return args.Int(0), args.Error(1)
}

func (m Mock) DeleteMany(ctx context.Context, filter interface{}, db, col string) (int, error){
	args := m.Called(ctx,filter,  db, col)
	return args.Int(0), args.Error(1)
}

func (m Mock) FindSpots(ctx context.Context, db, col string) (result []models.Spot, err error){
	args := m.Called(ctx,  db, col)
	return nil, args.Error(1)
}

func (m Mock) FindPaths(ctx context.Context, db, col string) (result []models.Path, err error){
	args := m.Called(ctx,  db, col)
	return nil, args.Error(1)
}

func (m Mock) FindOrigin(ctx context.Context, db, col string) (result []models.Origin, err error){
	args := m.Called(ctx, db, col)
	return nil, args.Error(1)
}

func (m Mock) EstimatedDocumentCount(ctx context.Context, db, collection string) (int, error){
	args := m.Called(ctx,  db, collection)
	return args.Int(0), args.Error(1)
}

