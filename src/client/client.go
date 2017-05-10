package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

const apiPrefix string = "https://0chan.hk/api"

type apiReponse struct {
	Ok bool
}

type loginForm struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type postOuter struct {
	Post Post
}

type Post struct {
	OpPostId  string
	IsOpPost  bool
	IsDeleted bool
	Message   string
}

func Auth(username string, password string) (session string, err error) {
	url := apiPrefix + "/user/login"
	form := loginForm{username, password}
	data, _ := json.Marshal(form)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil || resp.StatusCode != 200 {
		return "", errors.New("cannot login")
	}
	defer resp.Body.Close()
	// TODO(Kagami): Use json.Decoder.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("cannot read response")
	}
	var info apiReponse
	err = json.Unmarshal(body, &info)
	if err != nil {
		return "", errors.New("cannot parse response")
	}
	if !info.Ok {
		return "", errors.New("cannot login")
	}
	return resp.Header["X-Session"][0], nil
}

func GetPost(session string, id string) (p Post, err error) {
	v := url.Values{}
	v.Set("session", session)
	v.Set("post", id)
	url := apiPrefix + "/post?" + v.Encode()
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		return p, errors.New("cannot get post")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return p, errors.New("cannot read response")
	}
	var po postOuter
	err = json.Unmarshal(body, &po)
	if err != nil {
		return p, errors.New("cannot parse response")
	}
	return po.Post, nil
}

func DeletePost(session string, id string) error {
	v := url.Values{}
	v.Set("session", session)
	v.Set("post", id)
	url := apiPrefix + "/moderation/deletePost?" + v.Encode()
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		return errors.New("cannot delete post")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("cannot read response")
	}
	var info apiReponse
	err = json.Unmarshal(body, &info)
	if err != nil {
		return errors.New("cannot parse response")
	}
	if !info.Ok {
		return errors.New("cannot delete post")
	}
	return nil
}

func RestorePost(session string, id string) error {
	v := url.Values{}
	v.Set("session", session)
	v.Set("post", id)
	url := apiPrefix + "/moderation/restorePost?" + v.Encode()
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		return errors.New("cannot restore post")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("cannot read response")
	}
	var info apiReponse
	err = json.Unmarshal(body, &info)
	if err != nil {
		return errors.New("cannot parse response")
	}
	if !info.Ok {
		return errors.New("cannot restore post")
	}
	return nil
}
