// Code generated by MockGen. DO NOT EDIT.
// Source: app.go

// Package app is a generated GoMock package.
package app

import (
	game "github.com/Djarvur/allcups-itrally-2020-task/internal/app/game"
	gomock "github.com/golang/mock/gomock"
	io "io"
	reflect "reflect"
	time "time"
)

// MockAppl is a mock of Appl interface
type MockAppl struct {
	ctrl     *gomock.Controller
	recorder *MockApplMockRecorder
}

// MockApplMockRecorder is the mock recorder for MockAppl
type MockApplMockRecorder struct {
	mock *MockAppl
}

// NewMockAppl creates a new mock instance
func NewMockAppl(ctrl *gomock.Controller) *MockAppl {
	mock := &MockAppl{ctrl: ctrl}
	mock.recorder = &MockApplMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAppl) EXPECT() *MockApplMockRecorder {
	return m.recorder
}

// HealthCheck mocks base method
func (m *MockAppl) HealthCheck(arg0 Ctx) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HealthCheck", arg0)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HealthCheck indicates an expected call of HealthCheck
func (mr *MockApplMockRecorder) HealthCheck(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HealthCheck", reflect.TypeOf((*MockAppl)(nil).HealthCheck), arg0)
}

// Start mocks base method
func (m *MockAppl) Start(arg0 time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start
func (mr *MockApplMockRecorder) Start(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockAppl)(nil).Start), arg0)
}

