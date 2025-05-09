// Code generated by mockery v2.53.1. DO NOT EDIT.

package servicesmocks

import (
	context "github.com/a-novel-kit/context"
	daoai "github.com/a-novel/service-story-schematics/internal/daoai"
	mock "github.com/stretchr/testify/mock"

	models "github.com/a-novel/service-story-schematics/models"
)

// MockGenerateLoglinesSource is an autogenerated mock type for the GenerateLoglinesSource type
type MockGenerateLoglinesSource struct {
	mock.Mock
}

type MockGenerateLoglinesSource_Expecter struct {
	mock *mock.Mock
}

func (_m *MockGenerateLoglinesSource) EXPECT() *MockGenerateLoglinesSource_Expecter {
	return &MockGenerateLoglinesSource_Expecter{mock: &_m.Mock}
}

// GenerateLoglines provides a mock function with given fields: ctx, request
func (_m *MockGenerateLoglinesSource) GenerateLoglines(ctx context.Context, request daoai.GenerateLoglinesRequest) ([]models.LoglineIdea, error) {
	ret := _m.Called(ctx, request)

	if len(ret) == 0 {
		panic("no return value specified for GenerateLoglines")
	}

	var r0 []models.LoglineIdea
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, daoai.GenerateLoglinesRequest) ([]models.LoglineIdea, error)); ok {
		return rf(ctx, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, daoai.GenerateLoglinesRequest) []models.LoglineIdea); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.LoglineIdea)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, daoai.GenerateLoglinesRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockGenerateLoglinesSource_GenerateLoglines_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GenerateLoglines'
type MockGenerateLoglinesSource_GenerateLoglines_Call struct {
	*mock.Call
}

// GenerateLoglines is a helper method to define mock.On call
//   - ctx context.Context
//   - request daoai.GenerateLoglinesRequest
func (_e *MockGenerateLoglinesSource_Expecter) GenerateLoglines(ctx interface{}, request interface{}) *MockGenerateLoglinesSource_GenerateLoglines_Call {
	return &MockGenerateLoglinesSource_GenerateLoglines_Call{Call: _e.mock.On("GenerateLoglines", ctx, request)}
}

func (_c *MockGenerateLoglinesSource_GenerateLoglines_Call) Run(run func(ctx context.Context, request daoai.GenerateLoglinesRequest)) *MockGenerateLoglinesSource_GenerateLoglines_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(daoai.GenerateLoglinesRequest))
	})
	return _c
}

func (_c *MockGenerateLoglinesSource_GenerateLoglines_Call) Return(_a0 []models.LoglineIdea, _a1 error) *MockGenerateLoglinesSource_GenerateLoglines_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockGenerateLoglinesSource_GenerateLoglines_Call) RunAndReturn(run func(context.Context, daoai.GenerateLoglinesRequest) ([]models.LoglineIdea, error)) *MockGenerateLoglinesSource_GenerateLoglines_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockGenerateLoglinesSource creates a new instance of MockGenerateLoglinesSource. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockGenerateLoglinesSource(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockGenerateLoglinesSource {
	mock := &MockGenerateLoglinesSource{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
