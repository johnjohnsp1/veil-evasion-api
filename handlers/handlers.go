package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"path"

	"github.com/mholt/binding"
	"github.com/tomsteele/veil-evasion-api/veil"
	"github.com/unrolled/render"
)

type H struct {
	C   *rpc.Client
	Raw *net.Conn
	R   *render.Render
}

func New(h *H) *H {
	h.R = render.New()
	return h
}

type Errors struct {
	Errors []map[string]string `json:"errors"`
}

func (h *H) JSON(w http.ResponseWriter, v interface{}) {
	h.R.JSON(w, http.StatusOK, v)
}

func (h *H) JSON500(w http.ResponseWriter, err error) {
	h.R.JSON(w, http.StatusInternalServerError, Errors{Errors: []map[string]string{map[string]string{"system error": err.Error()}}})
}

func (h *H) JSON400(w http.ResponseWriter, err error) {
	h.R.JSON(w, http.StatusBadRequest, Errors{Errors: []map[string]string{map[string]string{"invalid_value": err.Error()}}})
}

func (h *H) Version(w http.ResponseWriter, req *http.Request) {
	var v veil.Version
	err := h.C.Call("version", []string{}, &v)
	if err != nil {
		h.JSON500(w, err)
		return
	}
	h.JSON(w, map[string]veil.Version{"version": v})
}

func (h *H) Payloads(w http.ResponseWriter, req *http.Request) {
	var payloads veil.Payloads
	err := h.C.Call("payloads", []string{}, &payloads)
	if err != nil {
		h.JSON500(w, err)
		return
	}
	h.JSON(w, payloads)
}

type PayloadOption struct {
	Key          string `json:"key"`
	DefaultValue string `json:"default_value"`
	Value        string `json:"value"`
}

func (h *H) PayloadOptions(w http.ResponseWriter, req *http.Request) {
	p := req.URL.Query().Get("payload")
	if p == "" {
		h.JSON400(w, errors.New("payload query parameter missing"))
		return
	}
	var options veil.PayloadOptions
	var thing json.RawMessage
	err := h.C.Call("payload_options", []string{p}, &thing)
	if err != nil {
		h.JSON400(w, err)
		return
	}
	d, err := thing.MarshalJSON()
	if err != nil {
		h.JSON400(w, err)
	}
	if err := json.Unmarshal(d, &options); err != nil {
		h.JSON400(w, err)
		return
	}
	opts := []PayloadOption{}
	for _, o := range options {
		if len(o) < 3 {
			continue
		}
		opts = append(opts, PayloadOption{Key: o[0], DefaultValue: o[1], Value: o[2]})
	}
	opts = append(opts, PayloadOption{
		Key:          "pwnstaller",
		DefaultValue: "N",
		Value:        "Use pwnstaller",
	})
	opts = append(opts, PayloadOption{
		Key:          "outputbase",
		DefaultValue: "",
		Value:        "Output base for generated payloads",
	})
	opts = append(opts, PayloadOption{
		Key:          "overwrite",
		DefaultValue: "true",
		Value:        "Overwrite existing files",
	})
	h.JSON(w, opts)
}

type generateReq struct {
	Options []PayloadOption `"options"`
}

func (g *generateReq) FieldMap() binding.FieldMap {
	return binding.FieldMap{}
}

func (h *H) Generate(w http.ResponseWriter, req *http.Request) {
	opts := &generateReq{}
	if errs := binding.Bind(req, opts); errs.Handle(w) {
		return
	}
	vopts := []string{}
	for _, o := range opts.Options {
		vopts = append(vopts, fmt.Sprintf("%s=%s", o.Key, o.Value))
	}
	var paths string
	err := h.C.Call("generate", vopts, &paths)
	if err != nil {
		h.JSON400(w, errors.New("could not parse the response from veil, likely no options or an invalid payload"))
		return
	}
	if paths == "" {
		h.JSON400(w, errors.New("no payload was generated"))
		return
	}
	fullDir := path.Dir(paths)
	rootDir := path.Base(fullDir)
	f := path.Base(paths)
	httpath := path.Join(rootDir, f)
	h.JSON(w, map[string]string{"result": httpath})
}
