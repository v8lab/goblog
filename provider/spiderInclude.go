package provider

import (
	"crypto/tls"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"kandaoni.com/anqicms/config"
	"kandaoni.com/anqicms/dao"
	"kandaoni.com/anqicms/model"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// QuerySpiderInclude 记录 网站的收录情况
// https://www.baidu.com/s?wd=site%3Awww.baidu.com&tn=json&rn=10
func QuerySpiderInclude() {
	if dao.DB == nil {
		return
	}
	link, _ := url.Parse(config.JsonData.System.BaseUrl)
	includeLog := model.SpiderInclude{
		BaiduCount:  GetBaiduInclude(link.Host),
		SogouCount:  GetSogouInclude(link.Host),
		SoCount:     GetSoInclude(link.Host),
		BingCount:   GetBingInclude(link.Host),
		GoogleCount: GetGoogleInclude(link.Host),
	}

	dao.DB.Create(&includeLog)
}

func GetBaiduInclude(serverHost string) (total int) {
	link := fmt.Sprintf("https://www.baidu.com/s?wd=site%%3A%s&tn=json&rn=10", serverHost)

	req := gorequest.New().SetDoNotClearSuperAgent(true).TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).Timeout(5 * time.Second)
	//set key header
	req.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.81 Safari/537.36")
	_, body, errs := req.Get(link).End()
	if errs != nil {
		return total
	}

	// "all": 173525,
	reg := regexp.MustCompile(`(?i)"all":\s*(\d+),`)
	matches := reg.FindStringSubmatch(body)
	if len(matches) > 0 {
		total, _ = strconv.Atoi(matches[1])
	}

	return total
}

// GetSogouInclude https://www.sogou.com/web?query=site%3Awww.baidu.com
func GetSogouInclude(serverHost string) (total int) {
	link := fmt.Sprintf("https://www.sogou.com/web?query=site%%3A%s", serverHost)

	req := gorequest.New().SetDoNotClearSuperAgent(true).TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).Timeout(5 * time.Second)
	//set key header
	req.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.81 Safari/537.36")
	_, body, errs := req.Get(link).End()
	if errs != nil {
		return total
	}

	// 'rn':'6051451',
	reg := regexp.MustCompile(`(?i)'rn':'(\d+)',`)
	matches := reg.FindStringSubmatch(body)
	if len(matches) > 0 {
		total, _ = strconv.Atoi(matches[1])
	}

	return total
}

func GetSoInclude(serverHost string) (total int) {
	link := fmt.Sprintf("https://www.so.com/s?q=site%%3A%s", serverHost)

	req := gorequest.New().SetDoNotClearSuperAgent(true).TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).Timeout(5 * time.Second)
	//set key header
	req.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.81 Safari/537.36")
	_, body, errs := req.Get(link).End()
	if errs != nil {
		return total
	}

	// 找到相关结果约79,000,000个
	reg := regexp.MustCompile(`(?i)找到相关结果约([0-9,]+)个`)
	matches := reg.FindStringSubmatch(body)
	if len(matches) > 0 {
		total, _ = strconv.Atoi(strings.ReplaceAll(matches[1], ",", ""))
	}

	return total
}

func GetBingInclude(serverHost string) (total int) {
	link := fmt.Sprintf("https://www.bing.com/search?q=site%%3A%s", serverHost)

	req := gorequest.New().SetDoNotClearSuperAgent(true).TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).Timeout(5 * time.Second)
	//set key header
	req.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.81 Safari/537.36")
	_, body, errs := req.Get(link).End()
	if errs != nil {
		return total
	}

	// 37,200 条结果
	reg := regexp.MustCompile(`(?i)([0-9,]+)\s*条结果`)
	matches := reg.FindStringSubmatch(body)
	if len(matches) > 0 {
		total, _ = strconv.Atoi(strings.ReplaceAll(matches[1], ",", ""))
	}

	return total
}

func GetGoogleInclude(serverHost string) (total int) {
	link := fmt.Sprintf("https://www.google.com/search?q=site%%3A%s", serverHost)

	req := gorequest.New().SetDoNotClearSuperAgent(true).TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).Timeout(5 * time.Second)
	//set key header
	req.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.81 Safari/537.36")
	_, body, errs := req.Get(link).End()
	if errs != nil {
		return total
	}

	// 37,200 条结果
	reg := regexp.MustCompile(`(?i)About\s([0-9,]+)\sresults`)
	matches := reg.FindStringSubmatch(body)
	if len(matches) > 0 {
		total, _ = strconv.Atoi(strings.ReplaceAll(matches[1], ",", ""))
	}

	return total
}
