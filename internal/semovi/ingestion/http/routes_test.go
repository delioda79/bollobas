package http

import (
	"context"
	"testing"

	phttp "github.com/beatlabs/patron/component/http"
	"github.com/stretchr/testify/assert"
)

func defaultInput() (context.Context, *phttp.Request) {
	return context.Background(), &phttp.Request{}
}

func TestGetAggregatedRides(t *testing.T) {
	ctx, req := defaultInput()

	rsp, err := getAggregatedRides(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, rsp)
	assert.EqualValues(t, phttp.Response{}, *rsp)
}

func TestGetOperatorStats(t *testing.T) {
	ctx, req := defaultInput()

	rsp, err := getOperatorStats(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, rsp)
	assert.EqualValues(t, phttp.Response{}, *rsp)
}

func TestGetTransitsMade(t *testing.T) {
	ctx, req := defaultInput()

	rsp, err := getTransitsMade(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, rsp)
	assert.EqualValues(t, phttp.Response{}, *rsp)
}
