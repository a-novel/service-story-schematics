// Code generated by mockery v2.53.1. DO NOT EDIT.

package servicesmocks

import (
	context "context"

	dao "github.com/a-novel/service-story-schematics/internal/dao"
	daoai "github.com/a-novel/service-story-schematics/internal/daoai"

	mock "github.com/stretchr/testify/mock"

	models "github.com/a-novel/service-story-schematics/models"

	uuid "github.com/google/uuid"
)

// MockExpandBeatSource is an autogenerated mock type for the ExpandBeatSource type
type MockExpandBeatSource struct {
	mock.Mock
}

type MockExpandBeatSource_Expecter struct {
	mock *mock.Mock
}

func (_m *MockExpandBeatSource) EXPECT() *MockExpandBeatSource_Expecter {
	return &MockExpandBeatSource_Expecter{mock: &_m.Mock}
}

// ExpandBeat provides a mock function with given fields: ctx, request
func (_m *MockExpandBeatSource) ExpandBeat(ctx context.Context, request daoai.ExpandBeatRequest) (*models.Beat, error) {
	ret := _m.Called(ctx, request)

	if len(ret) == 0 {
		panic("no return value specified for ExpandBeat")
	}

	var r0 *models.Beat
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, daoai.ExpandBeatRequest) (*models.Beat, error)); ok {
		return rf(ctx, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, daoai.ExpandBeatRequest) *models.Beat); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Beat)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, daoai.ExpandBeatRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockExpandBeatSource_ExpandBeat_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExpandBeat'
type MockExpandBeatSource_ExpandBeat_Call struct {
	*mock.Call
}

// ExpandBeat is a helper method to define mock.On call
//   - ctx context.Context
//   - request daoai.ExpandBeatRequest
func (_e *MockExpandBeatSource_Expecter) ExpandBeat(ctx interface{}, request interface{}) *MockExpandBeatSource_ExpandBeat_Call {
	return &MockExpandBeatSource_ExpandBeat_Call{Call: _e.mock.On("ExpandBeat", ctx, request)}
}

func (_c *MockExpandBeatSource_ExpandBeat_Call) Run(run func(ctx context.Context, request daoai.ExpandBeatRequest)) *MockExpandBeatSource_ExpandBeat_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(daoai.ExpandBeatRequest))
	})
	return _c
}

func (_c *MockExpandBeatSource_ExpandBeat_Call) Return(_a0 *models.Beat, _a1 error) *MockExpandBeatSource_ExpandBeat_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockExpandBeatSource_ExpandBeat_Call) RunAndReturn(run func(context.Context, daoai.ExpandBeatRequest) (*models.Beat, error)) *MockExpandBeatSource_ExpandBeat_Call {
	_c.Call.Return(run)
	return _c
}

// SelectBeatsSheet provides a mock function with given fields: ctx, data
func (_m *MockExpandBeatSource) SelectBeatsSheet(ctx context.Context, data uuid.UUID) (*dao.BeatsSheetEntity, error) {
	ret := _m.Called(ctx, data)

	if len(ret) == 0 {
		panic("no return value specified for SelectBeatsSheet")
	}

	var r0 *dao.BeatsSheetEntity
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*dao.BeatsSheetEntity, error)); ok {
		return rf(ctx, data)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *dao.BeatsSheetEntity); ok {
		r0 = rf(ctx, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dao.BeatsSheetEntity)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockExpandBeatSource_SelectBeatsSheet_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SelectBeatsSheet'
type MockExpandBeatSource_SelectBeatsSheet_Call struct {
	*mock.Call
}

// SelectBeatsSheet is a helper method to define mock.On call
//   - ctx context.Context
//   - data uuid.UUID
func (_e *MockExpandBeatSource_Expecter) SelectBeatsSheet(ctx interface{}, data interface{}) *MockExpandBeatSource_SelectBeatsSheet_Call {
	return &MockExpandBeatSource_SelectBeatsSheet_Call{Call: _e.mock.On("SelectBeatsSheet", ctx, data)}
}

func (_c *MockExpandBeatSource_SelectBeatsSheet_Call) Run(run func(ctx context.Context, data uuid.UUID)) *MockExpandBeatSource_SelectBeatsSheet_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockExpandBeatSource_SelectBeatsSheet_Call) Return(_a0 *dao.BeatsSheetEntity, _a1 error) *MockExpandBeatSource_SelectBeatsSheet_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockExpandBeatSource_SelectBeatsSheet_Call) RunAndReturn(run func(context.Context, uuid.UUID) (*dao.BeatsSheetEntity, error)) *MockExpandBeatSource_SelectBeatsSheet_Call {
	_c.Call.Return(run)
	return _c
}

