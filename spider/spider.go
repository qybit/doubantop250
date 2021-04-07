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
