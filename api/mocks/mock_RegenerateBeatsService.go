// Code generated by mockery v2.53.1. DO NOT EDIT.

package apimocks

import (
	context "github.com/a-novel-kit/context"
	mock "github.com/stretchr/testify/mock"

	models "github.com/a-novel/service-story-schematics/models"

	services "github.com/a-novel/service-story-schematics/internal/services"
)

// MockRegenerateBeatsService is an autogenerated mock type for the RegenerateBeatsService type
type MockRegenerateBeatsService struct {
	mock.Mock
}

type MockRegenerateBeatsService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockRegenerateBeatsService) EXPECT() *MockRegenerateBeatsService_Expecter {
	return &MockRegenerateBeatsService_Expecter{mock: &_m.Mock}
}

// RegenerateBeats provides a mock function with given fields: ctx, request
func (_m *MockRegenerateBeatsService) RegenerateBeats(ctx context.Context, request services.RegenerateBeatsRequest) ([]models.Beat, error) {
	ret := _m.Called(ctx, request)

	if len(ret) == 0 {
		panic("no return value specified for RegenerateBeats")
	}

	var r0 []models.Beat
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, services.RegenerateBeatsRequest) ([]models.Beat, error)); ok {
		return rf(ctx, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, services.RegenerateBeatsRequest) []models.Beat); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Beat)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, services.RegenerateBeatsRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRegenerateBeatsService_RegenerateBeats_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RegenerateBeats'
type MockRegenerateBeatsService_RegenerateBeats_Call struct {
	*mock.Call
}

// RegenerateBeats is a helper method to define mock.On call
//   - ctx context.Context
//   - request services.RegenerateBeatsRequest
func (_e *MockRegenerateBeatsService_Expecter) RegenerateBeats(ctx interface{}, request interface{}) *MockRegenerateBeatsService_RegenerateBeats_Call {
	return &MockRegenerateBeatsService_RegenerateBeats_Call{Call: _e.mock.On("RegenerateBeats", ctx, request)}
}

func (_c *MockRegenerateBeatsService_RegenerateBeats_Call) Run(run func(ctx context.Context, request services.RegenerateBeatsRequest)) *MockRegenerateBeatsService_RegenerateBeats_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(services.RegenerateBeatsRequest))
	})
	return _c
}

func (_c *MockRegenerateBeatsService_RegenerateBeats_Call) Return(_a0 []models.Beat, _a1 error) *MockRegenerateBeatsService_RegenerateBeats_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRegenerateBeatsService_RegenerateBeats_Call) RunAndReturn(run func(context.Context, services.RegenerateBeatsRequest) ([]models.Beat, error)) *MockRegenerateBeatsService_RegenerateBeats_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockRegenerateBeatsService creates a new instance of MockRegenerateBeatsService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockRegenerateBeatsService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockRegenerateBeatsService {
	mock := &MockRegenerateBeatsService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
