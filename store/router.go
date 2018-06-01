package store

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var controller = &Controller{Repository: Repository{}}

// Route defines a route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes defines the list of routes of our API
type Routes []Route

var routes = Routes{
	Route{
		"Login",
		"POST",
		"/Login",
		controller.Login,
	},
	Route{
		"Ideas",
		"GET",
		"/Ideas",
		controller.GetIdeas,
	},
	Route{
		"AddIdea",
		"POST",
		"/Ideas",
		AuthenticationMiddleware(controller.AddIdea),
	},
	Route{
		"UpdateIdea",
		"PUT",
		"/Ideas",
		AuthenticationMiddleware(controller.UpdateIdea),
	},
	// Get Product by username {name}
	Route{
		"GetProductByUsername",
		"GET",
		"/Ideas/User/{name}",
		controller.GetIdeasByUsername,
	},
	// Delete Product by {id}
	Route{
		"DeleteProduct",
		"DELETE",
		"/Ideas",
		AuthenticationMiddleware(controller.DeleteIdea),
	},
	// Search product with string
	Route{
		"SearchProduct",
		"GET",
		"/Ideas/Search/{query}",
		controller.GetIdeasByString,
	}}

// NewRouter configures a new router to the API
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		log.Println(route.Name)
		handler = route.HandlerFunc

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}
