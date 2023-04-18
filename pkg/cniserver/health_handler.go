package cniserver

import (
	"net/http"

	"github.com/jboelensns/openstack-cni/pkg/openstack"
)

// HealthHandler handles all /health related requests
type HealthHandler struct {
	OsClient openstack.OpenstackClient
}

// HandleRequest executes health checks and returns the results
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
	_, err := me.OsClient.GetServerByName("serverthatdoesntexist")
	if err == openstack.ErrServerNotFound {
		return resp
	}
	resp.IsHealthy = false
	resp.Error = err.Error()

	return resp
}

// HealthResponse is returned for GET /health
type HealthResponse struct {
	IsHealthy bool                  `json:"is_healthy,omitempty"`
	Checks    []HealthResponseCheck `json:"checks,omitempty"`
}

// HealthResponse represents an invididual health check
type HealthResponseCheck struct {
	Name      string `json:"name,omitempty"`
	IsHealthy bool   `json:"is_healthy,omitempty"`
	Error     string `json:"error,omitempty"`
}
