package domain

//自定义要发布的内容 结构体
type Content struct {
	Title     *string // 仓库名
	Subtitle  *string
	URL       *string
	ExtraData []string
}
