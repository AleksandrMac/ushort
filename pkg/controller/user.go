package controller

// func (env *Env) setUserHandlers(r *chi.Mux) {
// 	r.Route("/user", func(r chi.Router) {
// 		r.Get("/", env.userList)
// 		r.Post("/", env.createUser)

// 		r.Route("/{userId}", func(r chi.Router) {
// 			r.Get("/", env.getUser)
// 			r.Patch("/", env.updateUser)
// 			r.Delete("/", env.deleteUser)
// 			r.Post("/", env.createUserToken)
// 			r.Delete("/{api_key}", env.deleteUserToken)
// 		})
// 	})
// }

// func (env *Env) userList(w http.ResponseWriter, r *http.Request) {
// }
// func (env *Env) createUser(w http.ResponseWriter, r *http.Request) {

// }
// func (env *Env) getUser(w http.ResponseWriter, r *http.Request)         {}
// func (env *Env) deleteUser(w http.ResponseWriter, r *http.Request)      {}
// func (env *Env) updateUser(w http.ResponseWriter, r *http.Request)      {}
// func (env *Env) createUserToken(w http.ResponseWriter, r *http.Request) {}
// func (env *Env) deleteUserToken(w http.ResponseWriter, r *http.Request) {}
