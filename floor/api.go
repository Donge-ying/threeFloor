package floor

import (
	"errors"
	"fmt"
	"github.com/guonaihong/gout"
	"net/url"
	"strconv"
	"strings"
	"threeFloor/utils"
	"time"
)

var (
	ErrHttpCode      = errors.New("httpCode is not 200")
	ErrCommentRepeat = errors.New("请勿提交重复评论")
)

type Header map[string]string

type Client struct {
	Account    string
	Pwd        string
	Key        string
	SessionKey string
	UserInfo   UserInfo
}

func NewClient(account, pwd string) *Client {
	return &Client{
		Account: account,
		Pwd:     pwd,
	}
}

func (c *Client) SetKey(key string) {
	c.Key = key
}

func (c *Client) Login() error {
	var code int
	var m map[string]interface{}

	err := gout.POST("http://floor.huluxia.com/account/login/ANDROID/4.0?platform=2&gkey=000000&app_version=4.1.0.1.1&versioncode=20141443&market_id=floor_tencent&_key=&device_code=%5Bw%5D02%3A00%3A00%3A00%3A00%3A00%5Bd%5D1806d476-1533-45fd-8fff-5262ade07222").
		SetBody("account=" + url.QueryEscape(c.Account) + "&login_type=2&password=" + utils.HexMd5(c.Pwd)).
		SetTimeout(5 * time.Second).
		BindJSON(&m).
		SetHeader(Header{"Content-Type": "application/x-www-form-urlencoded", "User-Agent": "okhttp/3.8.1"}).
		Code(&code).
		Do()

	if err != nil {
		return err
	}

	if code != 200 {
		return ErrHttpCode
	}

	//判断是否登录成功
	status := int(m["status"].(float64))
	if status != 1 {
		return errors.New(m["msg"].(string))
	}

	c.Key = m["_key"].(string)
	c.SessionKey = m["session_key"].(string)
	userInfo := m["user"].(map[string]interface{})
	c.UserInfo.Age = int(userInfo["age"].(float64))
	c.UserInfo.Avatar = userInfo["avatar"].(string)
	c.UserInfo.Birthday = int(userInfo["birthday"].(float64))
	c.UserInfo.Gender = int(userInfo["gender"].(float64))
	c.UserInfo.Id = int(userInfo["userID"].(float64))
	c.UserInfo.Nick = userInfo["nick"].(string)
	c.UserInfo.Role = int(userInfo["role"].(float64))
	//fmt.Printf("%+v", c)
	return nil
}
func (c *Client) GetUserInfo() error {

	var code int
	var m map[string]interface{}

	err := gout.GET("http://floor.huluxia.com/user/info/ANDROID/2.1?platform=2&gkey=000000&app_version=4.1.0.1.1&versioncode=20141443&market_id=floor_tencent&_key=" + c.Key + "&device_code=%5Bw%5D02%3A00%3A00%3A00%3A00%3A00%5Bd%5D1806d476-1533-45fd-8fff-5262ade07222&user_id=" + strconv.Itoa(c.UserInfo.Id)).
		SetTimeout(5 * time.Second).
		//Debug(true).
		BindJSON(&m).
		SetHeader(Header{"User-Agent": "okhttp/3.8.1"}).
		Code(&code).
		Do()

	if err != nil {
		return err
	}

	if code != 200 {
		return ErrHttpCode
	}

	status := int(m["status"].(float64))

	if status != 1 {
		return errors.New(m["msg"].(string))
	}

	userInfo := m
	c.UserInfo.Age = int(userInfo["age"].(float64))
	c.UserInfo.Avatar = userInfo["avatar"].(string)
	c.UserInfo.Birthday = int(userInfo["birthday"].(float64))
	c.UserInfo.Gender = int(userInfo["gender"].(float64))
	c.UserInfo.Id = int(userInfo["userID"].(float64))
	c.UserInfo.Nick = userInfo["nick"].(string)
	c.UserInfo.Role = int(userInfo["role"].(float64))
	return nil
}

