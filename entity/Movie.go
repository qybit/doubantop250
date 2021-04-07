package entity

type Movie struct {
	Title    string `json:"title"`// 中文名
	Subtitle string `json:"subtitle"`// 英文名
	Other    string `json:"other"`// 港澳台翻译名
	Cover    string `json:"cover"`// 电影封面
	Desc     string `json:"desc"`// 描述
	Year     string `json:"year"`// 上映年份
	Area     string `json:"area"`// 属于哪个国家
	Tag      string `json:"tag"`// 属于哪一类型的电影
	Star     string `json:"star"`// 评分
	Comment  string `json:"comment"`// 参与评分的人数
	Quote    string `json:"quote"`// 宣传标语
}

