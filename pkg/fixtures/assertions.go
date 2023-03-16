package fixtures

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/containernetworking/cni/pkg/types"
	currentcni "github.com/containernetworking/cni/pkg/types/040"
	"github.com/jboelensns/openstack-cni/pkg/cniserver"
	. "github.com/pepinns/go-hamcrest"
)

type Assertions struct {
	t *testing.T
}

func (me *Assertions) IsCniError(resp *http.Response, err error) *types.Error {
	return me.IsCniErrorWithCode(resp, err, 500)
}

func (me *Assertions) IsCniErrorWithCode(resp *http.Response, err error, httpCode int) *types.Error {
	Assert(me.t).That(resp.StatusCode, Equals(httpCode))
	result := Is[types.Error](me.t, resp, err)
	return result
}

func (me *Assertions) IsCniResult(resp *http.Response, err error) *currentcni.Result {
	Assert(me.t).That(resp.StatusCode, Equals(200))
	result := Is[currentcni.Result](me.t, resp, err)
	return result
}

func (me *Assertions) IsGoodHealthResult(resp *http.Response, err error) *cniserver.HealthResponse {
	Assert(me.t).That(resp.StatusCode, Equals(200))
	result := Is[cniserver.HealthResponse](me.t, resp, err)
	Assert(me.t).That(result.IsHealthy, IsTrue())
	return result
}

func (me *Assertions) IsSickHealthResult(resp *http.Response, err error) *cniserver.HealthResponse {
	Assert(me.t).That(resp.StatusCode, Equals(500))
	result := Is[cniserver.HealthResponse](me.t, resp, err)
	Assert(me.t).That(result.IsHealthy, IsFalse())
	return result
}

func Is[T any](t *testing.T, resp *http.Response, err error) *T {
	t.Helper()

	Assert(t).That(err, IsNil())

	b, err := io.ReadAll(resp.Body)
	Assert(t).That(err, IsNil())

	var data = new(T)
	if err := json.Unmarshal(b, &data); err != nil {
		Assert(t).That(err, IsNil())
	}
	return data
}

func (me *Assertions) CniErrorHasCode(resp *http.Response, err error, code uint) {
	me.t.Helper()
	cerr := me.IsCniError(resp, err)
	Assert(me.t).That(cerr.Code, Equals(types.ErrInternal))
}