func (c *Client) GetMyPostList(count int) ([]Post, error) {
	var code int
	var m map[string]interface{}

	err := gout.GET("http://floor.huluxia.com/post/create/list/ANDROID/2.0?platform=2&gkey=000000&app_version=4.1.0.1.1&versioncode=20141443&market_id=floor_tencent&_key=" + c.Key + "&device_code=%5Bw%5D02%3A00%3A00%3A00%3A00%3A00%5Bd%5D1806d476-1533-45fd-8fff-5262ade07222&start=0&count=" + strconv.Itoa(count) + "&user_id=" + strconv.Itoa(c.UserInfo.Id)).
		SetTimeout(5 * time.Second).
		//Debug(true).
		BindJSON(&m).
		SetHeader(Header{"User-Agent": "okhttp/3.8.1"}).
		Code(&code).
		Do()

	if err != nil {
		return nil, err
	}

	if code != 200 {
		return nil, ErrHttpCode
	}

	status := int(m["status"].(float64))

	if status != 1 {
		return nil, errors.New(m["msg"].(string))
	}
	//解析
	post := m["posts"].([]interface{})

	postlist := make([]Post, len(post))
	for i, v := range post {
		obj := v.(map[string]interface{})
		postlist[i].Id = int(obj["postID"].(float64))
		postlist[i].Title = obj["title"].(string)
		postlist[i].Category = obj["category"].(map[string]interface{})["title"].(string)
	}

	return postlist, nil
}

func (c *Client) GetCategoryList() ([]Category, error) {

	var code int
	var m map[string]interface{}

	err := gout.GET("http://floor.huluxia.com/category/list/ANDROID/2.0?platform=2&gkey=000000&app_version=4.1.0.1.1&versioncode=20141443&market_id=floor_tencent&_key=" + c.Key + "&device_code=%5Bw%5D02%3A00%3A00%3A00%3A00%3A00%5Bd%5D1806d476-1533-45fd-8fff-5262ade07222&is_hidden=1").
		SetTimeout(5 * time.Second).
		//Debug(true).
		BindJSON(&m).
		SetHeader(Header{"User-Agent": "okhttp/3.8.1"}).
		Code(&code).
		Do()

	if err != nil {
		return nil, err
	}
	if code != 200 {
		return nil, ErrHttpCode
	}
	status := int(m["status"].(float64))
	if status != 1 {
		return nil, errors.New(m["msg"].(string))
	}
	//解析
	categorys := m["categories"].([]interface{})
	categoryList := make([]Category, len(categorys))

	for i, v := range categorys {
		obj := v.(map[string]interface{})
		categoryList[i].Id = int(obj["categoryID"].(float64))
		categoryList[i].Title = obj["title"].(string)
		categoryList[i].IsSubscribe = int(obj["isSubscribe"].(float64))
	}

	return categoryList, err

}

func (c *Client) Signin(categoryId int) error {
	var code int
	var m map[string]interface{}

	err := gout.GET("http://floor.huluxia.com/user/signin/ANDROID/4.0?platform=2&gkey=000000&app_version=4.1.0.1.1&versioncode=20141443&market_id=floor_tencent&_key=" + c.Key + "&device_code=%5Bw%5D02%3A00%3A00%3A00%3A00%3A00%5Bd%5D1806d476-1533-45fd-8fff-5262ade07222&cat_id=" + strconv.Itoa(categoryId)).
		SetTimeout(5 * time.Second).
		//Debug(true).
		BindJSON(&m).
		SetHeader(Header{"User-Agent": "okhttp/3.8.1"}).
		Code(&code).
		Do()

	if err != nil {
		return err
	}
	if code != 200 {
		return ErrHttpCode
	}
	status := int(m["status"].(float64))
	if status != 1 {
		return errors.New(m["msg"].(string))
	}
	return nil
}

func (c *Client) SiginAll() error {
	//获取列表
	list, err := c.GetCategoryList()
	if err != nil {
		return err
	}
	for _, v := range list {
		if v.Id != 0 {
			//签到
			err := c.Signin(v.Id)
			if err != nil {
				fmt.Printf("[%d: %s] 签到失败：%s\n", v.Id, v.Title, err.Error())
			} else {
				fmt.Printf("[%d: %s] 签到成功\n", v.Id, v.Title)
			}
			time.Sleep(500 * time.Millisecond)
		}
	}
	fmt.Println("全部签到完成！")
	return nil
}

