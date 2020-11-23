package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Spot struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	XCoordinate float64            `json:"x_coordinate,omitempty" bson:"x_coordinate,omitempty"`
	YCoordinate float64            `json:"y_coordinate,omitempty" bson:"y_coordinate,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Number      int                `json:"number,omitempty" bson:"number,omitempty"`
}

type Path struct{
	ID   primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	PointA primitive.ObjectID `json:"point_a,omitempty" bson:"point_a,omitempty"`
	PointB primitive.ObjectID `json:"point_b,omitempty" bson:"point_b,omitempty"`
}