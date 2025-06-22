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

// POST /admin/services   {name, targets, plugins}
// GET  /admin/services
// GET  /admin/services/{id}
// PUT  /admin/services/{id}
// DELETE /admin/services/{id}

func ServicesHandler(w http.ResponseWriter, r *http.Request) {
	collection := db.GetClient().Database(os.Getenv("MONGO_DATABASE")).Collection("services")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch r.Method {
	case http.MethodPost:
		var svc models.Service
		if err := json.NewDecoder(r.Body).Decode(&svc); err != nil {
			utils.Error(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		res, err := collection.InsertOne(ctx, svc)
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		svc.ID = res.InsertedID.(primitive.ObjectID)
		utils.JSON(w, http.StatusCreated, svc)

	case http.MethodGet:
		// Si hay un id en la URL, buscar uno solo
		id := r.URL.Query().Get("id")
		if id != "" {
			objID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				utils.Error(w, http.StatusBadRequest, "ID inválido")
				return
			}
			var svc models.Service
			err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&svc)
			if err != nil {
				utils.Error(w, http.StatusNotFound, "Servicio no encontrado")
				return
			}
			utils.JSON(w, http.StatusOK, svc)
			return
		}
		// Si no, listar todos
		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		defer cursor.Close(ctx)
		var services []models.Service
		for cursor.Next(ctx) {
			var svc models.Service
			if err := cursor.Decode(&svc); err == nil {
				services = append(services, svc)
			}
		}
		utils.JSON(w, http.StatusOK, services)

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
		var svc models.Service
		if err := json.NewDecoder(r.Body).Decode(&svc); err != nil {
			utils.Error(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		_, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": svc})
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.JSON(w, http.StatusOK, map[string]string{"result": "Servicio actualizado"})

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
		utils.JSON(w, http.StatusOK, map[string]string{"result": "Servicio eliminado"})

	default:
		utils.Error(w, http.StatusMethodNotAllowed, "Método no permitido")
	}
}
