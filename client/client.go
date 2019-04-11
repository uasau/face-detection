package client

import (
	"encoding/json"
	"image"
	"io"
	"net"
	"net/http"
	"time"
)

// Client interface for interaction with the service
type Client interface {
	DetectFaces(r io.Reader)
}

// FaceDetect is a http client for the face detection service
type FaceDetect struct {
	httpClient *http.Client
	location   string
}

// Response is returned by a successful face detecion call
type Response struct {
	Faces  []image.Rectangle
	Bounds image.Rectangle
}

// NewClient creates a new FaceDetect
func NewClient(l string) *FaceDetect {
	hc := http.DefaultClient
	hc.Timeout = 60 * time.Second
	// force transport to use ipv4
	hc.Transport = &http.Transport{
		Dial: (func(network, addr string) (net.Conn, error) {
			return (&net.Dialer{
				Timeout:   3 * time.Second,
				LocalAddr: nil,
				DualStack: false,
			}).Dial("tcp4", addr)
		}),
	}

	return &FaceDetect{hc, l}
}

// DetectFaces sends a request to the face detection service
func (h *FaceDetect) DetectFaces(r io.Reader) (*Response, error) {
	req, err := http.NewRequest("POST", h.location, r)
	if err != nil {
		return nil, err
	}

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	faces := &Response{}
	defer resp.Body.Close()
	d := json.NewDecoder(resp.Body)

	err = d.Decode(faces)
	if err != nil {
		return nil, err
	}

	return faces, nil
}
