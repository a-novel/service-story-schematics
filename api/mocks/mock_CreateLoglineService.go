// Code generated by mockery v2.53.1. DO NOT EDIT.

package apimocks

import (
	context "github.com/a-novel-kit/context"
	mock "github.com/stretchr/testify/mock"

	models "github.com/a-novel/service-story-schematics/models"

	services "github.com/a-novel/service-story-schematics/internal/services"
)

// MockCreateLoglineService is an autogenerated mock type for the CreateLoglineService type
type MockCreateLoglineService struct {
	mock.Mock
}

type MockCreateLoglineService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockCreateLoglineService) EXPECT() *MockCreateLoglineService_Expecter {
	return &MockCreateLoglineService_Expecter{mock: &_m.Mock}
}

// CreateLogline provides a mock function with given fields: ctx, request
func (_m *MockCreateLoglineService) CreateLogline(ctx context.Context, request services.CreateLoglineRequest) (*models.Logline, error) {
	ret := _m.Called(ctx, request)

	if len(ret) == 0 {
		panic("no return value specified for CreateLogline")
	}

	var r0 *models.Logline
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, services.CreateLoglineRequest) (*models.Logline, error)); ok {
		return rf(ctx, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, services.CreateLoglineRequest) *models.Logline); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Logline)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, services.CreateLoglineRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCreateLoglineService_CreateLogline_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateLogline'
type MockCreateLoglineService_CreateLogline_Call struct {
	*mock.Call
}

// CreateLogline is a helper method to define mock.On call
//   - ctx context.Context
//   - request services.CreateLoglineRequest
func (_e *MockCreateLoglineService_Expecter) CreateLogline(ctx interface{}, request interface{}) *MockCreateLoglineService_CreateLogline_Call {
	return &MockCreateLoglineService_CreateLogline_Call{Call: _e.mock.On("CreateLogline", ctx, request)}
}

func (_c *MockCreateLoglineService_CreateLogline_Call) Run(run func(ctx context.Context, request services.CreateLoglineRequest)) *MockCreateLoglineService_CreateLogline_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(services.CreateLoglineRequest))
	})
	return _c
}

func (_c *MockCreateLoglineService_CreateLogline_Call) Return(_a0 *models.Logline, _a1 error) *MockCreateLoglineService_CreateLogline_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCreateLoglineService_CreateLogline_Call) RunAndReturn(run func(context.Context, services.CreateLoglineRequest) (*models.Logline, error)) *MockCreateLoglineService_CreateLogline_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockCreateLoglineService creates a new instance of MockCreateLoglineService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockCreateLoglineService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockCreateLoglineService {
	mock := &MockCreateLoglineService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
