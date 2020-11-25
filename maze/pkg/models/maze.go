package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Spot struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	XCoordinate float64            `json:"x_coordinate,omitempty" bson:"x_coordinate,omitempty"`
	YCoordinate float64            `json:"y_coordinate,omitempty" bson:"y_coordinate,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Number      int                `json:"number,omitempty" bson:"number,omitempty"`
}

type Path struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	PointA   primitive.ObjectID `json:"point_a,omitempty" bson:"point_a,omitempty"`
	PointB   primitive.ObjectID `json:"point_b,omitempty" bson:"point_b,omitempty"`
	Distance float64            `json:"distance,omitempty" bson:"distance,omitempty"`
}

type Origin struct {
	XOrigin float64 `json:"x_origin,omitempty" bson:"x_origin,omitempty"`
	YOrigin float64 `json:"y_origin,omitempty" bson:"y_origin,omitempty"`
}

type Quadrant struct {
	Quadrant string `json:"name"`
}

type CreatePathRequest struct{
	PointA string `json:"point_a"`
	PointB string `json:"point_b"`
}

type CreateObjectResponse struct {
	ID string `json:"id"`
}

type ModifyObjectResponse struct {
	AffectedItems int `json:"affected_items"`
}
