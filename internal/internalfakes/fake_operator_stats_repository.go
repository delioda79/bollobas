// Code generated by counterfeiter. DO NOT EDIT.
package internalfakes

import (
	"context"
	"sync"

	"github.com/taxibeat/bollobas/internal"
)

type FakeOperatorStatsRepository struct {
	AddStub        func(context.Context, *internal.OperatorStats) error
	addMutex       sync.RWMutex
	addArgsForCall []struct {
		arg1 context.Context
		arg2 *internal.OperatorStats
	}
	addReturns struct {
		result1 error
	}
	addReturnsOnCall map[int]struct {
		result1 error
	}
	GetAllStub        func(context.Context, internal.DateFilter, internal.Pagination) ([]internal.OperatorStats, error)
	getAllMutex       sync.RWMutex
	getAllArgsForCall []struct {
		arg1 context.Context
		arg2 internal.DateFilter
		arg3 internal.Pagination
	}
	getAllReturns struct {
		result1 []internal.OperatorStats
		result2 error
	}
	getAllReturnsOnCall map[int]struct {
		result1 []internal.OperatorStats
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeOperatorStatsRepository) Add(arg1 context.Context, arg2 *internal.OperatorStats) error {
	fake.addMutex.Lock()
	ret, specificReturn := fake.addReturnsOnCall[len(fake.addArgsForCall)]
	fake.addArgsForCall = append(fake.addArgsForCall, struct {
		arg1 context.Context
		arg2 *internal.OperatorStats
	}{arg1, arg2})
	fake.recordInvocation("Add", []interface{}{arg1, arg2})
	fake.addMutex.Unlock()
	if fake.AddStub != nil {
		return fake.AddStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.addReturns
	return fakeReturns.result1
}

func (fake *FakeOperatorStatsRepository) AddCallCount() int {
	fake.addMutex.RLock()
	defer fake.addMutex.RUnlock()
	return len(fake.addArgsForCall)
}

func (fake *FakeOperatorStatsRepository) AddCalls(stub func(context.Context, *internal.OperatorStats) error) {
	fake.addMutex.Lock()
	defer fake.addMutex.Unlock()
	fake.AddStub = stub
}

func (fake *FakeOperatorStatsRepository) AddArgsForCall(i int) (context.Context, *internal.OperatorStats) {
	fake.addMutex.RLock()
	defer fake.addMutex.RUnlock()
	argsForCall := fake.addArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeOperatorStatsRepository) AddReturns(result1 error) {
	fake.addMutex.Lock()
	defer fake.addMutex.Unlock()
	fake.AddStub = nil
	fake.addReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeOperatorStatsRepository) AddReturnsOnCall(i int, result1 error) {
	fake.addMutex.Lock()
	defer fake.addMutex.Unlock()
	fake.AddStub = nil
	if fake.addReturnsOnCall == nil {
		fake.addReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.addReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeOperatorStatsRepository) GetAll(arg1 context.Context, arg2 internal.DateFilter, arg3 internal.Pagination) ([]internal.OperatorStats, error) {
	fake.getAllMutex.Lock()
	ret, specificReturn := fake.getAllReturnsOnCall[len(fake.getAllArgsForCall)]
	fake.getAllArgsForCall = append(fake.getAllArgsForCall, struct {
		arg1 context.Context
		arg2 internal.DateFilter
		arg3 internal.Pagination
	}{arg1, arg2, arg3})
	fake.recordInvocation("GetAll", []interface{}{arg1, arg2, arg3})
	fake.getAllMutex.Unlock()
	if fake.GetAllStub != nil {
		return fake.GetAllStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getAllReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeOperatorStatsRepository) GetAllCallCount() int {
	fake.getAllMutex.RLock()
	defer fake.getAllMutex.RUnlock()
	return len(fake.getAllArgsForCall)
}

func (fake *FakeOperatorStatsRepository) GetAllCalls(stub func(context.Context, internal.DateFilter, internal.Pagination) ([]internal.OperatorStats, error)) {
	fake.getAllMutex.Lock()
	defer fake.getAllMutex.Unlock()
	fake.GetAllStub = stub
}

func (fake *FakeOperatorStatsRepository) GetAllArgsForCall(i int) (context.Context, internal.DateFilter, internal.Pagination) {
	fake.getAllMutex.RLock()
	defer fake.getAllMutex.RUnlock()
	argsForCall := fake.getAllArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeOperatorStatsRepository) GetAllReturns(result1 []internal.OperatorStats, result2 error) {
	fake.getAllMutex.Lock()
	defer fake.getAllMutex.Unlock()
	fake.GetAllStub = nil
	fake.getAllReturns = struct {
		result1 []internal.OperatorStats
		result2 error
	}{result1, result2}
}

func (fake *FakeOperatorStatsRepository) GetAllReturnsOnCall(i int, result1 []internal.OperatorStats, result2 error) {
	fake.getAllMutex.Lock()
	defer fake.getAllMutex.Unlock()
	fake.GetAllStub = nil
	if fake.getAllReturnsOnCall == nil {
		fake.getAllReturnsOnCall = make(map[int]struct {
			result1 []internal.OperatorStats
			result2 error
		})
	}
	fake.getAllReturnsOnCall[i] = struct {
		result1 []internal.OperatorStats
		result2 error
	}{result1, result2}
}

func (fake *FakeOperatorStatsRepository) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.addMutex.RLock()
	defer fake.addMutex.RUnlock()
	fake.getAllMutex.RLock()
	defer fake.getAllMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeOperatorStatsRepository) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ internal.OperatorStatsRepository = new(FakeOperatorStatsRepository)