// Balance mocks base method
func (m *MockAppl) Balance(arg0 Ctx) (int, []int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Balance", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].([]int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Balance indicates an expected call of Balance
func (mr *MockApplMockRecorder) Balance(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Balance", reflect.TypeOf((*MockAppl)(nil).Balance), arg0)
}

// Licenses mocks base method
func (m *MockAppl) Licenses(arg0 Ctx) ([]game.License, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Licenses", arg0)
	ret0, _ := ret[0].([]game.License)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Licenses indicates an expected call of Licenses
func (mr *MockApplMockRecorder) Licenses(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Licenses", reflect.TypeOf((*MockAppl)(nil).Licenses), arg0)
}

// IssueLicense mocks base method
func (m *MockAppl) IssueLicense(arg0 Ctx, wallet []int) (game.License, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IssueLicense", arg0, wallet)
	ret0, _ := ret[0].(game.License)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IssueLicense indicates an expected call of IssueLicense
func (mr *MockApplMockRecorder) IssueLicense(arg0, wallet interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IssueLicense", reflect.TypeOf((*MockAppl)(nil).IssueLicense), arg0, wallet)
}

// ExploreArea mocks base method
func (m *MockAppl) ExploreArea(arg0 Ctx, area game.Area) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExploreArea", arg0, area)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExploreArea indicates an expected call of ExploreArea
func (mr *MockApplMockRecorder) ExploreArea(arg0, area interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExploreArea", reflect.TypeOf((*MockAppl)(nil).ExploreArea), arg0, area)
}

// Dig mocks base method
func (m *MockAppl) Dig(arg0 Ctx, licenseID int, pos game.Coord) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Dig", arg0, licenseID, pos)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Dig indicates an expected call of Dig
func (mr *MockApplMockRecorder) Dig(arg0, licenseID, pos interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Dig", reflect.TypeOf((*MockAppl)(nil).Dig), arg0, licenseID, pos)
}

// Cash mocks base method
func (m *MockAppl) Cash(arg0 Ctx, treasure string) ([]int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Cash", arg0, treasure)
	ret0, _ := ret[0].([]int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Cash indicates an expected call of Cash
func (mr *MockApplMockRecorder) Cash(arg0, treasure interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Cash", reflect.TypeOf((*MockAppl)(nil).Cash), arg0, treasure)
}

// MockRepo is a mock of Repo interface
type MockRepo struct {
	ctrl     *gomock.Controller
	recorder *MockRepoMockRecorder
}

// MockRepoMockRecorder is the mock recorder for MockRepo
type MockRepoMockRecorder struct {
	mock *MockRepo
}

// NewMockRepo creates a new mock instance
func NewMockRepo(ctrl *gomock.Controller) *MockRepo {
	mock := &MockRepo{ctrl: ctrl}
	mock.recorder = &MockRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRepo) EXPECT() *MockRepoMockRecorder {
	return m.recorder
}

// LoadStartTime mocks base method
func (m *MockRepo) LoadStartTime() (*time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadStartTime")
	ret0, _ := ret[0].(*time.Time)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadStartTime indicates an expected call of LoadStartTime
func (mr *MockRepoMockRecorder) LoadStartTime() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadStartTime", reflect.TypeOf((*MockRepo)(nil).LoadStartTime))
}

// SaveStartTime mocks base method
func (m *MockRepo) SaveStartTime(t time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveStartTime", t)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveStartTime indicates an expected call of SaveStartTime
func (mr *MockRepoMockRecorder) SaveStartTime(t interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveStartTime", reflect.TypeOf((*MockRepo)(nil).SaveStartTime), t)
}

// LoadTreasureKey mocks base method
func (m *MockRepo) LoadTreasureKey() ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadTreasureKey")
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadTreasureKey indicates an expected call of LoadTreasureKey
func (mr *MockRepoMockRecorder) LoadTreasureKey() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadTreasureKey", reflect.TypeOf((*MockRepo)(nil).LoadTreasureKey))
}

// SaveTreasureKey mocks base method
func (m *MockRepo) SaveTreasureKey(arg0 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveTreasureKey", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveTreasureKey indicates an expected call of SaveTreasureKey
func (mr *MockRepoMockRecorder) SaveTreasureKey(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveTreasureKey", reflect.TypeOf((*MockRepo)(nil).SaveTreasureKey), arg0)
}

// LoadGame mocks base method
func (m *MockRepo) LoadGame() (ReadSeekCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadGame")
	ret0, _ := ret[0].(ReadSeekCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadGame indicates an expected call of LoadGame
func (mr *MockRepoMockRecorder) LoadGame() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadGame", reflect.TypeOf((*MockRepo)(nil).LoadGame))
}

// SaveGame mocks base method
func (m *MockRepo) SaveGame(arg0 io.WriterTo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveGame", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveGame indicates an expected call of SaveGame
func (mr *MockRepoMockRecorder) SaveGame(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveGame", reflect.TypeOf((*MockRepo)(nil).SaveGame), arg0)
}

// SaveResult mocks base method
func (m *MockRepo) SaveResult(arg0 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveResult", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveResult indicates an expected call of SaveResult
func (mr *MockRepoMockRecorder) SaveResult(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveResult", reflect.TypeOf((*MockRepo)(nil).SaveResult), arg0)
}

// SaveError mocks base method
func (m *MockRepo) SaveError(msg string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveError", msg)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveError indicates an expected call of SaveError
func (mr *MockRepoMockRecorder) SaveError(msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveError", reflect.TypeOf((*MockRepo)(nil).SaveError), msg)
}

// MockGameFactory is a mock of GameFactory interface
type MockGameFactory struct {
	ctrl     *gomock.Controller
	recorder *MockGameFactoryMockRecorder
}

// MockGameFactoryMockRecorder is the mock recorder for MockGameFactory
type MockGameFactoryMockRecorder struct {
	mock *MockGameFactory
}

// NewMockGameFactory creates a new mock instance
func NewMockGameFactory(ctrl *gomock.Controller) *MockGameFactory {
	mock := &MockGameFactory{ctrl: ctrl}
	mock.recorder = &MockGameFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockGameFactory) EXPECT() *MockGameFactoryMockRecorder {
	return m.recorder
}

// New mocks base method
func (m *MockGameFactory) New(arg0 Ctx, arg1 game.Config) (game.Game, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "New", arg0, arg1)
	ret0, _ := ret[0].(game.Game)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// New indicates an expected call of New
func (mr *MockGameFactoryMockRecorder) New(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "New", reflect.TypeOf((*MockGameFactory)(nil).New), arg0, arg1)
}

// Continue mocks base method
func (m *MockGameFactory) Continue(arg0 Ctx, arg1 io.ReadSeeker) (game.Game, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Continue", arg0, arg1)
	ret0, _ := ret[0].(game.Game)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Continue indicates an expected call of Continue
func (mr *MockGameFactoryMockRecorder) Continue(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Continue", reflect.TypeOf((*MockGameFactory)(nil).Continue), arg0, arg1)
}

// MockCPU is a mock of CPU interface
type MockCPU struct {
	ctrl     *gomock.Controller
	recorder *MockCPUMockRecorder
}

// MockCPUMockRecorder is the mock recorder for MockCPU
type MockCPUMockRecorder struct {
	mock *MockCPU
}

// NewMockCPU creates a new mock instance
func NewMockCPU(ctrl *gomock.Controller) *MockCPU {
	mock := &MockCPU{ctrl: ctrl}
	mock.recorder = &MockCPUMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCPU) EXPECT() *MockCPUMockRecorder {
	return m.recorder
}

// Consume mocks base method
func (m *MockCPU) Consume(arg0 Ctx, arg1 time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Consume", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Consume indicates an expected call of Consume
func (mr *MockCPUMockRecorder) Consume(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Consume", reflect.TypeOf((*MockCPU)(nil).Consume), arg0, arg1)
}

// MockLicenseSvc is a mock of LicenseSvc interface
type MockLicenseSvc struct {
	ctrl     *gomock.Controller
	recorder *MockLicenseSvcMockRecorder
}

// MockLicenseSvcMockRecorder is the mock recorder for MockLicenseSvc
type MockLicenseSvcMockRecorder struct {
	mock *MockLicenseSvc
}

// NewMockLicenseSvc creates a new mock instance
func NewMockLicenseSvc(ctrl *gomock.Controller) *MockLicenseSvc {
	mock := &MockLicenseSvc{ctrl: ctrl}
	mock.recorder = &MockLicenseSvcMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLicenseSvc) EXPECT() *MockLicenseSvcMockRecorder {
	return m.recorder
}

// Call mocks base method
func (m *MockLicenseSvc) Call(ctx Ctx, percentFail int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Call", ctx, percentFail)
	ret0, _ := ret[0].(error)
	return ret0
}

// Call indicates an expected call of Call
func (mr *MockLicenseSvcMockRecorder) Call(ctx, percentFail interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Call", reflect.TypeOf((*MockLicenseSvc)(nil).Call), ctx, percentFail)
}

// MockReadSeekCloser is a mock of ReadSeekCloser interface
type MockReadSeekCloser struct {
	ctrl     *gomock.Controller
	recorder *MockReadSeekCloserMockRecorder
}

// MockReadSeekCloserMockRecorder is the mock recorder for MockReadSeekCloser
type MockReadSeekCloserMockRecorder struct {
	mock *MockReadSeekCloser
}

// NewMockReadSeekCloser creates a new mock instance
func NewMockReadSeekCloser(ctrl *gomock.Controller) *MockReadSeekCloser {
	mock := &MockReadSeekCloser{ctrl: ctrl}
	mock.recorder = &MockReadSeekCloserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockReadSeekCloser) EXPECT() *MockReadSeekCloserMockRecorder {
	return m.recorder
}

// Read mocks base method
func (m *MockReadSeekCloser) Read(p []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", p)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read
func (mr *MockReadSeekCloserMockRecorder) Read(p interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockReadSeekCloser)(nil).Read), p)
}

// Seek mocks base method
func (m *MockReadSeekCloser) Seek(offset int64, whence int) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Seek", offset, whence)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Seek indicates an expected call of Seek
func (mr *MockReadSeekCloserMockRecorder) Seek(offset, whence interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Seek", reflect.TypeOf((*MockReadSeekCloser)(nil).Seek), offset, whence)
}

// Close mocks base method
func (m *MockReadSeekCloser) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockReadSeekCloserMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockReadSeekCloser)(nil).Close))
}