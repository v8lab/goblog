package manageController

import (
	"github.com/jinzhu/now"
	"github.com/kataras/iris/v12"
	"kandaoni.com/anqicms/config"
	"kandaoni.com/anqicms/dao"
	"kandaoni.com/anqicms/model"
	"kandaoni.com/anqicms/provider"
	"kandaoni.com/anqicms/response"
	"time"
)

// StatisticSpider 蜘蛛爬行情况
func StatisticSpider(ctx iris.Context) {
	//支持按天，按小时区分
	separate := ctx.URLParam("separate")

	result := provider.StatisticSpider(separate)

	ctx.JSON(iris.Map{
		"code": config.StatusOK,
		"msg":  "",
		"data": result,
	})
}

func StatisticTraffic(ctx iris.Context) {
	//支持按天，按小时区分
	separate := ctx.URLParam("separate")

	result := provider.StatisticTraffic(separate)

	ctx.JSON(iris.Map{
		"code": config.StatusOK,
		"msg":  "",
		"data": result,
	})
}

func StatisticDetail(ctx iris.Context) {
	currentPage := ctx.URLParamIntDefault("current", 1)
	pageSize := ctx.URLParamIntDefault("pageSize", 20)
	isSpider, _ := ctx.URLParamBool("is_spider")

	list, total, _ := provider.StatisticDetail(isSpider, currentPage, pageSize)

	ctx.JSON(iris.Map{
		"code":  config.StatusOK,
		"msg":   "",
		"total": total,
		"data":  list,
	})
}

func GetSpiderIncludeDetail(ctx iris.Context) {
	currentPage := ctx.URLParamIntDefault("current", 1)
	pageSize := ctx.URLParamIntDefault("pageSize", 20)
	var list []*model.SpiderInclude
	var total int64

	if currentPage < 1 {
		currentPage = 1
	}
	offset := (currentPage - 1) * pageSize

	builder := dao.DB.Model(&model.SpiderInclude{})

	builder.Count(&total).Limit(pageSize).Offset(offset).Order("`id` desc").Find(&list)

	ctx.JSON(iris.Map{
		"code":  config.StatusOK,
		"msg":   "",
		"total": total,
		"data":  list,
	})
}

func GetSpiderInclude(ctx iris.Context) {
	var result = make([]response.ChartData, 0, 30*5)

	timeStamp := now.BeginningOfDay().AddDate(0, 0, -30).Unix()

	var includeLogs []model.SpiderInclude
	dao.DB.Model(&model.SpiderInclude{}).Where("`created_time` >= ?", timeStamp).
		Order("created_time asc").
		Scan(&includeLogs)

	lastDate := ""
	for _, v := range includeLogs {
		date := time.Unix(v.CreatedTime, 0).Format("01-02")
		if date == lastDate {
			continue
		}
		lastDate = date
		result = append(result, response.ChartData{
			Date:  date,
			Label: "百度",
			Value: v.BaiduCount,
		}, response.ChartData{
			Date:  date,
			Label: "搜狗",
			Value: v.SogouCount,
		}, response.ChartData{
			Date:  date,
			Label: "搜搜",
			Value: v.SoCount,
		}, response.ChartData{
			Date:  date,
			Label: "必应",
			Value: v.BingCount,
		}, response.ChartData{
			Date:  date,
			Label: "谷歌",
			Value: v.GoogleCount,
		})
	}

	ctx.JSON(iris.Map{
		"code": config.StatusOK,
		"msg":  "",
		"data": result,
	})
}
