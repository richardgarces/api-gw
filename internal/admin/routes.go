package admin

import (
	"api-gw/internal/db"
	"api-gw/internal/models"
	"api-gw/internal/utils"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// POST /admin/routes   {path, service_id, plugins}
// GET  /admin/routes
// GET  /admin/routes/{id}
// PUT  /admin/routes/{id}
// DELETE /admin/routes/{id}

func RoutesHandler(w http.ResponseWriter, r *http.Request) {
	collection := db.GetClient().Database(os.Getenv("MONGO_DATABASE")).Collection("routes")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch r.Method {
	case http.MethodPost:
		var route models.Route
		if err := json.NewDecoder(r.Body).Decode(&route); err != nil {
			utils.Error(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		res, err := collection.InsertOne(ctx, route)
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		route.ID = res.InsertedID.(primitive.ObjectID)
		utils.JSON(w, http.StatusCreated, route)

	case http.MethodGet:
		id := r.URL.Query().Get("id")
		if id != "" {
			objID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				utils.Error(w, http.StatusBadRequest, "ID inválido")
				return
			}
			var route models.Route
			err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&route)
			if err != nil {
				utils.Error(w, http.StatusNotFound, "Ruta no encontrada 4")
				return
			}
			utils.JSON(w, http.StatusOK, route)
			return
		}
		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		defer cursor.Close(ctx)
		var routes []models.Route
		for cursor.Next(ctx) {
			var route models.Route
			if err := cursor.Decode(&route); err == nil {
				routes = append(routes, route)
			}
		}
		utils.JSON(w, http.StatusOK, routes)

	case http.MethodPut:
		id := r.URL.Query().Get("id")
		if id == "" {
			utils.Error(w, http.StatusBadRequest, "ID requerido")
			return
		}
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "ID inválido")
			return
		}
		var route models.Route
		if err := json.NewDecoder(r.Body).Decode(&route); err != nil {
			utils.Error(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		_, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": route})
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.JSON(w, http.StatusOK, map[string]string{"result": "Ruta actualizada"})

	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if id == "" {
			utils.Error(w, http.StatusBadRequest, "ID requerido")
			return
		}
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "ID inválido")
			return
		}
		_, err = collection.DeleteOne(ctx, bson.M{"_id": objID})
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.JSON(w, http.StatusOK, map[string]string{"result": "Ruta eliminada"})

	default:
		utils.Error(w, http.StatusMethodNotAllowed, "Método no permitido")
	}
}
