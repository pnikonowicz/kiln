// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pivotal-cf/kiln/fetcher"
)

type S3ObjectLister struct {
	ListObjectsPagesStub        func(*s3.ListObjectsInput, func(*s3.ListObjectsOutput, bool) bool) error
	listObjectsPagesMutex       sync.RWMutex
	listObjectsPagesArgsForCall []struct {
		arg1 *s3.ListObjectsInput
		arg2 func(*s3.ListObjectsOutput, bool) bool
	}
	listObjectsPagesReturns struct {
		result1 error
	}
	listObjectsPagesReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *S3ObjectLister) ListObjectsPages(arg1 *s3.ListObjectsInput, arg2 func(*s3.ListObjectsOutput, bool) bool) error {
	fake.listObjectsPagesMutex.Lock()
	ret, specificReturn := fake.listObjectsPagesReturnsOnCall[len(fake.listObjectsPagesArgsForCall)]
	fake.listObjectsPagesArgsForCall = append(fake.listObjectsPagesArgsForCall, struct {
		arg1 *s3.ListObjectsInput
		arg2 func(*s3.ListObjectsOutput, bool) bool
	}{arg1, arg2})
	fake.recordInvocation("ListObjectsPages", []interface{}{arg1, arg2})
	fake.listObjectsPagesMutex.Unlock()
	if fake.ListObjectsPagesStub != nil {
		return fake.ListObjectsPagesStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.listObjectsPagesReturns
	return fakeReturns.result1
}

func (fake *S3ObjectLister) ListObjectsPagesCallCount() int {
	fake.listObjectsPagesMutex.RLock()
	defer fake.listObjectsPagesMutex.RUnlock()
	return len(fake.listObjectsPagesArgsForCall)
}

func (fake *S3ObjectLister) ListObjectsPagesCalls(stub func(*s3.ListObjectsInput, func(*s3.ListObjectsOutput, bool) bool) error) {
	fake.listObjectsPagesMutex.Lock()
	defer fake.listObjectsPagesMutex.Unlock()
	fake.ListObjectsPagesStub = stub
}

func (fake *S3ObjectLister) ListObjectsPagesArgsForCall(i int) (*s3.ListObjectsInput, func(*s3.ListObjectsOutput, bool) bool) {
	fake.listObjectsPagesMutex.RLock()
	defer fake.listObjectsPagesMutex.RUnlock()
	argsForCall := fake.listObjectsPagesArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *S3ObjectLister) ListObjectsPagesReturns(result1 error) {
	fake.listObjectsPagesMutex.Lock()
	defer fake.listObjectsPagesMutex.Unlock()
	fake.ListObjectsPagesStub = nil
	fake.listObjectsPagesReturns = struct {
		result1 error
	}{result1}
}

func (fake *S3ObjectLister) ListObjectsPagesReturnsOnCall(i int, result1 error) {
	fake.listObjectsPagesMutex.Lock()
	defer fake.listObjectsPagesMutex.Unlock()
	fake.ListObjectsPagesStub = nil
	if fake.listObjectsPagesReturnsOnCall == nil {
		fake.listObjectsPagesReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.listObjectsPagesReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *S3ObjectLister) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.listObjectsPagesMutex.RLock()
	defer fake.listObjectsPagesMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *S3ObjectLister) recordInvocation(key string, args []interface{}) {
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

var _ fetcher.S3ObjectLister = new(S3ObjectLister)
