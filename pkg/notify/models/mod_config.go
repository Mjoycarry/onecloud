// Copyright 2019 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package models

import (
	"context"
	"fmt"

	"yunion.io/x/jsonutils"
	"yunion.io/x/log"
	"yunion.io/x/pkg/errors"
	"yunion.io/x/sqlchemy"

	"yunion.io/x/onecloud/pkg/cloudcommon/db"
	"yunion.io/x/onecloud/pkg/mcclient"
	_interface "yunion.io/x/onecloud/pkg/notify/interface"
	"yunion.io/x/onecloud/pkg/notify/options"
)

type SConfigManager struct {
	SStatusStandaloneResourceBaseManager
}

var ConfigManager *SConfigManager

func init() {
	ConfigManager = &SConfigManager{
		SStatusStandaloneResourceBaseManager: NewStatusStandaloneResourceBaseManager(
			SConfig{},
			"notify_t_config",
			"config",
			"configs",
		),
	}
	ConfigManager.SetVirtualObject(ConfigManager)
}

// SConfig is a table which storage (k,v) and its type.
// The three important concepts are key, value and type.
// Key and type uniquely identify a value.
type SConfig struct {
	SStatusStandaloneResourceBase

	Type      string `width:"15" nullable:"false" create:"required" list:"user"`
	KeyText   string `width:"30" nullable:"false" create:"required" list:"user"`
	ValueText string `width:"100" nullable:"false" create:"required" list:"user"`
}

// ListItemFilter is a hook function belong to IModelManager interface when Listing.
// This will Called in yunion.io/x/onecloud/pkg/cloudcommon/db.List function.
func (self *SConfigManager) ListItemFilter(ctx context.Context, q *sqlchemy.SQuery, userCred mcclient.TokenCredential, query jsonutils.JSONObject) (*sqlchemy.SQuery, error) {
	if !query.Contains("type") {
		return q, nil
	}
	contactType, _ := query.GetString("type")
	q.Filter(sqlchemy.Equals(q.Field("type"), contactType))
	return q, nil
}

// GetValue fetch the SConfig struct corresponding to key and type.
func (self *SConfigManager) GetValue(key, contactType string) (*SConfig, error) {
	q := self.Query()
	q.Filter(sqlchemy.AND(sqlchemy.Equals(q.Field("type"), contactType), sqlchemy.Equals(q.Field("key_text"), key)))
	configs := make([]SConfig, 0, 1)
	err := db.FetchModelObjects(self, q, &configs)
	if err != nil {
		return nil, errors.Wrap(err, "Fetch SConfig by key and type failed")
	}
	if len(configs) == 0 {
		return nil, errors.Error("There is no SConfig whose key and type meet the requirement")
	}
	return &configs[0], nil
}

func (self *SConfigManager) InitializeData() error {
	q := self.Query()
	q = q.Filter(sqlchemy.OR(sqlchemy.IsNotNull(q.Field("updated_at")), sqlchemy.IsNotNull(q.Field("created_at")),
		sqlchemy.IsNotNull(q.Field("deleted_at")), sqlchemy.IsTrue(q.Field("deleted")),
		sqlchemy.IsNotNull(q.Field("update_by")), sqlchemy.IsNotNull(q.Field("delete_by"))))
	n, err := q.CountWithError()
	if err != nil {
		return err
	}
	if n <= 0 {
		log.Debugf("need to init data for %s", self.TableSpec().Name())
		// need to update
		sql := fmt.Sprintf("update %s set updated_at=gmt_modified, deleted=is_deleted, created_at=gmt_create, deleted_at=gmt_deleted, update_by=modified_by, delete_by=deleted_by", self.TableSpec().Name())
		q := sqlchemy.NewRawQuery(sql, "")
		q.Row()
		sql = fmt.Sprintf("update %s set type='mobile' where type='sms_aliyun'", self.TableSpec().Name())
		q = sqlchemy.NewRawQuery(sql, "")
		q.Row()
	}

	// init webconsole's config
	sql := fmt.Sprintf("update %s set deleted='1' where type='webconsole'", self.TableSpec().Name())
	q = sqlchemy.NewRawQuery(sql)
	q.Row()
	configs := []SConfig{
		{
			Type:      "webconsole",
			KeyText:   "auth_uri",
			ValueText: options.Options.AuthURL,
		},
		{
			Type:      "webconsole",
			KeyText:   "admin_user",
			ValueText: options.Options.AdminUser,
		},
		{
			Type:      "webconsole",
			KeyText:   "admin_password",
			ValueText: options.Options.AdminPassword,
		},
		{
			Type:      "webconsole",
			KeyText:   "admin_tenant_name",
			ValueText: options.Options.AdminProject,
		},
	}
	for _, config := range configs {
		err := self.TableSpec().Insert(&config)
		if err != nil {
			return err
		}
	}
	return nil
}

// Fetch all SConfig struct which type is contactType.
func (self *SConfigManager) GetConfigByType(contactType string) ([]SConfig, error) {
	q := self.Query()
	q.Filter(sqlchemy.Equals(q.Field("type"), contactType))
	configs := make([]SConfig, 0, 5)
	err := db.FetchModelObjects(self, q, &configs)
	if err != nil {
		return nil, errors.Wrap(err, "Fetch SConfigs by type failed")
	}
	//if len(configs) == 0 {
	//	return nil, errors.Error("There is no SConfig whose type meet the requirement")
	//}
	return configs, nil
}

func (self *SConfigManager) GetConfig(contactType string) (_interface.SConfig, error) {
	configs, err := self.GetConfigByType(contactType)
	if err != nil {
		return nil, err
	}
	ret := make(map[string]string)
	for i := range configs {
		ret[configs[i].KeyText] = configs[i].ValueText
	}
	return ret, nil
}

func (self *SConfigManager) SetConfig(contactType string, config _interface.SConfig) error {
	return fmt.Errorf("SetConfig Not Implemented")
}
