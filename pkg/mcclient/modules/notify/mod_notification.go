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

package notify

import (
	"yunion.io/x/jsonutils"
	"yunion.io/x/pkg/errors"
	"yunion.io/x/pkg/util/sets"

	"yunion.io/x/onecloud/pkg/mcclient"
	"yunion.io/x/onecloud/pkg/mcclient/modulebase"
	"yunion.io/x/onecloud/pkg/mcclient/modules"
)

var (
	Notifications NotificationManager
)

type SNotifyMessage struct {
	Uid         []string        `json:"uid,omitempty"`
	Gid         []string        `json:"gid,omitempty"`
	ContactType TNotifyChannel  `json:"contact_type,omitempty"`
	Topic       string          `json:"topic,omitempty"`
	Priority    TNotifyPriority `json:"priority,omitempty"`
	Msg         string          `json:"msg,omitempty"`
	Remark      string          `json:"remark,omitempty"`
	Broadcast   bool            `json:"broadcast,omitempty"`
}

type SNotifyV2Message struct {
	ReceiverIds []string `json:"receiver_ids"`
	ContactType string   `json:"contact_type"`
	Topic       string   `json:"topic"`
	Priority    string   `json:"priority"`
	Message     string   `json:"message"`
}

type NotificationManager struct {
	modulebase.ResourceManager
}

func (manager *NotificationManager) Send(s *mcclient.ClientSession, msg SNotifyMessage) error {
	receiverIds := make([]string, 0, len(msg.Uid))
	if len(msg.Gid) > 0 {
		// fetch uid
		uidSet := sets.NewString()
		for _, gid := range msg.Gid {
			users, err := modules.Groups.GetUsers(s, gid, nil)
			if err != nil {
				return errors.Wrapf(err, "Groups.GetUsers for group %q", gid)
			}
			for i := range users.Data {
				id, _ := users.Data[i].GetString("id")
				uidSet.Insert(id)
			}
		}
		for _, uid := range uidSet.UnsortedList() {
			receiverIds = append(receiverIds, uid)
		}
	}
	receiverIds = append(receiverIds, msg.Uid...)

	v2msg := SNotifyV2Message{
		ReceiverIds: receiverIds,
		ContactType: string(msg.ContactType),
		Topic:       msg.Topic,
		Priority:    string(msg.Priority),
		Message:     msg.Msg,
	}
	params := jsonutils.Marshal(&v2msg)

	_, err := manager.Create(s, params)
	return err
}

func init() {
	Notifications = NotificationManager{
		modules.Notification,
	}
}
