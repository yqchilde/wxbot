package covid19

type ApiResponse struct {
	Result []struct {
		DisplayData struct {
			ResultData struct {
				TplData struct {
					Desc     string `json:"desc"`
					DataList []struct {
						TotalDesc string `json:"total_desc"`
						TotalNum  string `json:"total_num"`
					} `json:"data_list"`
					DynamicList []struct {
						DataList []struct {
							TotalDesc  string `json:"total_desc"`
							TotalNum   string `json:"total_num"`
							ChangeDesc string `json:"change_desc,omitempty"`
							ChangeNum  string `json:"change_num,omitempty"`
						} `json:"data_list"`
					} `json:"dynamic_list"`
				} `json:"tplData"`
			} `json:"resultData"`
		} `json:"DisplayData"`
	} `json:"Result"`
}

type EpidemicData struct {
	LastUpdateTime  string // 更新时间
	LocalAdd        string // 新增本土
	LocalNow        string // 现有本土
	LocalAddWzz     string // 新增本土无症状
	LocalNowWzz     string // 现有本土无症状
	ForeignAdd      string // 新增境外
	ForeignNow      string // 现有境外
	HkMacTwAdd      string // 港澳台新增
	ConfirmNow      string // 现有确诊
	ConfirmTotal    string // 累计确诊
	ConfirmTotalAdd string // 累计确诊（较昨日）
	ForeignTotal    string // 累计境外
	ForeignTotalAdd string // 累计境外（较昨日）
	HealTotal       string // 累计治愈
	HealTotalAdd    string // 累计治愈（较昨日）
	DeadTotal       string // 累计死亡
	DeadTotalAdd    string // 累计死亡（较昨日）
}
