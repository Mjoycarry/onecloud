package tasks

import (
	"context"

	"yunion.io/x/jsonutils"
	"yunion.io/x/log"

	"yunion.io/x/onecloud/pkg/cloudcommon/db"
	"yunion.io/x/onecloud/pkg/cloudcommon/db/taskman"
	"yunion.io/x/onecloud/pkg/compute/models"
	"yunion.io/x/onecloud/pkg/util/logclient"
)

type GuestStopTask struct {
	SGuestBaseTask
}

func init() {
	taskman.RegisterTask(GuestStopTask{})
}

func (self *GuestStopTask) OnInit(ctx context.Context, obj db.IStandaloneModel, data jsonutils.JSONObject) {
	guest := obj.(*models.SGuest)
	db.OpsLog.LogEvent(guest, db.ACT_STOPPING, nil, self.UserCred)
	self.stopGuest(ctx, guest)
}

func (self *GuestStopTask) isSubtask() bool {
	return jsonutils.QueryBoolean(self.Params, "subtask", false)
}

func (self *GuestStopTask) stopGuest(ctx context.Context, guest *models.SGuest) {
	host := guest.GetHost()
	if host == nil {
		self.OnGuestStopTaskCompleteFailed(ctx, guest, jsonutils.NewString("no associated host"))
		return
	}
	if !self.isSubtask() {
		guest.SetStatus(self.UserCred, models.VM_STOPPING, "")
	}
	self.SetStage("OnMasterStopTaskComplete", nil)
	err := guest.GetDriver().RequestStopOnHost(ctx, guest, host, self)
	if err != nil {
		log.Errorf("RequestStopOnHost fail %s", err)
		self.OnGuestStopTaskCompleteFailed(ctx, guest, jsonutils.NewString(err.Error()))
	}
}

func (self *GuestStopTask) OnMasterStopTaskComplete(ctx context.Context, guest *models.SGuest, data jsonutils.JSONObject) {
	if len(guest.BackupHostId) > 0 {
		host := models.HostManager.FetchHostById(guest.BackupHostId)
		self.SetStage("OnGuestStopTaskComplete", nil)
		err := guest.GetDriver().RequestStopOnHost(ctx, guest, host, self)
		if err != nil {
			log.Errorf("RequestStopOnHost fail %s", err)
			self.OnGuestStopTaskCompleteFailed(ctx, guest, jsonutils.NewString(err.Error()))
		}
	} else {
		self.OnGuestStopTaskComplete(ctx, guest, data)
	}
}

func (self *GuestStopTask) OnMasterStopTaskCompleteFailed(ctx context.Context, obj db.IStandaloneModel, reason jsonutils.JSONObject) {
	guest := obj.(*models.SGuest)
	self.OnGuestStopTaskCompleteFailed(ctx, guest, reason)
}

func (self *GuestStopTask) OnGuestStopTaskComplete(ctx context.Context, guest *models.SGuest, data jsonutils.JSONObject) {
	if !self.isSubtask() {
		guest.SetStatus(self.UserCred, models.VM_READY, "")
	}
	db.OpsLog.LogEvent(guest, db.ACT_STOP, guest.GetShortDesc(ctx), self.UserCred)
	models.HostManager.ClearSchedDescCache(guest.HostId)
	self.SetStageComplete(ctx, nil)
	if guest.Status == models.VM_READY && guest.DisableDelete.IsFalse() && guest.ShutdownBehavior == models.SHUTDOWN_TERMINATE {
		guest.StartAutoDeleteGuestTask(ctx, self.UserCred, "")
	}
	logclient.AddActionLogWithStartable(self, guest, logclient.ACT_VM_STOP, "", self.UserCred, true)
}

func (self *GuestStopTask) OnGuestStopTaskCompleteFailed(ctx context.Context, guest *models.SGuest, reason jsonutils.JSONObject) {
	guest.SetStatus(self.UserCred, models.VM_STOP_FAILED, reason.String())
	db.OpsLog.LogEvent(guest, db.ACT_STOP_FAIL, reason.String(), self.UserCred)
	self.SetStageFailed(ctx, reason.String())
	logclient.AddActionLogWithStartable(self, guest, logclient.ACT_VM_STOP, reason.String(), self.UserCred, false)
}
