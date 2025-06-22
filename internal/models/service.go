package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type PluginConfig struct {
	Type   string                 `bson:"type" json:"type"`
	Config map[string]interface{} `bson:"config" json:"config"`
}

type Service struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name    string             `bson:"name" json:"name"`
	Targets []string           `bson:"targets" json:"targets"`
	Plugins []PluginConfig     `bson:"plugins" json:"plugins"`
}
