package main

import (
	"bytes"
	"compress/zlib"
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"
)

func main() {
	r := gin.Default()
	r.Use(Cors()).GET(
		"/", downComment)
	r.Run("127.0.0.1:1188")
}

func Cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, x-token")
		context.Header("Access-Control-Allow-Methods", "*")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		context.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
		}
		context.Next()
	}
}

// 合并 map[string][]string
func mergeMaps(map1, map2 map[string][]string) map[string][]string {
	result := make(map[string][]string)
	for key, value := range map1 {
		result[key] = value
	}
	for key, value := range map2 {
		result[key] = value
	}
	return result
}

// 解压 zlib
func zlibDecode(compressedData []byte) map[string][]string {
	// 首先，使用zlib.NewReader创建一个读取器
	r, err := zlib.NewReader(bytes.NewBuffer(compressedData))
	if err != nil {
		panic(err)
	}
	defer r.Close()

	// 然后，读取解压后的数据
	uncompressedData, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}

	// 打印解压后的数据
	// fmt.Println(string(uncompressedData))
	return xmlDecode(uncompressedData)
}

type CommentList struct {
	XMLName xml.Name `xml:"danmu"`
	Text    string   `xml:",chardata"`
	Code    string   `xml:"code"`
	Data    struct {
		Text  string `xml:",chardata"`
		Entry []struct {
			Text string `xml:",chardata"`
			Int  string `xml:"int"`
			List struct {
				Text       string `xml:",chardata"`
				BulletInfo []struct {
					Text             string `xml:",chardata"`
					ContentId        string `xml:"contentId"`
					Content          string `xml:"content"`
					ParentId         string `xml:"parentId"`
					ShowTime         string `xml:"showTime"`
					Font             string `xml:"font"`
					Color            string `xml:"color"`
					Opacity          string `xml:"opacity"`
					Position         string `xml:"position"`
					Background       string `xml:"background"`
					VariableEffectId string `xml:"variableEffectId"`
					IsReply          string `xml:"isReply"`
					LikeCount        string `xml:"likeCount"`
					PlusCount        string `xml:"plusCount"`
					DissCount        string `xml:"dissCount"`
					IsShowLike       string `xml:"isShowLike"`
					IsShowLikeTest   string `xml:"isShowLikeTest"`
					IsShowReplyFlag  string `xml:"isShowReplyFlag"`
					ReplyCnt         string `xml:"replyCnt"`
					UserInfo         struct {
						Text           string `xml:",chardata"`
						SenderAvatar   string `xml:"senderAvatar"`
						Uid            string `xml:"uid"`
						Udid           string `xml:"udid"`
						Name           string `xml:"name"`
						AvatarId       string `xml:"avatarId"`
						AvatarVipLevel string `xml:"avatarVipLevel"`
						PicL           string `xml:"picL"`
					} `xml:"userInfo"`
					ContentType    string `xml:"contentType"`
					SubType        string `xml:"subType"`
					Src            string `xml:"src"`
					Spoiler        string `xml:"spoiler"`
					HalfScreenShow string `xml:"halfScreenShow"`
					ScoreLevel     string `xml:"scoreLevel"`
					EmotionType    string `xml:"emotionType"`
					MinVersion     struct {
						Text   string `xml:",chardata"`
						IPhone string `xml:"iPhone"`
						IPad   string `xml:"iPad"`
						GPhone string `xml:"GPhone"`
						GPad   string `xml:"GPad"`
					} `xml:"minVersion"`
					SpecialEffectType string `xml:"specialEffectType"`
					BtnType           string `xml:"btnType"`
					BtnExtJson        string `xml:"btnExtJson"`
				} `xml:"bulletInfo"`
			} `xml:"list"`
		} `xml:"entry"`
	} `xml:"data"`
	Sum      string `xml:"sum"`
	ValidSum string `xml:"validSum"`
	Duration string `xml:"duration"`
	Ts       string `xml:"ts"`
}

// 解析XML
func xmlDecode(data []byte) map[string][]string {
	// 解析XML
	var comment CommentList
	err := xml.Unmarshal(data, &comment)
	if err != nil {
		log.Fatalf("unable to unmarshal XML: %v", err)
	}
	// 应该给cap的，懒狗
	m := make(map[string][]string)
	// 打印解析后的数据
	for _, entry := range comment.Data.Entry {
		for _, info := range entry.List.BulletInfo {
			m["_"+info.ShowTime] = append(m["_"+info.ShowTime], info.Content)
		}
	}
	return m
}

func downComment(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.String(405, "id 不能为空")
	}
	// 获取tvid
	tvid := getTvid(id)

	idLength := len(tvid)
	// 协程 获取 弹幕，5分钟一个包 .z
	page := 30
	wg := sync.WaitGroup{}
	ch := make(chan map[string][]string, page)
	// 限制并发数
	reqNums := 10
	reqCh := make(chan struct{}, reqNums)
	wg.Add(page)
	for i := 1; i <= page; i++ {
		go func(i int) {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println("err:", err)
				}
				wg.Done()
				<-reqCh
			}()
			reqCh <- struct{}{}
			url := "https://cmts.iqiyi.com/bullet/" + string(tvid[idLength-4:idLength-2]) + "/" + tvid[idLength-2:] + "/" + tvid + "_300_" + strconv.Itoa(i) + ".z"
			fmt.Println("page:", url, "i:", i)
			// 下载
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println("page:" + strconv.Itoa(i) + "err:" + err.Error())
				return
			}

			res, err := io.ReadAll(resp.Body)
			fmt.Println("page:"+strconv.Itoa(i), ",StatusCode:", resp.StatusCode)
			if resp.StatusCode != 200 {
				return
			}

			if err != nil {
				panic(err)
			}
			m := zlibDecode(res)
			ch <- m
		}(i)
	}
	fmt.Println("等待协程结束")
	wg.Wait()
	close(ch)
	close(reqCh)
	fmt.Println("协程结束")
	ms := make(map[string][]string)
	for m := range ch {
		ms = mergeMaps(ms, m)
	}
	c.JSON(200, ms)
}

func getTvid(id string) string {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("err:", err)
		}
	}()

	jsUrl := "https://mesh.if.iqiyi.com/player/lw/lwplay/accelerator.js"
	referer := "https://www.iqiyi.com/" + id + ".html"
	fmt.Println(referer)
	client := &http.Client{
		Timeout: 2 * time.Second,
	}
	req, err := http.NewRequest("GET", jsUrl, bytes.NewBuffer([]byte("")))
	if err != nil {
		panic(err)
	}
	req.Header.Set("referer", referer)
	// req.Header.Set("Content-Type", "application/json")
	// 发送HTTP请求
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close() // 确保关闭响应体

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	r, _ := regexp.Compile(`"tvid":(\d+),`)
	tvidStr := r.FindString(string(body))

	reg := regexp.MustCompile(`\d+`)

	submatches := reg.FindSubmatch([]byte(tvidStr))

	tvid := string(submatches[0])

	fmt.Println("tvid:", tvid)

	return tvid

}
