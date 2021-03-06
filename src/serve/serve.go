package serve

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"

	"client"
	"common"
)

var session string
var modByToken = make(map[string]common.Mod)

type lmError struct {
	Error string `json:"error"`
}

type postForm struct {
	Id string
}

func sendError(w http.ResponseWriter, status int, err error) bool {
	if err != nil {
		w.WriteHeader(status)
		lmErr := lmError{fmt.Sprint(err)}
		data, _ := json.Marshal(lmErr)
		w.Write(data)
		return true
	}
	return false
}

func auth(w http.ResponseWriter, r *http.Request) (mod common.Mod, authed bool) {
	if token, ok := r.Header["X-Token"]; ok {
		mod, ok = modByToken[token[0]]
		if ok {
			return mod, true
		}
	}
	sendError(w, http.StatusForbidden, errors.New("unknown token"))
	return mod, false
}

func ensurePrefix(w http.ResponseWriter, mod common.Mod, post client.Post) bool {
	msg := strings.TrimSpace(post.Message)
	idx := strings.Index(msg, mod.Prefix)
	nl := strings.Index(msg, "\n")
	if idx >= 0 && (nl < 0 || idx < nl) {
		return true
	}
	sendError(w, http.StatusForbidden, errors.New("bad OP prefix"))
	return false
}

func managePost(w http.ResponseWriter, r *http.Request, ps httprouter.Params, del bool) {
	w.Header().Set("Content-Type", "application/json")

	mod, authed := auth(w, r)
	if !authed {
		return
	}

	var f postForm
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&f)
	if err != nil {
		sendError(w, http.StatusBadRequest, errors.New("bad form"))
		return
	}
	numId, err := strconv.Atoi(f.Id)
	if err != nil || numId < 1 {
		sendError(w, http.StatusBadRequest, errors.New("bad form"))
		return
	}

	post, err := client.GetPost(session, f.Id)
	if sendError(w, http.StatusInternalServerError, err) {
		return
	}
	if (del && post.IsDeleted) || (!del && !post.IsDeleted) || post.IsOpPost {
		sendError(w, http.StatusBadRequest, errors.New("bad post"))
		return
	}

	opPost, err := client.GetPost(session, post.OpPostId)
	if sendError(w, http.StatusInternalServerError, err) {
		return
	}
	if opPost.IsDeleted {
		sendError(w, http.StatusBadRequest, errors.New("bad thread"))
		return
	}

	if !ensurePrefix(w, mod, opPost) {
		return
	}

	if del {
		err = client.DeletePost(session, f.Id)
	} else {
		err = client.RestorePost(session, f.Id)
	}
	sendError(w, http.StatusInternalServerError, err)
}

func deletePost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	managePost(w, r, ps, true)
}

func restorePost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	managePost(w, r, ps, false)
}

func Serve(cfg common.Config) {
	var err error
	session, err = client.Auth(cfg.Auth.Username, cfg.Auth.Password)
	common.HandleError(err)

	for _, mod := range cfg.Mods {
		if mod.Token == "" || mod.Prefix == "" {
			log.Fatalln("Empty mod entry")
		}
		modByToken[mod.Token] = mod
	}

	// TODO(Kagami): Delete/restore thread.
	// TODO(Kagami): Set/unset NSFW.
	router := httprouter.New()
	router.POST("/api/post/delete", deletePost)
	router.POST("/api/post/restore", restorePost)
	addr := cfg.Serve.Host + ":" + strconv.Itoa(cfg.Serve.Port)

	log.Printf("Listening on %v", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
