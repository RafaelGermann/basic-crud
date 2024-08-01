package servidor

import (
	"basic-crud/banco"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
)

type user struct {
	ID    uint32 `json:"id"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	body, erro := io.ReadAll(r.Body)
	if erro != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("erro no body")))
		return
	}

	var user user

	if erro = json.Unmarshal(body, &user); erro != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Erro unmarshal")))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("erro ao conectar ao banco")))
		return
	}
	defer db.Close()

	var lastInsertId int
	erro = db.QueryRow("insert into usuarios (nome, email) values ($1,$2) returning id", user.Nome, user.Email).Scan(&lastInsertId)
	if erro != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(erro.Error())))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Usuario inserido com sucesso Id: %d", lastInsertId)))
}

func BuscarUsuarios(w http.ResponseWriter, r *http.Request) {
	db, erro := banco.Conectar()
	if erro != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("erro ao conectar ao banco")))
		return
	}
	defer db.Close()

	linhas, erro := db.Query("select * from usuarios")
	if erro != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("erro buscar usuarios")))
		return
	}
	defer linhas.Close()

	var usuarios []user
	for linhas.Next() {
		var user user
		if erro := linhas.Scan(&user.ID, &user.Nome, &user.Email); erro != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("erro ao fazer scan usuario")))
			return
		}
		usuarios = append(usuarios, user)
	}

	w.WriteHeader(http.StatusOK)
	if erro = json.NewEncoder(w).Encode(usuarios); erro != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("erro ao fazer o encoder")))
		return
	}
}

func BuscarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	ID, erro := strconv.ParseUint(parametros["id"], 10, 32)
	if erro != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("erro ao converter para inteiro")))
		return
	}
	db, erro := banco.Conectar()
	if erro != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("erro ao conectar ao banco")))
		return
	}
	defer db.Close()

	linha, erro := db.Query("select * from usuarios where id = $1", ID)
	if erro != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("erro ao conectar ao banco")))
		return
	}

	defer linha.Close()

	var usuario user
	if linha.Next() {
		if erro = linha.Scan(&usuario.ID, &usuario.Nome, &usuario.Email); erro != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("erro ao escanear usuario")))
			return
		}
	}

	if erro := json.NewEncoder(w).Encode(usuario); erro != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("erro ao converter usuario para json")))
		return
	}
}

func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	ID, erro := strconv.ParseUint(parametros["id"], 10, 32)
	if erro != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("erro ao converter para inteiro")))
		return
	}

	body, erro := io.ReadAll(r.Body)
	if erro != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("erro no body")))
		return
	}

	var user user

	if erro = json.Unmarshal(body, &user); erro != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Erro unmarshal")))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("erro ao conectar ao banco")))
		return
	}
	defer db.Close()

	if erro = db.QueryRow("update usuarios set nome = $1, email = $2 where id = $3", user.Nome, user.Email, ID).Err(); erro != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(erro.Error())))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func DeletarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	ID, erro := strconv.ParseUint(parametros["id"], 10, 32)
	if erro != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("erro ao converter para inteiro")))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("erro ao conectar ao banco")))
		return
	}
	defer db.Close()
	if erro = db.QueryRow("delete from usuarios where id = $1", ID).Err(); erro != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(erro.Error())))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
