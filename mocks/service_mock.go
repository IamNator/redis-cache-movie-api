// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/iamnator/movie-api/service (interfaces: IServices)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/iamnator/movie-api/model"
)

// MockIServices is a mock of IServices interface.
type MockIServices struct {
	ctrl     *gomock.Controller
	recorder *MockIServicesMockRecorder
}

// MockIServicesMockRecorder is the mock recorder for MockIServices.
type MockIServicesMockRecorder struct {
	mock *MockIServices
}

// NewMockIServices creates a new mock instance.
func NewMockIServices(ctrl *gomock.Controller) *MockIServices {
	mock := &MockIServices{ctrl: ctrl}
	mock.recorder = &MockIServicesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIServices) EXPECT() *MockIServicesMockRecorder {
	return m.recorder
}

// GetCharactersByMovieID mocks base method.
func (m *MockIServices) GetCharactersByMovieID(arg0 model.GetCharactersByMovieIDArgs) (*model.CharacterList, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCharactersByMovieID", arg0)
	ret0, _ := ret[0].(*model.CharacterList)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetCharactersByMovieID indicates an expected call of GetCharactersByMovieID.
func (mr *MockIServicesMockRecorder) GetCharactersByMovieID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCharactersByMovieID", reflect.TypeOf((*MockIServices)(nil).GetCharactersByMovieID), arg0)
}

// GetComment mocks base method.
func (m *MockIServices) GetComment(arg0, arg1, arg2 int) ([]model.Comment, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetComment", arg0, arg1, arg2)
	ret0, _ := ret[0].([]model.Comment)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetComment indicates an expected call of GetComment.
func (mr *MockIServicesMockRecorder) GetComment(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetComment", reflect.TypeOf((*MockIServices)(nil).GetComment), arg0, arg1, arg2)
}

// GetMovieByID mocks base method.
func (m *MockIServices) GetMovieByID(arg0 int) (*model.Movie, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMovieByID", arg0)
	ret0, _ := ret[0].(*model.Movie)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMovieByID indicates an expected call of GetMovieByID.
func (mr *MockIServicesMockRecorder) GetMovieByID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMovieByID", reflect.TypeOf((*MockIServices)(nil).GetMovieByID), arg0)
}

// GetMovies mocks base method.
func (m *MockIServices) GetMovies(arg0, arg1 int) ([]model.Movie, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMovies", arg0, arg1)
	ret0, _ := ret[0].([]model.Movie)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetMovies indicates an expected call of GetMovies.
func (mr *MockIServicesMockRecorder) GetMovies(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMovies", reflect.TypeOf((*MockIServices)(nil).GetMovies), arg0, arg1)
}

// SaveComment mocks base method.
func (m *MockIServices) SaveComment(arg0 int, arg1 model.Comment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveComment", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveComment indicates an expected call of SaveComment.
func (mr *MockIServicesMockRecorder) SaveComment(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveComment", reflect.TypeOf((*MockIServices)(nil).SaveComment), arg0, arg1)
}
