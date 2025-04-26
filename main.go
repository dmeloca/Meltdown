package main

import (
	"context"
	"log"
	"meltdown/contextKeys"
	"meltdown/pages"
	"meltdown/session"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			pages.IndexHandler(w, r)
		case "/login":
			pages.LoginHandler(w, r)
		case "/dashboard":
			withAuth(pages.DashboardHandler)(w, r)
		case "/dashboard/profile":
			withAuth(pages.ProfileHandler)(w, r)
		case "/dashboard/settings":
			withAuth(pages.SettingsHandler)(w, r)
		case "/logout":
			withAuth(pages.LogoutHandler)(w, r)
		default:
			pages.NotFoundHandler(w, r)
		}
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Println("Servidor corriendo en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func withAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Verificar que la sesión exista usando el paquete `session`
		data, exists := session.Get(cookie.Value)

		if !exists {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeys.UserContextKey, data.Username)
		ctx = context.WithValue(ctx, contextKeys.NameContextKey, data.Name)
		ctx = context.WithValue(ctx, contextKeys.RoleContextKey, data.Role)
		r = r.WithContext(ctx)
		// Aquí puedes también pasar el username por contexto si quieres
		handler.ServeHTTP(w, r)
	}
}
