// Code generated by counterfeiter. DO NOT EDIT.
package mixpanelfakes

import (
	"sync"

	mixpanela "github.com/dukex/mixpanel"
)

type FakeMixpanel struct {
	AliasStub        func(string, string) error
	aliasMutex       sync.RWMutex
	aliasArgsForCall []struct {
		arg1 string
		arg2 string
	}
	aliasReturns struct {
		result1 error
	}
	aliasReturnsOnCall map[int]struct {
		result1 error
	}
	TrackStub        func(string, string, *mixpanela.Event) error
	trackMutex       sync.RWMutex
	trackArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 *mixpanela.Event
	}
	trackReturns struct {
		result1 error
	}
	trackReturnsOnCall map[int]struct {
		result1 error
	}
	UpdateStub        func(string, *mixpanela.Update) error
	updateMutex       sync.RWMutex
	updateArgsForCall []struct {
		arg1 string
		arg2 *mixpanela.Update
	}
	updateReturns struct {
		result1 error
	}
	updateReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeMixpanel) Alias(arg1 string, arg2 string) error {
	fake.aliasMutex.Lock()
	ret, specificReturn := fake.aliasReturnsOnCall[len(fake.aliasArgsForCall)]
	fake.aliasArgsForCall = append(fake.aliasArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("Alias", []interface{}{arg1, arg2})
	fake.aliasMutex.Unlock()
	if fake.AliasStub != nil {
		return fake.AliasStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.aliasReturns
	return fakeReturns.result1
}

func (fake *FakeMixpanel) AliasCallCount() int {
	fake.aliasMutex.RLock()
	defer fake.aliasMutex.RUnlock()
	return len(fake.aliasArgsForCall)
}

func (fake *FakeMixpanel) AliasCalls(stub func(string, string) error) {
	fake.aliasMutex.Lock()
	defer fake.aliasMutex.Unlock()
	fake.AliasStub = stub
}

func (fake *FakeMixpanel) AliasArgsForCall(i int) (string, string) {
	fake.aliasMutex.RLock()
	defer fake.aliasMutex.RUnlock()
	argsForCall := fake.aliasArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeMixpanel) AliasReturns(result1 error) {
	fake.aliasMutex.Lock()
	defer fake.aliasMutex.Unlock()
	fake.AliasStub = nil
	fake.aliasReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeMixpanel) AliasReturnsOnCall(i int, result1 error) {
	fake.aliasMutex.Lock()
	defer fake.aliasMutex.Unlock()
	fake.AliasStub = nil
	if fake.aliasReturnsOnCall == nil {
		fake.aliasReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.aliasReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeMixpanel) Track(arg1 string, arg2 string, arg3 *mixpanela.Event) error {
	fake.trackMutex.Lock()
	ret, specificReturn := fake.trackReturnsOnCall[len(fake.trackArgsForCall)]
	fake.trackArgsForCall = append(fake.trackArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 *mixpanela.Event
	}{arg1, arg2, arg3})
	fake.recordInvocation("Track", []interface{}{arg1, arg2, arg3})
	fake.trackMutex.Unlock()
	if fake.TrackStub != nil {
		return fake.TrackStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.trackReturns
	return fakeReturns.result1
}

func (fake *FakeMixpanel) TrackCallCount() int {
	fake.trackMutex.RLock()
	defer fake.trackMutex.RUnlock()
	return len(fake.trackArgsForCall)
}

func (fake *FakeMixpanel) TrackCalls(stub func(string, string, *mixpanela.Event) error) {
	fake.trackMutex.Lock()
	defer fake.trackMutex.Unlock()
	fake.TrackStub = stub
}

func (fake *FakeMixpanel) TrackArgsForCall(i int) (string, string, *mixpanela.Event) {
	fake.trackMutex.RLock()
	defer fake.trackMutex.RUnlock()
	argsForCall := fake.trackArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeMixpanel) TrackReturns(result1 error) {
	fake.trackMutex.Lock()
	defer fake.trackMutex.Unlock()
	fake.TrackStub = nil
	fake.trackReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeMixpanel) TrackReturnsOnCall(i int, result1 error) {
	fake.trackMutex.Lock()
	defer fake.trackMutex.Unlock()
	fake.TrackStub = nil
	if fake.trackReturnsOnCall == nil {
		fake.trackReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.trackReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeMixpanel) Update(arg1 string, arg2 *mixpanela.Update) error {
	fake.updateMutex.Lock()
	ret, specificReturn := fake.updateReturnsOnCall[len(fake.updateArgsForCall)]
	fake.updateArgsForCall = append(fake.updateArgsForCall, struct {
		arg1 string
		arg2 *mixpanela.Update
	}{arg1, arg2})
	fake.recordInvocation("Update", []interface{}{arg1, arg2})
	fake.updateMutex.Unlock()
	if fake.UpdateStub != nil {
		return fake.UpdateStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.updateReturns
	return fakeReturns.result1
}

func (fake *FakeMixpanel) UpdateCallCount() int {
	fake.updateMutex.RLock()
	defer fake.updateMutex.RUnlock()
	return len(fake.updateArgsForCall)
}

func (fake *FakeMixpanel) UpdateCalls(stub func(string, *mixpanela.Update) error) {
	fake.updateMutex.Lock()
	defer fake.updateMutex.Unlock()
	fake.UpdateStub = stub
}

func (fake *FakeMixpanel) UpdateArgsForCall(i int) (string, *mixpanela.Update) {
	fake.updateMutex.RLock()
	defer fake.updateMutex.RUnlock()
	argsForCall := fake.updateArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeMixpanel) UpdateReturns(result1 error) {
	fake.updateMutex.Lock()
	defer fake.updateMutex.Unlock()
	fake.UpdateStub = nil
	fake.updateReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeMixpanel) UpdateReturnsOnCall(i int, result1 error) {
	fake.updateMutex.Lock()
	defer fake.updateMutex.Unlock()
	fake.UpdateStub = nil
	if fake.updateReturnsOnCall == nil {
		fake.updateReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.updateReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeMixpanel) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.aliasMutex.RLock()
	defer fake.aliasMutex.RUnlock()
	fake.trackMutex.RLock()
	defer fake.trackMutex.RUnlock()
	fake.updateMutex.RLock()
	defer fake.updateMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeMixpanel) recordInvocation(key string, args []interface{}) {
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

var _ mixpanela.Mixpanel = new(FakeMixpanel)
