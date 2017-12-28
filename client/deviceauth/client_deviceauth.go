// Copyright 2017 Northern.tech AS
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
package deviceauth

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/mendersoftware/go-lib-micro/log"
	"github.com/mendersoftware/go-lib-micro/rest_utils"
	"github.com/pkg/errors"

	"github.com/mendersoftware/deviceadm/client"
	"github.com/mendersoftware/deviceadm/utils"
)

const (
	// default device ID endpoint
	defaultDevAuthDevicesUri = "/api/management/v1/devauth/devices/{id}/auth/{aid}/status"
	// default preauthorize device endpoint
	defaultPreauthorizeDeviceUri = "/api/management/v1/devauth/devices"
	// default request timeout, 10s?
	defaultDevAuthReqTimeout = time.Duration(10) * time.Second
)

type Config struct {
	// root devauth address
	DevauthUrl string
	// request timeout
	Timeout time.Duration
}

type Client struct {
	client client.HttpRunner
	conf   Config
}

// devauth's status request
type StatusReq struct {
	DeviceId string `json:"-"`
	AuthId   string `json:"-"`
	Status   string `json:"status"`
}

type PreAuthReq struct {
	DeviceId  string `json:"device_id" valid:"required" bson:"device_id"`
	AuthSetId string `json:"auth_set_id" valid:"required" bson:"auth_set_id"`
	IdData    string `json:"id_data" valid:"required" bson:"id_data"`
	PubKey    string `json:"pubkey" valid:"required" bson:"pubkey"`
}

// TODO rename this and calling funcs to UpdateDeviceStatus etc.
// perhaps change the interface - the whole Device isn't needed
// leaving for later, requires large refact in tests etc.
func (d *Client) UpdateDevice(ctx context.Context, sreq StatusReq) error {
	l := log.FromContext(ctx)
	l.Debugf("update device %s", sreq.DeviceId)

	url := d.buildDevAuthUpdateUrl(sreq)

	statusReqJson, err := json.Marshal(sreq)
	if err != nil {
		return errors.Wrapf(err, "failed to prepare dev auth request")
	}

	reader := bytes.NewReader(statusReqJson)

	req, err := http.NewRequest(http.MethodPut, url, reader)
	if err != nil {
		return errors.Wrapf(err, "failed to prepare dev auth request")
	}

	req.Header.Set("Content-Type", "application/json")

	// set request timeout and setup cancellation
	ctx, cancel := context.WithTimeout(ctx, d.conf.Timeout)
	defer cancel()
	rsp, err := d.client.Do(req.WithContext(ctx))
	if err != nil {
		return errors.Wrapf(err, "failed to update device status")
	}
	defer rsp.Body.Close()

	switch rsp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusUnprocessableEntity:
		err := rest_utils.ParseApiError(rsp.Body)
		if rest_utils.IsApiError(err) {
			return utils.NewUsageError(err.Error())
		} else {
			return errors.Wrap(err, "device status update request failed")
		}
	default:
		return errors.Errorf("device status update request failed with status %v",
			rsp.Status)
	}
}

func (d *Client) PreauthorizeDevice(
	ctx context.Context, preAuthReq *PreAuthReq, authorizationHeader string) error {
	url := d.conf.DevauthUrl + defaultPreauthorizeDeviceUri

	preAuthReqJSON, err := json.Marshal(preAuthReq)
	if err != nil {
		return errors.Wrapf(err, "failed to prepare dev auth request")
	}

	reader := bytes.NewReader(preAuthReqJSON)

	req, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		return errors.Wrapf(err, "failed to prepare dev auth POST request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authorizationHeader)

	// set request timeout and setup cancellation
	ctx, cancel := context.WithTimeout(ctx, d.conf.Timeout)
	defer cancel()
	rsp, err := d.client.Do(req.WithContext(ctx))
	if err != nil {
		return errors.Wrapf(err, "failed to preauthorize device")
	}
	defer rsp.Body.Close()

	switch rsp.StatusCode {
	case http.StatusCreated:
		return nil
	default:
		return errors.Errorf("device preauthorize request failed with status %v",
			rsp.Status)
	}
}

func NewClient(c Config, client client.HttpRunner) *Client {

	// use default timeout if none was provided
	if c.Timeout == 0 {
		c.Timeout = defaultDevAuthReqTimeout
	}
	return &Client{
		client: client,
		conf:   c,
	}
}

func (d *Client) buildDevAuthUpdateUrl(req StatusReq) string {
	repl := strings.NewReplacer("{id}", req.DeviceId,
		"{aid}", req.AuthId)

	return repl.Replace(d.conf.DevauthUrl + defaultDevAuthDevicesUri)
}