// SelectLogline provides a mock function with given fields: ctx, data
func (_m *MockExpandBeatSource) SelectLogline(ctx context.Context, data dao.SelectLoglineData) (*dao.LoglineEntity, error) {
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

// MockExpandBeatSource_SelectLogline_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SelectLogline'
type MockExpandBeatSource_SelectLogline_Call struct {
	*mock.Call
}

// SelectLogline is a helper method to define mock.On call
//   - ctx context.Context
//   - data dao.SelectLoglineData
func (_e *MockExpandBeatSource_Expecter) SelectLogline(ctx interface{}, data interface{}) *MockExpandBeatSource_SelectLogline_Call {
	return &MockExpandBeatSource_SelectLogline_Call{Call: _e.mock.On("SelectLogline", ctx, data)}
}

func (_c *MockExpandBeatSource_SelectLogline_Call) Run(run func(ctx context.Context, data dao.SelectLoglineData)) *MockExpandBeatSource_SelectLogline_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(dao.SelectLoglineData))
	})
	return _c
}

func (_c *MockExpandBeatSource_SelectLogline_Call) Return(_a0 *dao.LoglineEntity, _a1 error) *MockExpandBeatSource_SelectLogline_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockExpandBeatSource_SelectLogline_Call) RunAndReturn(run func(context.Context, dao.SelectLoglineData) (*dao.LoglineEntity, error)) *MockExpandBeatSource_SelectLogline_Call {
	_c.Call.Return(run)
	return _c
}

// SelectStoryPlan provides a mock function with given fields: ctx, data
func (_m *MockExpandBeatSource) SelectStoryPlan(ctx context.Context, data uuid.UUID) (*dao.StoryPlanEntity, error) {
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

// MockExpandBeatSource_SelectStoryPlan_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SelectStoryPlan'
type MockExpandBeatSource_SelectStoryPlan_Call struct {
	*mock.Call
}

// SelectStoryPlan is a helper method to define mock.On call
//   - ctx context.Context
//   - data uuid.UUID
func (_e *MockExpandBeatSource_Expecter) SelectStoryPlan(ctx interface{}, data interface{}) *MockExpandBeatSource_SelectStoryPlan_Call {
	return &MockExpandBeatSource_SelectStoryPlan_Call{Call: _e.mock.On("SelectStoryPlan", ctx, data)}
}

func (_c *MockExpandBeatSource_SelectStoryPlan_Call) Run(run func(ctx context.Context, data uuid.UUID)) *MockExpandBeatSource_SelectStoryPlan_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockExpandBeatSource_SelectStoryPlan_Call) Return(_a0 *dao.StoryPlanEntity, _a1 error) *MockExpandBeatSource_SelectStoryPlan_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockExpandBeatSource_SelectStoryPlan_Call) RunAndReturn(run func(context.Context, uuid.UUID) (*dao.StoryPlanEntity, error)) *MockExpandBeatSource_SelectStoryPlan_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockExpandBeatSource creates a new instance of MockExpandBeatSource. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockExpandBeatSource(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockExpandBeatSource {
	mock := &MockExpandBeatSource{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
