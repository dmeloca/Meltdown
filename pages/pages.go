package pages

import (
	"crypto/rand"
	"encoding/hex"
	"html/template"
	"log"
	"meltdown/contextKeys"
	"meltdown/session"
	"meltdown/users"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "No se pudo cargar la plantilla", http.StatusInternalServerError)
		return
	}
	err1 := tmpl.Execute(w, nil)
	if err1 != nil {
		return
	}
}

func generateSessionID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")
		loadedUsers, err := users.LoadUsers()
		if err != nil {
			http.Error(w, "Error al cargar los usuarios", http.StatusInternalServerError)
			return
		}
		for _, user := range loadedUsers {
			if user.Username == username && user.Password == password {
				sessionID := generateSessionID()

				// Guardar la sesión
				session.Set(sessionID, user.Username, user.Name, user.Role)

				// Crear cookie
				http.SetCookie(w, &http.Cookie{
					Name:     "session_id",
					Value:    sessionID,
					Path:     "/",
					HttpOnly: true,
					MaxAge:   3600, // 1 hora
				})
				log.Printf("Usuario autenticado: %s", username)
				http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
				return
			}
		}
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`<script>alert("Credenciales inválidas"); window.location.href = "/login";</script>`))
		// http.Error(w, "Credenciales inválidas", http.StatusUnauthorized)
		return
	}
	// GET: mostrar login
	tmpl, _ := template.ParseFiles("templates/login.html")
	tmpl.Execute(w, nil)
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/dashboard.html")
	username := r.Context().Value(contextKeys.UserContextKey).(string)
	name := r.Context().Value(contextKeys.NameContextKey).(string)
	role := r.Context().Value(contextKeys.RoleContextKey).(string)
	if err != nil {
		http.Error(w, "No se pudo cargar la plantilla de dashboard", http.StatusInternalServerError)
		return
	}
	data := map[string]string{"Username": username, "Name": name, "Role": role}
	if role == "admin" {
		data["AdminMessage"] = "Bienvenido, administrador"
	} else if role == "user" {
		data["UserMessage"] = "Bienvenido, usuario"
	} else {
		data["Message"] = "Bienvenido, invitado"
	}
	err1 := tmpl.Execute(w, data)
	if err1 != nil {
		return
	}
}
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/Dashboard/profile.html")

	if err != nil {
		http.Error(w, "No se pudo cargar la página de perfil", http.StatusInternalServerError)
		return
	}
	username := r.Context().Value(contextKeys.UserContextKey).(string)
	name := r.Context().Value(contextKeys.NameContextKey).(string)
	role := r.Context().Value(contextKeys.RoleContextKey).(string)
	data := map[string]string{"Username": username, "Name": name, "Role": role}
	tmpl.Execute(w, data)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/404.html")
	if err != nil {
		http.Error(w, "Página no encontrada", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	tmpl.Execute(w, nil)
}
func SettingsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/Dashboard/settings.html")
		if err != nil {
			http.Error(w, "Error al cargar la página de configuración", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		// Aquí puedes procesar los datos del formulario
		nombre := r.FormValue("nombre")
		email := r.FormValue("email")
		password := r.FormValue("password")
		w.Write([]byte(`<script>alert("Cambios Recibidos"); window.location.href = "/dashboard";</script>`))
		log.Printf("Cambios recibidos - Nombre: %s, Email: %s, Contraseña: %s", nombre, email, password)
		// Redirigir después de guardar (puedes mostrar un mensaje también)
		http.Redirect(w, r, "/dashboard/profile", http.StatusSeeOther)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		// Borrar del mapa
		session.Delete(cookie.Value)
	}
	username := r.Context().Value(contextKeys.UserContextKey).(string)
	log.Printf("Usuario desconectado %s", username)
	// Eliminar cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
