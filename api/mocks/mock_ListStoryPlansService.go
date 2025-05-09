// Code generated by mockery v2.53.1. DO NOT EDIT.

package apimocks

import (
	context "github.com/a-novel-kit/context"
	mock "github.com/stretchr/testify/mock"

	models "github.com/a-novel/service-story-schematics/models"

	services "github.com/a-novel/service-story-schematics/internal/services"
)

// MockListStoryPlansService is an autogenerated mock type for the ListStoryPlansService type
type MockListStoryPlansService struct {
	mock.Mock
}

type MockListStoryPlansService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockListStoryPlansService) EXPECT() *MockListStoryPlansService_Expecter {
	return &MockListStoryPlansService_Expecter{mock: &_m.Mock}
}

// ListStoryPlans provides a mock function with given fields: ctx, request
func (_m *MockListStoryPlansService) ListStoryPlans(ctx context.Context, request services.ListStoryPlansRequest) ([]*models.StoryPlanPreview, error) {
	ret := _m.Called(ctx, request)

	if len(ret) == 0 {
		panic("no return value specified for ListStoryPlans")
	}

	var r0 []*models.StoryPlanPreview
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, services.ListStoryPlansRequest) ([]*models.StoryPlanPreview, error)); ok {
		return rf(ctx, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, services.ListStoryPlansRequest) []*models.StoryPlanPreview); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.StoryPlanPreview)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, services.ListStoryPlansRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockListStoryPlansService_ListStoryPlans_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListStoryPlans'
type MockListStoryPlansService_ListStoryPlans_Call struct {
	*mock.Call
}

// ListStoryPlans is a helper method to define mock.On call
//   - ctx context.Context
//   - request services.ListStoryPlansRequest
func (_e *MockListStoryPlansService_Expecter) ListStoryPlans(ctx interface{}, request interface{}) *MockListStoryPlansService_ListStoryPlans_Call {
	return &MockListStoryPlansService_ListStoryPlans_Call{Call: _e.mock.On("ListStoryPlans", ctx, request)}
}

func (_c *MockListStoryPlansService_ListStoryPlans_Call) Run(run func(ctx context.Context, request services.ListStoryPlansRequest)) *MockListStoryPlansService_ListStoryPlans_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(services.ListStoryPlansRequest))
	})
	return _c
}

func (_c *MockListStoryPlansService_ListStoryPlans_Call) Return(_a0 []*models.StoryPlanPreview, _a1 error) *MockListStoryPlansService_ListStoryPlans_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockListStoryPlansService_ListStoryPlans_Call) RunAndReturn(run func(context.Context, services.ListStoryPlansRequest) ([]*models.StoryPlanPreview, error)) *MockListStoryPlansService_ListStoryPlans_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockListStoryPlansService creates a new instance of MockListStoryPlansService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockListStoryPlansService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockListStoryPlansService {
	mock := &MockListStoryPlansService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
