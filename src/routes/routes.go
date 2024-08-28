package routes

import (
	"github.com/gorilla/mux"

	"github.com/tam-code/lrn/src/controllers"
	"github.com/tam-code/lrn/src/middleware"
	"github.com/tam-code/lrn/src/producers"
	"github.com/tam-code/lrn/src/repositories"
)

const (
	pathPrefix     = "/api/v1"
	imagePath      = "/images"
	statisticsPath = "/statistics"
	uploadLinkPath = "/upload-link"
)

func SetupRoutes(repositories *repositories.Repositories, producers *producers.Producers) *mux.Router {
	router := mux.NewRouter()

	imageController := controllers.NewImageController(repositories, producers)
	uploadLinkController := controllers.NewUploadLinkController(repositories, pathPrefix+imagePath)
	statisticsController := controllers.NewStatisticsController(repositories)

	subrouter := router.PathPrefix(pathPrefix).Subrouter()

	subrouter.HandleFunc(imagePath+"/{upload_link_id}", imageController.UploadImage).Methods("POST")
	subrouter.HandleFunc(imagePath+"/{image_id}", imageController.GetImage).Methods("GET")

	subrouterWithSecret := router.PathPrefix(pathPrefix).Subrouter()
	subrouterWithSecret.Use(middleware.ValidateSecretToken)

	subrouterWithSecret.HandleFunc(statisticsPath, statisticsController.GetStatistics).Methods("GET")
	subrouterWithSecret.HandleFunc(uploadLinkPath, uploadLinkController.CreateUploadLink).Methods("POST")

	return router
}
