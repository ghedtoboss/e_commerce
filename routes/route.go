package routes

import (
	"e_commerce/controller"
	"e_commerce/middleware"
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func InitRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/users/register", controller.RegisterHandler)                                                                                         //++
	r.HandleFunc("/users/login", controller.LoginHandler)                                                                                               //++
	r.Handle("/users/profile", middleware.JWTAuth(http.HandlerFunc(controller.GetProfile))).Methods("GET")                                              //++
	r.Handle("/users/profile", middleware.JWTAuth(http.HandlerFunc(controller.UpdateProfile))).Methods("PUT")                                           //++
	r.Handle("/users/profile/password", middleware.JWTAuth(http.HandlerFunc(controller.UpdatePassword))).Methods("PUT")                                 //++
	r.Handle("/users/{user_id}/delete", middleware.JWTAuth(middleware.Authorize("admin")(http.HandlerFunc(controller.DeleteUser)))).Methods("DELETE")   //++
	r.Handle("/users/close-account", middleware.JWTAuth(middleware.Authorize("customer")(http.HandlerFunc(controller.CloseAccount)))).Methods("DELETE") //++

	r.Handle("/shop", middleware.JWTAuth(middleware.Authorize("seller")(http.HandlerFunc(controller.CreateShop)))).Methods("POST")  //++
	r.Handle("/shop", middleware.JWTAuth(middleware.Authorize("seller")(http.HandlerFunc(controller.UpdateShop)))).Methods("PUT")   //++
	r.Handle("/shop/my", middleware.JWTAuth(middleware.Authorize("seller")(http.HandlerFunc(controller.GetMyShop)))).Methods("GET") //--
	r.Handle("/shop/{shop_id}", middleware.JWTAuth(http.HandlerFunc(controller.GetShop))).Methods("GET")                            //++

	r.Handle("/product", middleware.JWTAuth(middleware.Authorize("seller")(http.HandlerFunc(controller.AddProduct)))).Methods("POST")                     //++
	r.Handle("/product/{product_id}", middleware.JWTAuth(middleware.Authorize("seller")(http.HandlerFunc(controller.UpdateProduct)))).Methods("PUT")      //++
	r.Handle("/product/my-products", middleware.JWTAuth(middleware.Authorize("seller")(http.HandlerFunc(controller.GetProductsByMyShop)))).Methods("GET") //++
	r.Handle("/product/{product_id}", http.HandlerFunc(controller.GetProduct)).Methods("GET")                                                             //++
	r.Handle("/product", http.HandlerFunc(controller.GetProducts)).Methods("GET")                                                                         //++
	r.Handle("/product/{shop_id}/products", http.HandlerFunc(controller.GetProductsByShop)).Methods("GET")                                                //++

	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/docs/swagger.json"), // The url pointing to API definition
	))

	// Static files endpoint for serving the swagger docs
	fs := http.FileServer(http.Dir("./docs"))
	r.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", fs))

	return r
}
