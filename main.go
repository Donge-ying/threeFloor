package main

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
	_ "threeFloor/config"
	"threeFloor/floor"
	"threeFloor/utils"
	"time"
)

func main() {

	client := floor.NewClient(viper.GetString("account"), viper.GetString("pwd"))

	/*err := client.Login()
	if err != nil {
		panic(err)
	}*/

	client.SetKey(viper.GetString("key"))
	client.UserInfo.Id = viper.GetInt("user_id")
	err := client.GetUserInfo()

	if err != nil {
		panic(err)
	}

	client.SiginAll()
	return
	err = AutoReply(client, 5*time.Minute)
	fmt.Println(err)
}

func AutoReply(client *floor.Client, dur time.Duration) error {

	postList, err := client.GetMyPostList(10)

	if err != nil {
		return err
	}

	for _, v := range postList {
		fmt.Printf("[%d : %s] %s\n", v.Id, v.Category, v.Title)
	}

	fmt.Println("只展示前10个帖子，输入帖子ID来自动顶贴：")
	var postId int

	for _, err = fmt.Scanln(&postId); err != nil; {
		fmt.Println("输入错误：", err, "请重新输入：")
	}

	fmt.Println("输入成功，正在自动暖贴...")
	randComment := viper.GetStringSlice("rand_comment")

	//发帖
	for {
		comment := floor.Prefix + randComment[utils.Rand(0, len(randComment)-1)] + strings.Repeat("[玫瑰]", utils.Rand(0, 4))
		time.Sleep(time.Millisecond)
		comment += strings.Repeat("[爱心]", utils.Rand(0, 4))
		err = client.Comment(postId, 0, comment)

		if err != nil {
			if err == floor.ErrCommentRepeat {

				fmt.Println("评论重复", comment)
				time.Sleep(dur)
				continue

			}
			return err
		}

		fmt.Println("已回复：", comment)
		//获取所有评论列表
		commentList, n, err := client.GetPostComments(postId, 1)

		if err != nil {
			return err
		}

		fmt.Println("共有：", n, "页")
		if n > 1 {
			for i := 2; i <= n; i++ {
				cl, _, err := client.GetPostComments(postId, i)
				if err != nil {
					panic(err)
				}
				commentList = append(commentList, cl...)
			}
		}

		time.Sleep(dur)
		//查找我的帖子
		index := floor.GetLastComment(commentList, client.UserInfo.Id)

		if index == -1 {
			fmt.Println("没有发现回帖")
			break
		}

		myComment := commentList[index]
		fmt.Println("已发现我的回帖：", myComment.Text)

		//删除
		err = client.DeleteComment(myComment.Id)
		if err != nil {
			return err
		}

		fmt.Println("已删除回帖：", myComment.Text)
	}

	return nil
}
