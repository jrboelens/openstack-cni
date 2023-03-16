package cniserver

import (
	"net/http"

	"github.com/jboelensns/openstack-cni/pkg/openstack"
)

type HealthHandler struct {
	osClient openstack.OpenstackClient
}

func NewHealthHandler(osApi openstack.OpenstackClient) *HealthHandler {
	return &HealthHandler{osApi}
}

func (me *HealthHandler) HandleRequest(w http.ResponseWriter, r *http.Request) {
	health := HealthResponse{
		IsHealthy: true,
		Checks: []HealthResponseCheck{
			me.checkOpenstack(),
		},
	}

	for _, check := range health.Checks {
		if !check.IsHealthy {
			health.IsHealthy = false
			break
		}

	}
	if !health.IsHealthy {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(asJson(health))
}

func (me HealthHandler) checkOpenstack() HealthResponseCheck {
	resp := HealthResponseCheck{
		Name:      "openstack",
		IsHealthy: true,
		Error:     "",
	}
	_, err := me.osClient.GetServerByName("serverthatdoesntexist")
	if err == openstack.ErrServerNotFound {
		return resp
	}
	resp.IsHealthy = false
	resp.Error = err.Error()

	return resp
}

type HealthResponse struct {
	IsHealthy bool                  `json:"is_healthy,omitempty"`
	Checks    []HealthResponseCheck `json:"checks,omitempty"`
}

type HealthResponseCheck struct {
	Name      string `json:"name,omitempty"`
	IsHealthy bool   `json:"is_healthy,omitempty"`
	Error     string `json:"error,omitempty"`
}
