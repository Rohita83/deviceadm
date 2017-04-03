// Copyright 2016 Mender Software AS
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
package mocks

import mock "github.com/stretchr/testify/mock"
import model "github.com/mendersoftware/deviceadm/model"
import store "github.com/mendersoftware/deviceadm/store"

// DataStore is an autogenerated mock type for the DataStore type
type DataStore struct {
	mock.Mock
}

// DeleteDeviceAuth provides a mock function with given fields: id
func (_m *DataStore) DeleteDeviceAuth(id model.AuthID) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(model.AuthID) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteDeviceAuthByDevice provides a mock function with given fields: id
func (_m *DataStore) DeleteDeviceAuthByDevice(id model.DeviceID) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(model.DeviceID) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetDeviceAuth provides a mock function with given fields: id
func (_m *DataStore) GetDeviceAuth(id model.AuthID) (*model.DeviceAuth, error) {
	ret := _m.Called(id)

	var r0 *model.DeviceAuth
	if rf, ok := ret.Get(0).(func(model.AuthID) *model.DeviceAuth); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.DeviceAuth)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(model.AuthID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDeviceAuths provides a mock function with given fields: skip, limit, filter
func (_m *DataStore) GetDeviceAuths(skip int, limit int, filter store.Filter) ([]model.DeviceAuth, error) {
	ret := _m.Called(skip, limit, filter)

	var r0 []model.DeviceAuth
	if rf, ok := ret.Get(0).(func(int, int, store.Filter) []model.DeviceAuth); ok {
		r0 = rf(skip, limit, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.DeviceAuth)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, int, store.Filter) error); ok {
		r1 = rf(skip, limit, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutDeviceAuth provides a mock function with given fields: dev
func (_m *DataStore) PutDeviceAuth(dev *model.DeviceAuth) error {
	ret := _m.Called(dev)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.DeviceAuth) error); ok {
		r0 = rf(dev)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
