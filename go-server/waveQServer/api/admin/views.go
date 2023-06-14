package admin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"waveQServer/src/comm"
	"waveQServer/src/config"
	"waveQServer/src/core/groups"
	queueImpl2 "waveQServer/src/core/queue/queueImpl"
	"waveQServer/src/identity"
	"waveQServer/src/utils"
	"waveQServer/src/utils/jwtutil"
	"waveQServer/src/utils/logutil"
)

// Login 管理员登录
func Login(c *gin.Context) {
	admin := identity.NewAdmin()
	err := c.ShouldBindJSON(admin)
	if err != nil {
		logutil.LogError(err.Error())
		fail := comm.Fail(err.Error())
		c.JSON(http.StatusBadRequest, fail)
		c.Abort()
		return
	}
	if utils.NotEquals(config.GetConfig().UserName, utils.Md5([]byte(admin.UserName))) {
		fail := comm.Fail("username error")
		c.JSON(http.StatusBadRequest, fail)
		c.Abort()
		return
	}
	if utils.NotEquals(config.GetConfig().Password, utils.Md5([]byte(admin.Password))) {
		fail := comm.Fail("password error")
		c.JSON(http.StatusBadRequest, fail)
		c.Abort()
		return
	}
	token, err := jwtutil.GetToken(admin.UserName, admin.Password)
	if err != nil {
		fail := comm.Fail(err.Error())
		c.JSON(http.StatusBadRequest, fail)
		c.Abort()
		return
	}
	m := make(map[string]string)
	m["XMD-TOKEN"] = token
	ok := comm.OK(m)
	c.JSON(http.StatusOK, ok)
	c.Abort()
	return
}

// CreateGroup 创建一个组
func CreateGroup(c *gin.Context) {
	group := make(map[string]string)
	err := c.ShouldBindJSON(group)
	if err != nil {
		comm.DisposeError(err, c)
		return
	}
	_, err = groups.NewGroup(group["groupId"])
	if err != nil {
		comm.DisposeError(err, c)
		return
	}
	c.JSON(http.StatusOK, comm.OK())
	return
}

func CreateQueue(c *gin.Context) {
	group := make(map[string]string)
	err := c.ShouldBindJSON(group)
	if err != nil {
		comm.DisposeError(err, c)
		return
	}
	queueType := group["queueType"]
	groupId := group["GroupId"]
	squeueId := group["queueId"]
	switch queueType {
	case "1":
		delayQueue, err := queueImpl2.NewDelayQueue([]byte(squeueId))
		if err != nil {
			comm.DisposeError(err, c)
			return
		}
		err = groups.GetGroupById(groupId).BindQueue(delayQueue)
		if err != nil {
			comm.DisposeError(err, c)
			return
		}
		break
	case "2":
		queue, err := queueImpl2.NewBroadcastQueue([]byte(squeueId))
		if err != nil {
			comm.DisposeError(err, c)
			return
		}
		err = groups.GetGroupById(groupId).BindQueue(queue)
		if err != nil {
			comm.DisposeError(err, c)
			return
		}
		break
	}

}
