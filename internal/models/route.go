package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Route struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Path      string             `bson:"path" json:"path"`
	ServiceID primitive.ObjectID `bson:"service_id" json:"service_id"`
	Plugins   []PluginConfig     `bson:"plugins" json:"plugins"`
}
