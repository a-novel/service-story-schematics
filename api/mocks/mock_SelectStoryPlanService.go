// Code generated by mockery v2.53.1. DO NOT EDIT.

package apimocks

import (
	context "github.com/a-novel-kit/context"
	mock "github.com/stretchr/testify/mock"

	models "github.com/a-novel/service-story-schematics/models"

	services "github.com/a-novel/service-story-schematics/internal/services"
)

// MockSelectStoryPlanService is an autogenerated mock type for the SelectStoryPlanService type
type MockSelectStoryPlanService struct {
	mock.Mock
}

type MockSelectStoryPlanService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockSelectStoryPlanService) EXPECT() *MockSelectStoryPlanService_Expecter {
	return &MockSelectStoryPlanService_Expecter{mock: &_m.Mock}
}

// SelectStoryPlan provides a mock function with given fields: ctx, request
func (_m *MockSelectStoryPlanService) SelectStoryPlan(ctx context.Context, request services.SelectStoryPlanRequest) (*models.StoryPlan, error) {
	ret := _m.Called(ctx, request)

	if len(ret) == 0 {
		panic("no return value specified for SelectStoryPlan")
	}

	var r0 *models.StoryPlan
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, services.SelectStoryPlanRequest) (*models.StoryPlan, error)); ok {
		return rf(ctx, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, services.SelectStoryPlanRequest) *models.StoryPlan); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.StoryPlan)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, services.SelectStoryPlanRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockSelectStoryPlanService_SelectStoryPlan_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SelectStoryPlan'
type MockSelectStoryPlanService_SelectStoryPlan_Call struct {
	*mock.Call
}

// SelectStoryPlan is a helper method to define mock.On call
//   - ctx context.Context
//   - request services.SelectStoryPlanRequest
func (_e *MockSelectStoryPlanService_Expecter) SelectStoryPlan(ctx interface{}, request interface{}) *MockSelectStoryPlanService_SelectStoryPlan_Call {
	return &MockSelectStoryPlanService_SelectStoryPlan_Call{Call: _e.mock.On("SelectStoryPlan", ctx, request)}
}

func (_c *MockSelectStoryPlanService_SelectStoryPlan_Call) Run(run func(ctx context.Context, request services.SelectStoryPlanRequest)) *MockSelectStoryPlanService_SelectStoryPlan_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(services.SelectStoryPlanRequest))
	})
	return _c
}

func (_c *MockSelectStoryPlanService_SelectStoryPlan_Call) Return(_a0 *models.StoryPlan, _a1 error) *MockSelectStoryPlanService_SelectStoryPlan_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockSelectStoryPlanService_SelectStoryPlan_Call) RunAndReturn(run func(context.Context, services.SelectStoryPlanRequest) (*models.StoryPlan, error)) *MockSelectStoryPlanService_SelectStoryPlan_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockSelectStoryPlanService creates a new instance of MockSelectStoryPlanService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockSelectStoryPlanService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockSelectStoryPlanService {
	mock := &MockSelectStoryPlanService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