//commentId = 0时回复楼主
func (c *Client) Comment(postId int, commentId int, text string) error {
	var code int
	var m map[string]interface{}

	err := gout.POST("http://floor.huluxia.com/comment/create/ANDROID/2.0?platform=2&gkey=000000&app_version=4.1.0.1.1&versioncode=20141443&market_id=floor_tencent&_key=" + c.Key + "&device_code=%5Bw%5D02%3A00%3A00%3A00%3A00%3A00%5Bd%5D1806d476-1533-45fd-8fff-5262ade07222").
		SetBody("post_id=" + strconv.Itoa(postId) + "&comment_id=" + strconv.Itoa(commentId) + "&text=" + url.QueryEscape(text) + "&patcha=&images=&remindUsers=").
		SetTimeout(5 * time.Second).
		//Debug(true).
		BindJSON(&m).
		SetHeader(Header{"Content-Type": "application/x-www-form-urlencoded", "User-Agent": "okhttp/3.8.1"}).
		Code(&code).
		Do()

	if err != nil {
		return err
	}
	if code != 200 {
		return ErrHttpCode
	}
	//判断是否成功
	status := int(m["status"].(float64))
	if status != 1 {
		code := int(m["code"].(float64))
		if code == 104 {
			return ErrCommentRepeat
		}
		return errors.New(m["msg"].(string))
	}
	return nil
}

func (c *Client) DeleteComment(commentId int) error {
	var code int
	var m map[string]interface{}

	err := gout.GET("http://floor.huluxia.com/comment/destroy/ANDROID/2.0?comment_id=" + strconv.Itoa(commentId) + "&platform=2&gkey=000000&app_version=4.1.0.1.1&versioncode=20141443&market_id=floor_tencent&_key=" + c.Key + "&device_code=%5Bw%5D02%3A00%3A00%3A00%3A00%3A00%5Bd%5D1806d476-1533-45fd-8fff-5262ade07222").
		SetTimeout(5 * time.Second).
		//Debug(true).
		BindJSON(&m).
		SetHeader(Header{"User-Agent": "okhttp/3.8.1"}).
		Code(&code).
		Do()

	if err != nil {
		return err
	}
	if code != 200 {
		return ErrHttpCode
	}
	status := int(m["status"].(float64))
	if status != 1 {
		return errors.New(m["msg"].(string))
	}
	return nil
}

func (c *Client) GetPostComments(postId int, pageNo int) ([]Comment, int, error) {
	var code int
	var m map[string]interface{}
	var totalPage int
	err := gout.GET("http://floor.huluxia.com/post/detail/ANDROID/2.3?platform=2&gkey=000000&app_version=4.1.0.1.1&versioncode=20141443&market_id=floor_tencent&_key=" + c.Key + "&device_code=%5Bw%5D02%3A00%3A00%3A00%3A00%3A00%5Bd%5D1806d476-1533-45fd-8fff-5262ade07222&post_id=" + strconv.Itoa(postId) + "&page_no=" + strconv.Itoa(pageNo) + "&page_size=20&doc=1").
		SetTimeout(5 * time.Second).
		//Debug(true).
		BindJSON(&m).
		SetHeader(Header{"User-Agent": "okhttp/3.8.1"}).
		Code(&code).
		Do()

	if err != nil {
		return nil, totalPage, err
	}
	if code != 200 {
		return nil, totalPage, ErrHttpCode
	}
	status := int(m["status"].(float64))
	if status != 1 {
		return nil, totalPage, errors.New(m["msg"].(string))
	}
	//解析
	totalPage = int(m["totalPage"].(float64))
	comments := m["comments"].([]interface{})
	commentList := make([]Comment, len(comments))
	for i, v := range comments {
		obj := v.(map[string]interface{})
		commentList[i].Id = int(obj["commentID"].(float64))
		commentList[i].Text = obj["text"].(string)
		userInfo := obj["user"].(map[string]interface{})
		commentList[i].UserInfo.Id = int(userInfo["userID"].(float64))
		commentList[i].UserInfo.Age = int(userInfo["age"].(float64))
		commentList[i].UserInfo.Avatar = userInfo["avatar"].(string)
		commentList[i].UserInfo.Gender = int(userInfo["gender"].(float64))
		commentList[i].UserInfo.Nick = userInfo["nick"].(string)
		commentList[i].UserInfo.Role = int(userInfo["role"].(float64))
	}
	return commentList, totalPage, nil
}

func GetLastComment(comments []Comment, userId int) int {
	for i := len(comments) - 1; i >= 0; i-- {
		if comments[i].UserInfo.Id == userId && strings.HasPrefix(comments[i].Text, Prefix) {
			return i
		}
	}
	return -1
}
