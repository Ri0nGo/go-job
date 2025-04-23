package email

type EmailTplType int

const (
	EmailBindVerfyCodeTpl = iota + 1
	EmailJobNotifyTpl
)

type EmailTpl struct {
	Subject string
	Content string
}

var emailTplMap = map[EmailTplType]EmailTpl{
	EmailBindVerfyCodeTpl: {
		Subject: "[go-job] 验证码",
		Content: `<!DOCTYPE html><html><head><meta charset="UTF-8"><title>邮箱验证码</title></head>
<body style="font-family: Arial, sans-serif; background-color: #fff; padding: 40px;">
<div style="max-width: 600px; margin: 0 auto; text-align: center;">
<p style="font-size: 18px;">你好，<strong style="color: #007BFF;">%s</strong></p>
<p style="font-size: 16px;">你的邮箱验证码为：</p>
<p style="font-size: 48px; font-weight: bold; color: #2c3e50; margin: 30px 0;">%s</p>
<p style="color: #999;">30分钟内有效，请勿向他人泄漏</p>
</div></body></html>`,
	},
	EmailJobNotifyTpl: {
		Subject: "[任务] %s",
		Content: `<!DOCTYPE html><html><head><meta charset="UTF-8">
<style>.c{padding-left:20px}.l{font-weight:bold}.o,.e{border-radius:5px;padding:10px;margin:5px 0 15px;white-space:pre-line}.o{background:#e6f4ea}.e{background:#fdecea;color:#d93025}</style></head>
<body><div class="c"><p><span class="l">任务名称：</span>%s</p><p><span class="l">状态：</span>%s</p><p><span class="l">开始执行时间：</span>%s</p><p><span class="l">耗时：</span>%.2f</p><p><span class="l">输出：</span></p>
<div class="o">%s</div><p><span class="l">异常情况：</span></p><div class="e">%s`,
	},
}

func GetEmailTpl(tplType EmailTplType) EmailTpl {
	return emailTplMap[tplType]
}
