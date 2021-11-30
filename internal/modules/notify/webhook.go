package notify

import (
	"fmt"
	"html"
	"time"

	"github.com/ouqiang/gocron/internal/models"
	"github.com/ouqiang/gocron/internal/modules/httpclient"
	"github.com/ouqiang/gocron/internal/modules/logger"
	"github.com/ouqiang/gocron/internal/modules/utils"
)

type WebHook struct{}

func (webHook *WebHook) Send(msg Message) {
	model := new(models.Setting)
	webHookSetting, err := model.Webhook()
	if err != nil {
		logger.Error("#webHook#从数据库获取webHook配置失败", err)
		return
	}
	if webHookSetting.Url == "" {
		logger.Error("#webHook#webhook-url为空")
		return
	}
	logger.Debugf("%+v", webHookSetting)
	msg["name"] = utils.EscapeJson(msg["name"].(string))
	msg["output"] = utils.EscapeJson(msg["output"].(string))
	msg["content"] = parseNotifyTemplate(webHookSetting.Template, msg)
	msg["content"] = html.UnescapeString(msg["content"].(string))
	webHook.send(msg, webHookSetting.Url)
}

func (webHook *WebHook) send(msg Message, url string) {
	// content := msg["content"].(string)
	// 更改为钉钉报警格式
	body := &httpclient.DingMsg{
		MsgType: "markdown",
		Markdown: httpclient.MarkDownModel{
			Title: "报警： 定时任务 告急",
			Text: fmt.Sprintf("## 报警： Task_center告急\n\n### 😭😭😭 小D提醒您，任务执行失败\n\n > 任务名称: %s\n\n> 错误日志: \n\n>  %s\n\n> 状态: %s\n\n> 备注: %s", msg["name"], msg["output"], msg["status"], msg["remark"]),
		},
	}
	timeout := 30
	maxTimes := 3
	i := 0
	for i < maxTimes {
		resp := httpclient.PostJson(url, body, timeout)
		if resp.StatusCode == 200 {
			break
		}
		i += 1
		time.Sleep(2 * time.Second)
		if i < maxTimes {
			logger.Errorf("webHook#发送消息失败#%s#消息内容-%s", resp.Body, msg["content"])
		}
	}
}
