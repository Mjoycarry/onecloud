package tasks

import (
	"context"
	"fmt"

	"yunion.io/x/jsonutils"

	"yunion.io/x/onecloud/pkg/cloudcommon/db"
	"yunion.io/x/onecloud/pkg/cloudcommon/db/taskman"
	"yunion.io/x/onecloud/pkg/cloudprovider"
	"yunion.io/x/onecloud/pkg/compute/models"
	"yunion.io/x/onecloud/pkg/util/logclient"
)

type EipDissociateTask struct {
	taskman.STask
}

func init() {
	taskman.RegisterTask(EipDissociateTask{})
}

func (self *EipDissociateTask) TaskFail(ctx context.Context, eip *models.SElasticip, msg string, vm *models.SGuest) {
	eip.SetStatus(self.UserCred, models.EIP_STATUS_READY, msg)
	self.SetStageFailed(ctx, msg)
	if vm != nil {
		vm.SetStatus(self.UserCred, models.VM_DISSOCIATE_EIP_FAILED, msg)
		db.OpsLog.LogEvent(vm, db.ACT_EIP_DETACH, msg, self.GetUserCred())
		logclient.AddActionLogWithStartable(self, vm, logclient.ACT_EIP_DISSOCIATE, msg, self.UserCred, false)
	}
}

func (self *EipDissociateTask) OnInit(ctx context.Context, obj db.IStandaloneModel, data jsonutils.JSONObject) {
	eip := obj.(*models.SElasticip)

	server := eip.GetAssociateVM()
	if server != nil {

		if server.Status != models.VM_DISSOCIATE_EIP {
			server.SetStatus(self.UserCred, models.VM_DISSOCIATE_EIP, "dissociate eip")
		}

		extEip, err := eip.GetIEip()
		if err != nil && err != cloudprovider.ErrNotFound {
			msg := fmt.Sprintf("fail to find iEIP for eip %s", err)
			self.TaskFail(ctx, eip, msg, server)
			return
		}

		if err == nil && len(extEip.GetAssociationExternalId()) > 0 {
			err = extEip.Dissociate()
			if err != nil {
				msg := fmt.Sprintf("fail to remote dissociate eip %s", err)
				self.TaskFail(ctx, eip, msg, server)
				return
			}
		}

		err = eip.Dissociate(ctx, self.UserCred)
		if err != nil {
			msg := fmt.Sprintf("fail to local dissociate eip %s", err)
			self.TaskFail(ctx, eip, msg, server)
			return
		}

		eip.SetStatus(self.UserCred, models.EIP_STATUS_READY, "dissociate")

		logclient.AddActionLogWithStartable(self, server, logclient.ACT_EIP_DISSOCIATE, nil, self.UserCred, true)
		server.StartSyncstatus(ctx, self.UserCred, "")
	}

	self.SetStageComplete(ctx, nil)

	autoDelete := jsonutils.QueryBoolean(self.GetParams(), "auto_delete", false)

	if eip.AutoDellocate.IsTrue() || autoDelete {
		eip.StartEipDeallocateTask(ctx, self.UserCred, "")
	}
}
