// Code generated by mockery v2.53.1. DO NOT EDIT.

package apimocks

import (
	context "github.com/a-novel-kit/context"
	mock "github.com/stretchr/testify/mock"

	models "github.com/a-novel/service-story-schematics/models"

	services "github.com/a-novel/service-story-schematics/internal/services"
)

// MockUpdateStoryPlanService is an autogenerated mock type for the UpdateStoryPlanService type
type MockUpdateStoryPlanService struct {
	mock.Mock
}

type MockUpdateStoryPlanService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockUpdateStoryPlanService) EXPECT() *MockUpdateStoryPlanService_Expecter {
	return &MockUpdateStoryPlanService_Expecter{mock: &_m.Mock}
}

// UpdateStoryPlan provides a mock function with given fields: ctx, request
func (_m *MockUpdateStoryPlanService) UpdateStoryPlan(ctx context.Context, request services.UpdateStoryPlanRequest) (*models.StoryPlan, error) {
	ret := _m.Called(ctx, request)

	if len(ret) == 0 {
		panic("no return value specified for UpdateStoryPlan")
	}

	var r0 *models.StoryPlan
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, services.UpdateStoryPlanRequest) (*models.StoryPlan, error)); ok {
		return rf(ctx, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, services.UpdateStoryPlanRequest) *models.StoryPlan); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.StoryPlan)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, services.UpdateStoryPlanRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockUpdateStoryPlanService_UpdateStoryPlan_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateStoryPlan'
type MockUpdateStoryPlanService_UpdateStoryPlan_Call struct {
	*mock.Call
}

// UpdateStoryPlan is a helper method to define mock.On call
//   - ctx context.Context
//   - request services.UpdateStoryPlanRequest
func (_e *MockUpdateStoryPlanService_Expecter) UpdateStoryPlan(ctx interface{}, request interface{}) *MockUpdateStoryPlanService_UpdateStoryPlan_Call {
	return &MockUpdateStoryPlanService_UpdateStoryPlan_Call{Call: _e.mock.On("UpdateStoryPlan", ctx, request)}
}

func (_c *MockUpdateStoryPlanService_UpdateStoryPlan_Call) Run(run func(ctx context.Context, request services.UpdateStoryPlanRequest)) *MockUpdateStoryPlanService_UpdateStoryPlan_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(services.UpdateStoryPlanRequest))
	})
	return _c
}

func (_c *MockUpdateStoryPlanService_UpdateStoryPlan_Call) Return(_a0 *models.StoryPlan, _a1 error) *MockUpdateStoryPlanService_UpdateStoryPlan_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockUpdateStoryPlanService_UpdateStoryPlan_Call) RunAndReturn(run func(context.Context, services.UpdateStoryPlanRequest) (*models.StoryPlan, error)) *MockUpdateStoryPlanService_UpdateStoryPlan_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockUpdateStoryPlanService creates a new instance of MockUpdateStoryPlanService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockUpdateStoryPlanService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockUpdateStoryPlanService {
	mock := &MockUpdateStoryPlanService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
