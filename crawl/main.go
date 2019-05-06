package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	ROOT = "https://qiita.com/api/v2/"
)

type Article []struct {
	RenderedBody   string      `json:"rendered_body"`
	Body           string      `json:"body"`
	Coediting      bool        `json:"coediting"`
	CommentsCount  int         `json:"comments_count"`
	CreatedAt      time.Time   `json:"created_at"`
	Group          interface{} `json:"group"`
	ID             string      `json:"id"`
	LikesCount     int         `json:"likes_count"`
	Private        bool        `json:"private"`
	ReactionsCount int         `json:"reactions_count"`
	Tags           []struct {
		Name     string        `json:"name"`
		Versions []interface{} `json:"versions"`
	} `json:"tags"`
	Title     string    `json:"title"`
	UpdatedAt time.Time `json:"updated_at"`
	URL       string    `json:"url"`
	User      struct {
		Description       interface{} `json:"description"`
		FacebookID        interface{} `json:"facebook_id"`
		FolloweesCount    int         `json:"followees_count"`
		FollowersCount    int         `json:"followers_count"`
		GithubLoginName   string      `json:"github_login_name"`
		ID                string      `json:"id"`
		ItemsCount        int         `json:"items_count"`
		LinkedinID        interface{} `json:"linkedin_id"`
		Location          interface{} `json:"location"`
		Name              string      `json:"name"`
		Organization      interface{} `json:"organization"`
		PermanentID       int         `json:"permanent_id"`
		ProfileImageURL   string      `json:"profile_image_url"`
		TeamOnly          bool        `json:"team_only"`
		TwitterScreenName interface{} `json:"twitter_screen_name"`
		WebsiteURL        interface{} `json:"website_url"`
	} `json:"user"`
	PageViewsCount interface{} `json:"page_views_count"`
}

type tag struct {
	tagname string
	sum     string
}

func CsvReader(ch chan tag) {
	file1, err := os.Open(`./org_inv_tags.csv`)
	if err != nil {
		log.Fatal("Error:", err)
		return
	}
	defer file1.Close()
	reader := csv.NewReader(file1)
	reader.FieldsPerRecord = -1
	for {
		record, err := reader.Read()
		if err == io.EOF {
			close(ch)
			break
		} else if err != nil {
			log.Fatal("Error:", err)
			close(ch)
			break
		}
		log.Println("====", record[0], record[1], "====")
		if _, err := strconv.Atoi(record[1]); err == nil {
			ch <- tag{record[0], record[1]}
		}

	}
}

func ReadArticles(ch chan tag) {
	for {
		v, ok := <-ch
		if !ok {
			log.Println("Channnel Closed")
			break
		}
		i, err := strconv.Atoi(v.sum)
		if err != nil {
			log.Println("parse error")
			continue
		}
		page := 1
		for err == nil && i/100+10 > page {
			file, err := os.OpenFile("./dataset/"+v.tagname, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				log.Fatal(err)
			}
			//example https://qiita.com/api/v2/items?page=1&per_page=100&query=tag%3ARuby
			url := ROOT + "items?page=" + strconv.Itoa(page) + "&per_page=100&query=tag%3A" + v.tagname
			log.Println(url)
			req, _ := http.NewRequest("GET", url, nil)
			req.Header.Set("Authorization", "Bearer "+os.Getenv("qiita_token"))
			client := new(http.Client)
			resp, err := client.Do(req)

			if err != nil {
				fmt.Println(err)
				file.Close()
				resp.Body.Close()

				log.Fatal("Error:", err)
			}

			time.Sleep(time.Second * 4)

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				file.Close()
				log.Fatal("Error:", err)
				resp.Body.Close()
			}
			var articles []Article
			if err := json.Unmarshal(body, &articles); err != nil {
				//log.Println(string(body))
			}

			if len(articles) == 0 {
				file.Close()
				resp.Body.Close()
				break
			}
			file.Write(([]byte)(body))
			file.Write(([]byte)("\n"))

			page = page + 1
			file.Close()
			resp.Body.Close()
		}
		s3put("./dataset/"+v.tagname, v.tagname)
	}

}

func s3put(filepath, objname string) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	cre := credentials.NewStaticCredentials(
		os.Getenv("ACCESS_KEY"),
		os.Getenv("SECRET_KEY"),
		"")

	cli := s3.New(session.New(), &aws.Config{
		Credentials: cre,
		Region:      aws.String("ap-northeast-1"),
	})

	_, err = cli.PutObject(&s3.PutObjectInput{
		Bucket: aws.String("*"),
		Key:    aws.String(objname),
		Body:   file,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Object Uploaded:" + objname)
}

func main() {
	ch1 := make(chan tag)
	go CsvReader(ch1)
	ReadArticles(ch1)

	/*
		file, err := os.Create(`./tags.csv`)
		file.Write(([]byte)("name,count,follower count\n"))
		err = GetTags(os.Stdout, 100000000, file)
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()

		if err != nil {
			fmt.Println("Error happened", err)
		}
	*/
}
