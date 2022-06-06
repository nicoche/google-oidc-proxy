package handler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"google.golang.org/api/idtoken"
)

type Handler struct {
	client *http.Client
	target string
}

func NewHandler(ctx context.Context, target string, audience string) (*Handler, error) {
	client, err := idtoken.NewClient(ctx, audience)
	if err != nil {
		return nil, err
	}

	handler := &Handler{
		client: client,
		target: target,
	}
	return handler, nil
}

func (h *Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	resp, err := h.doProxy(req)
	if err != nil {
		rw.WriteHeader(http.StatusBadGateway)
		rw.Write([]byte(fmt.Sprintf("proxy request failed: %+v", err)))
		log.Printf("failed to proxy request: %+v\n", err)
		return
	}

	defer resp.Body.Close()

	buf := bytes.Buffer{}
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(fmt.Sprintf("proxy response failed: %+v", err)))
		log.Printf("failed to proxy response: %+v\n", err)
		return
	}

	for name, values := range resp.Header {
		for _, item := range values {
			rw.Header().Add(name, item)
		}
	}
	rw.WriteHeader(resp.StatusCode)
	rw.Write(buf.Bytes())
}

func (h *Handler) doProxy(req *http.Request) (*http.Response, error) {
	url := *req.URL
	url.Scheme = "https"
	url.Host = h.target

	proxyReq, err := http.NewRequest(req.Method, url.String(), req.Body)
	if err != nil {
		return nil, err
	}

	proxyReq.Header = req.Header.Clone()

	resp, err := h.client.Do(proxyReq)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
