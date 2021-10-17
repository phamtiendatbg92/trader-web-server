package main

import (
	"fmt"
	"net/http"
	"trader-web-api/conf"
	"trader-web-api/routers"

	"go.uber.org/zap"
)

func init() {
	conf.InitLogger()
	conf.InitConfig()
}

// func main() {
// 	//db - controller.connect()
// 	//dbcontroller.Connect()
// 	// Create router
// 	router := mux.NewRouter()

// 	router.HandleFunc("/api/v1/get-list-tutorials", getAllTutorial).Methods(http.MethodGet)
// 	router.HandleFunc("/api/v1/detail-tutorial/{url}", getDetailTutorial).Methods(http.MethodGet)
// 	router.HandleFunc("/api/v1/get-hashtag", getAllHashTag).Methods(http.MethodGet)
// 	router.HandleFunc("/api/v1/upload-new-post", uploadNewPost).Methods(http.MethodPost, http.MethodOptions)
// 	router.HandleFunc("/api/v1/update-post", updateTutorial).Methods(http.MethodPut, http.MethodOptions)
// 	router.HandleFunc("/api/v1/delete-post/{id}", deleteTutorial).Methods(http.MethodDelete, http.MethodOptions)

// 	// router.PathPrefix("/images").Handler(http.FileServer(http.Dir("./public/images")))

// 	// This will serve files under http://localhost:8000/static/<filename>
// 	router.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("./public/images/"))))

// 	router.Use(accessControlMiddleware)
// 	// start Server
// 	log.Fatal(http.ListenAndServe(":5000", router))
// }
func main() {
	bindingAddr := fmt.Sprintf(":%v", conf.EnvConfig.Server.BindingPort)
	zap.S().Info("Starting api at: ", bindingAddr)
	router, err := routers.InitRouter()
	if err != nil {
		zap.S().Error("Error when init router, detail: ", err)
		panic(err)
	}
	_ = router.Run(bindingAddr)
}

// access control and  CORS middleware
func accessControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS, PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}
