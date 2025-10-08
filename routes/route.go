package routes

import (
	middleware "crud-app/Middleware"
	service "crud-app/app/Service"
	"crud-app/app/repository"
	"net/http"

	"github.com/gorilla/mux"
)

func UserRoutes(r *mux.Router, PekerjaanService *service.PekerjaanService, alumniService *service.AlumniService, authService *service.AuthService, userRepo *repository.UserRepository, userService *service.UserService) {
	r.HandleFunc("/register", authService.Register).Methods("POST")
	r.HandleFunc("/login", authService.Login).Methods("POST")

	// Alumni routes
	r.Handle("/alumni/{id}", middleware.AuthMiddleware(*userRepo, http.HandlerFunc(alumniService.GetByID))).Methods("GET")

	r.Handle("/alumni", middleware.AuthMiddleware(*userRepo,
		middleware.RoleMiddleware("admin", http.HandlerFunc(alumniService.Create)))).Methods("POST")

	r.Handle("/alumni/{id}", middleware.AuthMiddleware(*userRepo,
		middleware.RoleMiddleware("admin", http.HandlerFunc(alumniService.Update)))).Methods("PUT")

	r.Handle("/alumni/{id}", middleware.AuthMiddleware(*userRepo,
		middleware.RoleMiddleware("admin", http.HandlerFunc(alumniService.Delete)))).Methods("DELETE")

	// Pekerjaan routes
	r.Handle("/pekerjaan/{alumni_id}", middleware.AuthMiddleware(*userRepo, http.HandlerFunc(PekerjaanService.GetByAlumni))).Methods("GET")
	r.Handle("/pekerjaan", middleware.AuthMiddleware(*userRepo,
		middleware.RoleMiddleware("admin", http.HandlerFunc(PekerjaanService.Create)))).Methods("POST")
	r.Handle("/pekerjaan/{id}", middleware.AuthMiddleware(*userRepo,
		middleware.RoleMiddleware("admin", http.HandlerFunc(PekerjaanService.Update)))).Methods("PUT")
	// r.Handle("/pekerjaan/{id}", middleware.AuthMiddleware(*userRepo,
	// 	middleware.RoleMiddleware("admin", http.HandlerFunc(PekerjaanService.Delete)))).Methods("DELETE")
	r.Handle("/pekerjaan/{id}", middleware.AuthMiddleware(*userRepo, http.HandlerFunc(PekerjaanService.GetByID))).Methods("GET")
	

	r.Handle("/Users/{id}", middleware.AuthMiddleware(*userRepo, http.HandlerFunc(userService.SoftDeleteUser))).Methods("DELETE")
	// Routing with pagination , sort by dll
	r.Handle("/users", middleware.AuthMiddleware(*userRepo, http.HandlerFunc(userService.GetUsers))).Methods("GET")
	r.Handle("/alumni", middleware.AuthMiddleware(*userRepo, http.HandlerFunc(alumniService.GetAlumni))).Methods("GET")
	r.Handle("/pekerjaan", middleware.AuthMiddleware(*userRepo, http.HandlerFunc(PekerjaanService.GetPekerjaan))).Methods("GET")
	

r.Handle("/pekerjaan/{id}", 
    middleware.AuthMiddleware(*userRepo, http.HandlerFunc(PekerjaanService.SoftDeletePekerjaan)),
).Methods("DELETE")



// package routes

// import (
// 	middleware "crud-app/Middleware"
// 	service "crud-app/app/Service"
// 	"crud-app/app/repository"
// 	"net/http"

// 	"github.com/gorilla/mux"
// )

// func UserRoutes(
// 	r *mux.Router,
// 	PekerjaanService *service.PekerjaanService,
// 	alumniService *service.AlumniService,
// 	authService *service.AuthService,
// 	userRepo *repository.UserRepository,
// 	userService *service.UserService,
// ) {
// 	r.HandleFunc("/register", authService.Register).Methods("POST")
// 	r.HandleFunc("/login", authService.Login).Methods("POST")

// 	// Alumni routes
// 	r.Handle("/alumni/{id}", middleware.AuthMiddleware(*userRepo, http.HandlerFunc(alumniService.GetByID))).Methods("GET")

// 	r.Handle("/alumni", middleware.AuthMiddleware(*userRepo,
// 		middleware.RoleMiddleware("admin", http.HandlerFunc(alumniService.Create)))).Methods("POST")

// 	r.Handle("/alumni/{id}", middleware.AuthMiddleware(*userRepo,
// 		middleware.RoleMiddleware("admin", http.HandlerFunc(alumniService.Update)))).Methods("PUT")

// 	r.Handle("/alumni/{id}", middleware.AuthMiddleware(*userRepo,
// 		middleware.RoleMiddleware("admin", http.HandlerFunc(alumniService.Delete)))).Methods("DELETE")

// 	// ========================
// 	// Pekerjaan routes
// 	// ========================

// 	// ambil semua pekerjaan milik alumni tertentu
// 	r.Handle("/pekerjaan/{alumni_id}",
// 		middleware.AuthMiddleware(*userRepo, http.HandlerFunc(PekerjaanService.GetByAlumni)),
// 	).Methods("GET")

// 	// tambah pekerjaan (admin only)
// 	r.Handle("/pekerjaan",
// 		middleware.AuthMiddleware(*userRepo, middleware.RoleMiddleware("admin", http.HandlerFunc(PekerjaanService.Create))),
// 	).Methods("POST")

// 	// update pekerjaan (admin only)
// 	r.Handle("/pekerjaan/{id}",
// 		middleware.AuthMiddleware(*userRepo, middleware.RoleMiddleware("admin", http.HandlerFunc(PekerjaanService.Update))),
// 	).Methods("PUT")

// 	// SOFT DELETE oleh USER (hanya hapus pekerjaannya sendiri)
// 	r.Handle("/pekerjaan/user/{id}",
// 		middleware.AuthMiddleware(*userRepo, http.HandlerFunc(PekerjaanService.DeleteByUser)),
// 	).Methods("DELETE")

// 	// SOFT DELETE oleh ADMIN (hapus semua pekerjaan)
// 	r.Handle("/pekerjaan/admin/{id}",
// 		middleware.AuthMiddleware(*userRepo, middleware.RoleMiddleware("admin", http.HandlerFunc(PekerjaanService.DeleteByAdmin))),
// 	).Methods("DELETE")

// 	// ========================
// 	// Extra routes
// 	// ========================
// 	r.Handle("/users", middleware.AuthMiddleware(*userRepo, http.HandlerFunc(userService.GetUsers))).Methods("GET")
// 	r.Handle("/alumni", middleware.AuthMiddleware(*userRepo, http.HandlerFunc(alumniService.GetAlumni))).Methods("GET")
// 	r.Handle("/PekerjaanAlumni", middleware.AuthMiddleware(*userRepo, http.HandlerFunc(PekerjaanService.GetPekerjaan))).Methods("GET")
// }
	
	// TRASH MANAGEMENT ROUTES
	// Get trash data
	r.Handle("/trash/pekerjaan", 
		middleware.AuthMiddleware(*userRepo, http.HandlerFunc(PekerjaanService.GetTrash)),
	).Methods("GET")
	
	// Restore from trash
	r.Handle("/trash/pekerjaan/{id}/restore", 
		middleware.AuthMiddleware(*userRepo, http.HandlerFunc(PekerjaanService.RestorePekerjaan)),
	).Methods("PUT")
	
	// Hard delete (permanen hapus)
	r.Handle("/trash/pekerjaan/{id}/hard-delete", 
		middleware.AuthMiddleware(*userRepo, http.HandlerFunc(PekerjaanService.HardDeletePekerjaan)),
	).Methods("DELETE")
}