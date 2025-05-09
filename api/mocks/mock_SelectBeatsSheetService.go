// Code generated by mockery v2.53.1. DO NOT EDIT.

package apimocks

import (
	context "github.com/a-novel-kit/context"
	mock "github.com/stretchr/testify/mock"

	models "github.com/a-novel/service-story-schematics/models"

	services "github.com/a-novel/service-story-schematics/internal/services"
)

// MockSelectBeatsSheetService is an autogenerated mock type for the SelectBeatsSheetService type
type MockSelectBeatsSheetService struct {
	mock.Mock
}

type MockSelectBeatsSheetService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockSelectBeatsSheetService) EXPECT() *MockSelectBeatsSheetService_Expecter {
	return &MockSelectBeatsSheetService_Expecter{mock: &_m.Mock}
}

// SelectBeatsSheet provides a mock function with given fields: ctx, request
func (_m *MockSelectBeatsSheetService) SelectBeatsSheet(ctx context.Context, request services.SelectBeatsSheetRequest) (*models.BeatsSheet, error) {
	ret := _m.Called(ctx, request)

	if len(ret) == 0 {
		panic("no return value specified for SelectBeatsSheet")
	}

	var r0 *models.BeatsSheet
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, services.SelectBeatsSheetRequest) (*models.BeatsSheet, error)); ok {
		return rf(ctx, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, services.SelectBeatsSheetRequest) *models.BeatsSheet); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.BeatsSheet)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, services.SelectBeatsSheetRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockSelectBeatsSheetService_SelectBeatsSheet_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SelectBeatsSheet'
type MockSelectBeatsSheetService_SelectBeatsSheet_Call struct {
	*mock.Call
}

// SelectBeatsSheet is a helper method to define mock.On call
//   - ctx context.Context
//   - request services.SelectBeatsSheetRequest
func (_e *MockSelectBeatsSheetService_Expecter) SelectBeatsSheet(ctx interface{}, request interface{}) *MockSelectBeatsSheetService_SelectBeatsSheet_Call {
	return &MockSelectBeatsSheetService_SelectBeatsSheet_Call{Call: _e.mock.On("SelectBeatsSheet", ctx, request)}
}

func (_c *MockSelectBeatsSheetService_SelectBeatsSheet_Call) Run(run func(ctx context.Context, request services.SelectBeatsSheetRequest)) *MockSelectBeatsSheetService_SelectBeatsSheet_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(services.SelectBeatsSheetRequest))
	})
	return _c
}

func (_c *MockSelectBeatsSheetService_SelectBeatsSheet_Call) Return(_a0 *models.BeatsSheet, _a1 error) *MockSelectBeatsSheetService_SelectBeatsSheet_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockSelectBeatsSheetService_SelectBeatsSheet_Call) RunAndReturn(run func(context.Context, services.SelectBeatsSheetRequest) (*models.BeatsSheet, error)) *MockSelectBeatsSheetService_SelectBeatsSheet_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockSelectBeatsSheetService creates a new instance of MockSelectBeatsSheetService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockSelectBeatsSheetService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockSelectBeatsSheetService {
	mock := &MockSelectBeatsSheetService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
