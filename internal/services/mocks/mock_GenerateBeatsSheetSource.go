// Code generated by mockery v2.53.1. DO NOT EDIT.

package servicesmocks

import (
	context "github.com/a-novel-kit/context"
	dao "github.com/a-novel/service-story-schematics/internal/dao"
	daoai "github.com/a-novel/service-story-schematics/internal/daoai"

	mock "github.com/stretchr/testify/mock"

	models "github.com/a-novel/service-story-schematics/models"

	uuid "github.com/google/uuid"
)

// MockGenerateBeatsSheetSource is an autogenerated mock type for the GenerateBeatsSheetSource type
type MockGenerateBeatsSheetSource struct {
	mock.Mock
}

type MockGenerateBeatsSheetSource_Expecter struct {
	mock *mock.Mock
}

func (_m *MockGenerateBeatsSheetSource) EXPECT() *MockGenerateBeatsSheetSource_Expecter {
	return &MockGenerateBeatsSheetSource_Expecter{mock: &_m.Mock}
}

// GenerateBeatsSheet provides a mock function with given fields: ctx, request
func (_m *MockGenerateBeatsSheetSource) GenerateBeatsSheet(ctx context.Context, request daoai.GenerateBeatsSheetRequest) ([]models.Beat, error) {
	ret := _m.Called(ctx, request)

	if len(ret) == 0 {
		panic("no return value specified for GenerateBeatsSheet")
	}

	var r0 []models.Beat
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, daoai.GenerateBeatsSheetRequest) ([]models.Beat, error)); ok {
		return rf(ctx, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, daoai.GenerateBeatsSheetRequest) []models.Beat); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Beat)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, daoai.GenerateBeatsSheetRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockGenerateBeatsSheetSource_GenerateBeatsSheet_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GenerateBeatsSheet'
type MockGenerateBeatsSheetSource_GenerateBeatsSheet_Call struct {
	*mock.Call
}

// GenerateBeatsSheet is a helper method to define mock.On call
//   - ctx context.Context
//   - request daoai.GenerateBeatsSheetRequest
func (_e *MockGenerateBeatsSheetSource_Expecter) GenerateBeatsSheet(ctx interface{}, request interface{}) *MockGenerateBeatsSheetSource_GenerateBeatsSheet_Call {
	return &MockGenerateBeatsSheetSource_GenerateBeatsSheet_Call{Call: _e.mock.On("GenerateBeatsSheet", ctx, request)}
}

func (_c *MockGenerateBeatsSheetSource_GenerateBeatsSheet_Call) Run(run func(ctx context.Context, request daoai.GenerateBeatsSheetRequest)) *MockGenerateBeatsSheetSource_GenerateBeatsSheet_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(daoai.GenerateBeatsSheetRequest))
	})
	return _c
}

func (_c *MockGenerateBeatsSheetSource_GenerateBeatsSheet_Call) Return(_a0 []models.Beat, _a1 error) *MockGenerateBeatsSheetSource_GenerateBeatsSheet_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockGenerateBeatsSheetSource_GenerateBeatsSheet_Call) RunAndReturn(run func(context.Context, daoai.GenerateBeatsSheetRequest) ([]models.Beat, error)) *MockGenerateBeatsSheetSource_GenerateBeatsSheet_Call {
	_c.Call.Return(run)
	return _c
}

// SelectLogline provides a mock function with given fields: ctx, data
func (_m *MockGenerateBeatsSheetSource) SelectLogline(ctx context.Context, data dao.SelectLoglineData) (*dao.LoglineEntity, error) {
	ret := _m.Called(ctx, data)

	if len(ret) == 0 {
		panic("no return value specified for SelectLogline")
	}

	var r0 *dao.LoglineEntity
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, dao.SelectLoglineData) (*dao.LoglineEntity, error)); ok {
		return rf(ctx, data)
	}
	if rf, ok := ret.Get(0).(func(context.Context, dao.SelectLoglineData) *dao.LoglineEntity); ok {
		r0 = rf(ctx, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dao.LoglineEntity)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, dao.SelectLoglineData) error); ok {
		r1 = rf(ctx, data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockGenerateBeatsSheetSource_SelectLogline_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SelectLogline'
type MockGenerateBeatsSheetSource_SelectLogline_Call struct {
	*mock.Call
}

// SelectLogline is a helper method to define mock.On call
//   - ctx context.Context
//   - data dao.SelectLoglineData
func (_e *MockGenerateBeatsSheetSource_Expecter) SelectLogline(ctx interface{}, data interface{}) *MockGenerateBeatsSheetSource_SelectLogline_Call {
	return &MockGenerateBeatsSheetSource_SelectLogline_Call{Call: _e.mock.On("SelectLogline", ctx, data)}
}

func (_c *MockGenerateBeatsSheetSource_SelectLogline_Call) Run(run func(ctx context.Context, data dao.SelectLoglineData)) *MockGenerateBeatsSheetSource_SelectLogline_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(dao.SelectLoglineData))
	})
	return _c
}

func (_c *MockGenerateBeatsSheetSource_SelectLogline_Call) Return(_a0 *dao.LoglineEntity, _a1 error) *MockGenerateBeatsSheetSource_SelectLogline_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockGenerateBeatsSheetSource_SelectLogline_Call) RunAndReturn(run func(context.Context, dao.SelectLoglineData) (*dao.LoglineEntity, error)) *MockGenerateBeatsSheetSource_SelectLogline_Call {
	_c.Call.Return(run)
	return _c
}

// SelectStoryPlan provides a mock function with given fields: ctx, data
func (_m *MockGenerateBeatsSheetSource) SelectStoryPlan(ctx context.Context, data uuid.UUID) (*dao.StoryPlanEntity, error) {
	ret := _m.Called(ctx, data)

	if len(ret) == 0 {
		panic("no return value specified for SelectStoryPlan")
	}

	var r0 *dao.StoryPlanEntity
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*dao.StoryPlanEntity, error)); ok {
		return rf(ctx, data)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *dao.StoryPlanEntity); ok {
		r0 = rf(ctx, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dao.StoryPlanEntity)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockGenerateBeatsSheetSource_SelectStoryPlan_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SelectStoryPlan'
type MockGenerateBeatsSheetSource_SelectStoryPlan_Call struct {
	*mock.Call
}

// SelectStoryPlan is a helper method to define mock.On call
//   - ctx context.Context
//   - data uuid.UUID
func (_e *MockGenerateBeatsSheetSource_Expecter) SelectStoryPlan(ctx interface{}, data interface{}) *MockGenerateBeatsSheetSource_SelectStoryPlan_Call {
	return &MockGenerateBeatsSheetSource_SelectStoryPlan_Call{Call: _e.mock.On("SelectStoryPlan", ctx, data)}
}

func (_c *MockGenerateBeatsSheetSource_SelectStoryPlan_Call) Run(run func(ctx context.Context, data uuid.UUID)) *MockGenerateBeatsSheetSource_SelectStoryPlan_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockGenerateBeatsSheetSource_SelectStoryPlan_Call) Return(_a0 *dao.StoryPlanEntity, _a1 error) *MockGenerateBeatsSheetSource_SelectStoryPlan_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockGenerateBeatsSheetSource_SelectStoryPlan_Call) RunAndReturn(run func(context.Context, uuid.UUID) (*dao.StoryPlanEntity, error)) *MockGenerateBeatsSheetSource_SelectStoryPlan_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockGenerateBeatsSheetSource creates a new instance of MockGenerateBeatsSheetSource. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockGenerateBeatsSheetSource(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockGenerateBeatsSheetSource {
	mock := &MockGenerateBeatsSheetSource{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
