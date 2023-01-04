博客地址：[README](https://qybit.gitee.io/2021/04/07/doubantop250/)

## 前言

集中 学习+复习 Go 语言有一个多星期了，也该写点东西了。

说下使用的 Go 语言的感受吧，直观上来说。Go 刚上手是比较反人类的，比如变量名和变量类型的位置是反着的，已经函数和方法的返回值的位置就更加的奇怪了。但是，总体的上手难度不是很大。相反，我认为这些也是 Go 的一大特色吧。总体的学习门槛是比较低的，而且 Go 身上也有很多 C/C++ 的影子 (比如指针类型，还有结构体)。最让我感到意外的就是，Go 语言中的接口的设计，真正的做到了低耦合。因为只要任意一个结构体或者类型实现了接口中的方法后，就算是真正意义上的实现了一个接口。而当我们从代码里 "拿走" 这个接口时，是不会影响到实现了该接口的结构体或者类型，因为那只是它们的方法而已。

你可能会疑惑，为什么要从爬虫开始实践？

我认为兴趣是最好的老师，我不喜欢死板的去写 ”xxx通讯录管理系统“ 或者 ”xxx管理系统“ 之类的无聊 demo。我是兴趣驱动，我更愿意从爬虫入手去学习。

本来打算拿我的看家本领 Senlium 呢，结果到官网一查还不支持 Go。

好的，说的有点多了，下面开始我们的 Go 语言爬虫实践吧。

项目地址：[项目地址](https://github.com/qybit/doubantop250)

## 准备工作

开发环境

- go version go1.16.2 windows/amd64
- goland
- 第三方库 goquery（一个类似jQuery可以操作DOM的库）

## 知识点

- 结构体
- 函数 / 方法
- http
- 懂得 DOM 编程，至少会 JavaScript 中的 DOM 部分
- 正则表达式
- 异常处理
- json
- io 处理

## 工作目录

我们的工作目录长这样👇。

![](https://cdn.jsdelivr.net/gh/qybit/CDN@master/Photo/my/doubantop250_0.png)

## 发起http请求

作为一只合格的网络爬虫，我们必须要可以发起基本的 http 请求获得网页数据。

正规的网站，一般最基础的防御就是通过 User-Agent / Agent 字段的校验，来检测是不是真人用户操作。我们只需要在请求头中加入这一字段，把自己 "伪装" 成真人用户，具体内容可以在浏览器查看。这里不在赘述。

```go
// 获取内容
func fetchSinglePageContent(url string, start string) (io.Reader, error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", url+start, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.114 Safari/537.36")
	request.Header.Add("Referer", "https://movie.douban.com/top250")
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
```

## 使用获取的响应内容，构建 DOM 树

这一步，我们将会使用 goquery 库，它会帮助我们把网络请求的响应内容解析成一颗 DOM 树。并提供和 JavaScript 和 jQuery 类似的 API 供我们访问某个节点。

这里的 content 实际上就是 响应的Body 部分。

```go
// 解析获取的内容为 DOM 树
func generatorDomTree(content io.Reader) (*goquery.Document, error) {
	reader, err := goquery.NewDocumentFromReader(content)
	if err != nil {
		return nil, err
	}
	return reader, nil
}
```

## 解析单个页面的所有电影内容

我们拿到上面生成的 goquery.Document 对象， goquery 提供的选择器的功能和 jQuery 几乎一模一样。所以有过 jQuery 使用经验的话，上手 goquery 是非常容易的。这里的稍微麻烦点的就是我们处理字符串的时候。

```go
// 获取所有的电影对应的 li 标签
func parseSinglePage(doc *goquery.Document) ([]*entity.Movie, error) {
	var ret []*entity.Movie
	doc.Find("#content > div > div.article > ol > li").Each(func(i int, s *goquery.Selection) {
		cover, _ := s.Find(".pic a img").Eq(0).Attr("src")

		title := s.Find(".hd a span").Eq(0).Text()
		subtitle := s.Find(".hd a span").Eq(1).Text()
		subtitle = strings.TrimLeft(subtitle, "  /  ")

		other := s.Find(".hd a span").Eq(2).Text()
		other = strings.TrimLeft(other, "  /  ")

		desc := strings.TrimSpace(s.Find(".bd p").Eq(0).Text())
		DescInfo := strings.Split(desc, "\n")

		desc = DescInfo[0]

		movieDesc := strings.Split(DescInfo[1], "/")
		year := strings.TrimSpace(movieDesc[0])
		area := strings.TrimSpace(movieDesc[1])
		tag := strings.TrimSpace(movieDesc[2])

		star := s.Find(".bd star .rating_num").Text()

		comment := strings.TrimSpace(s.Find(".bd star span").Eq(3).Text())
		compile := regexp.MustCompile("[0-9]")
		comment = strings.Join(compile.FindAllString(comment, -1), "")

		quote := s.Find(".quote .inq").Text()

		movie := &entity.Movie{
			Title:    title,
			Subtitle: subtitle,
			Other:    other,
			Cover:    cover,
			Desc:     desc,
			Year:     year,
			Area:     area,
			Tag:      tag,
			Star:     star,
			Comment:  comment,
			Quote:    quote,
		}
		ret = append(ret, movie)
	})
	return ret, nil
}
```

## 汇总

这里将会调用上面的所有方法，然后根据用户输入的信息，进行决策。

比如从哪一页开始获取，以及是否需要持久化等

```go
// 解析单一页面
func parseOnePage(start string, page int, ok bool) {
   content, err := fetchSinglePageContent(URL, start)
   if err != nil {
      fmt.Println("获取内容时出错！")
      return
   }
   dom, err := generatorDomTree(content)
   if err != nil {
      fmt.Println("解析成 DOM 树的过程中出错！")
      return
   }
   books, err := parseSinglePage(dom)
   for _, book := range books {
      fmt.Println(book)
   }
   if ok {
      data, _ := json.Marshal(books)
      err := ioutil.WriteFile("page"+strconv.Itoa(page)+".txt", data, 0644)
      if err != nil {
         panic(err)
      }
   }
}
```

## 最终效果

![](https://cdn.jsdelivr.net/gh/qybit/CDN@master/Photo/my/toubantop250.png)

## 代码

### 依赖

```go
module qybit.com/doubantop250

go 1.16

require github.com/PuerkitoBio/goquery v1.6.1

```

### entiy 包

```go
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
```

### spider 包

```go
package spider

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"net/http"
	"qybit.com/doubantop250/entity"
	"regexp"
	"strconv"
	"strings"
)

const (
	URL          string = "https://movie.douban.com/top250?start="
	DefaultCover string = "https://img.imgdb.cn/item/601fdca33ffa7d37b326de61.jpg"
)


// 获取内容
func fetchSinglePageContent(url string, start string) (io.Reader, error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", url+start, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.114 Safari/537.36")
	request.Header.Add("Referer", "https://movie.douban.com/top250")
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// 获取所有的电影对应的 li 标签
func parseSinglePage(doc *goquery.Document) ([]*entity.Movie, error) {
	var ret []*entity.Movie
	doc.Find("#content > div > div.article > ol > li").Each(func(i int, s *goquery.Selection) {
		cover, _ := s.Find(".pic a img").Eq(0).Attr("src")

		title := s.Find(".hd a span").Eq(0).Text()
		subtitle := s.Find(".hd a span").Eq(1).Text()
		subtitle = strings.TrimLeft(subtitle, "  /  ")

		other := s.Find(".hd a span").Eq(2).Text()
		other = strings.TrimLeft(other, "  /  ")

		desc := strings.TrimSpace(s.Find(".bd p").Eq(0).Text())
		DescInfo := strings.Split(desc, "\n")

		desc = DescInfo[0]

		movieDesc := strings.Split(DescInfo[1], "/")
		year := strings.TrimSpace(movieDesc[0])
		area := strings.TrimSpace(movieDesc[1])
		tag := strings.TrimSpace(movieDesc[2])

		star := s.Find(".bd star .rating_num").Text()

		comment := strings.TrimSpace(s.Find(".bd star span").Eq(3).Text())
		compile := regexp.MustCompile("[0-9]")
		comment = strings.Join(compile.FindAllString(comment, -1), "")

		quote := s.Find(".quote .inq").Text()

		movie := &entity.Movie{
			Title:    title,
			Subtitle: subtitle,
			Other:    other,
			Cover:    cover,
			Desc:     desc,
			Year:     year,
			Area:     area,
			Tag:      tag,
			Star:     star,
			Comment:  comment,
			Quote:    quote,
		}
		ret = append(ret, movie)
	})
	return ret, nil
}

// 解析获取的内容为 DOM 树
func generatorDomTree(content io.Reader) (*goquery.Document, error) {
	reader, err := goquery.NewDocumentFromReader(content)
	if err != nil {
		return nil, err
	}
	return reader, nil
}

// 解析单一页面
func parseOnePage(start string, page int, ok bool) {
	content, err := fetchSinglePageContent(URL, start)
	if err != nil {
		fmt.Println("获取内容时出错！")
		return
	}
	dom, err := generatorDomTree(content)
	if err != nil {
		fmt.Println("解析成 DOM 树的过程中出错！")
		return
	}
	books, err := parseSinglePage(dom)
	for _, book := range books {
		fmt.Println(book)
	}
	if ok {
		data, _ := json.Marshal(books)
		err := ioutil.WriteFile("page"+strconv.Itoa(page)+".txt", data, 0644)
		if err != nil {
			panic(err)
		}
	}
}

func Run(page int, ok bool) {
	var k int = 0
	for i := 1; i <= page; i++ {
		parseOnePage(strconv.Itoa(k), i, ok)
		k += 25
	}
}

```

### app

```go
package main

import (
	"bufio"
	"fmt"
	"os"
	"qybit.com/doubantop250/spider"
	"strconv"
)

func main() {
	fmt.Print("请输入要爬取的页数，最大10页：")
	cin := bufio.NewScanner(os.Stdin)
	cin.Scan()
	page, err := strconv.Atoi(cin.Text())
	if err != nil {
		fmt.Println("输入数据不合法，请按照要求输入！")
		os.Exit(1)
	}
	fmt.Print("是否需要持久化？（请输入 yes/y 或者 no/n）")
	cin.Scan()
	ok := cin.Text()
	isOk := false
	if ok == "yes" || ok == "y" {
		isOk = true
	}
	spider.Run(page, isOk)
}

```

