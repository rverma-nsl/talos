// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

//nolint:dupl
package network_test

import (
	"context"
	"log"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/cosi-project/runtime/pkg/controller/runtime"
	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/cosi-project/runtime/pkg/state"
	"github.com/cosi-project/runtime/pkg/state/impl/inmem"
	"github.com/cosi-project/runtime/pkg/state/impl/namespaced"
	"github.com/stretchr/testify/suite"
	"github.com/talos-systems/go-retry/retry"

	netctrl "github.com/talos-systems/talos/internal/app/machined/pkg/controllers/network"
	"github.com/talos-systems/talos/pkg/logging"
	"github.com/talos-systems/talos/pkg/machinery/constants"
	"github.com/talos-systems/talos/pkg/machinery/resources/network"
)

type TimeServerSpecSuite struct {
	suite.Suite

	state state.State

	runtime *runtime.Runtime
	wg      sync.WaitGroup

	ctx       context.Context
	ctxCancel context.CancelFunc
}

func (suite *TimeServerSpecSuite) SetupTest() {
	suite.ctx, suite.ctxCancel = context.WithTimeout(context.Background(), 3*time.Minute)

	suite.state = state.WrapCore(namespaced.NewState(inmem.Build))

	var err error

	suite.runtime, err = runtime.NewRuntime(suite.state, logging.Wrap(log.Writer()))
	suite.Require().NoError(err)

	suite.Require().NoError(suite.runtime.RegisterController(&netctrl.TimeServerSpecController{}))

	suite.startRuntime()
}

func (suite *TimeServerSpecSuite) startRuntime() {
	suite.wg.Add(1)

	go func() {
		defer suite.wg.Done()

		suite.Assert().NoError(suite.runtime.Run(suite.ctx))
	}()
}

func (suite *TimeServerSpecSuite) assertStatus(id string, servers ...string) error {
	r, err := suite.state.Get(suite.ctx, resource.NewMetadata(network.NamespaceName, network.TimeServerStatusType, id, resource.VersionUndefined))
	if err != nil {
		if state.IsNotFoundError(err) {
			return retry.ExpectedError(err)
		}

		return err
	}

	status := r.(*network.TimeServerStatus) //nolint:errcheck,forcetypeassert

	if !reflect.DeepEqual(status.TypedSpec().NTPServers, servers) {
		return retry.ExpectedErrorf("server list mismatch: %q != %q", status.TypedSpec().NTPServers, servers)
	}

	return nil
}

func (suite *TimeServerSpecSuite) TestSpec() {
	spec := network.NewTimeServerSpec(network.NamespaceName, "timeservers")
	*spec.TypedSpec() = network.TimeServerSpecSpec{
		NTPServers:  []string{constants.DefaultNTPServer},
		ConfigLayer: network.ConfigDefault,
	}

	for _, res := range []resource.Resource{spec} {
		suite.Require().NoError(suite.state.Create(suite.ctx, res), "%v", res.Spec())
	}

	suite.Assert().NoError(retry.Constant(3*time.Second, retry.WithUnits(100*time.Millisecond)).Retry(
		func() error {
			return suite.assertStatus("timeservers", constants.DefaultNTPServer)
		}))
}

func (suite *TimeServerSpecSuite) TearDownTest() {
	suite.T().Log("tear down")

	suite.ctxCancel()

	suite.wg.Wait()

	// trigger updates in resources to stop watch loops
	suite.Assert().NoError(suite.state.Create(context.Background(), network.NewTimeServerSpec(network.NamespaceName, "bar")))
}

func TestTimeServerSpecSuite(t *testing.T) {
	suite.Run(t, new(TimeServerSpecSuite))
}
