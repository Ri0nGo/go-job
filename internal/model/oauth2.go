package model

// OAuth2State 开始请求时的信息
type OAuth2State struct {
	Uid      int
	State    string
	Scene    Auth2Scene // 场景，表示从那个页面发起的操作
	Platform string
	Used     bool
}

// OAuth2TempCode 回调结束后的临时code存储的信息
type OAuth2TempCode struct {
	Uid          int        `json:"uid"`
	Name         string     `json:"name"`     // 第三方平台授权用户的名称
	Identify     string     `json:"identify"` // 第三方平台的唯一标识
	Platform     string     `json:"platform"` // 第三方平台的名称
	RedirectPage string     `json:"redirect_page"`
	Scene        Auth2Scene `json:"scene"`
	Err          string     `json:"err"`
}
