package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/machinebox/graphql"
	"github.com/prometheus/common/log"
)

type Service interface {
	command(http.ResponseWriter, *http.Request)
}

type service struct {
	timeout time.Duration
	token   string
}

type Response struct {
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	Text string `json:"text"`
}

type GqlResponse struct {
	Repository struct {
		Name        string `json:"name"`
		PullRequest struct {
			Nodes []struct {
				Title  string `json:"title"`
				Url    string `json:"url"`
				Author struct {
					Login string `json:"login"`
				} `json:"author"`
			} `json:"nodes"`
		} `json:"pullRequests"`
	} `json:"repository"`
}

func NewService(token string) Service {
	return &service{
		timeout: 10,
		token:   token,
	}
}

func (s *service) command(w http.ResponseWriter, r *http.Request) {

	if r.Body == nil {
		return
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("error happen %+v", err)
	}

	query, err := url.ParseQuery(string(body))
	if err != nil {
		log.Errorf("error happen %+v", err)
	}

	// get pr list from github api
	params := strings.Split(query["text"][0], " ")
	if len(params) < 3 {
		log.Errorf("params less than 3")
		return
	}

	limit, _ := strconv.Atoi(params[2])
	gqlResp, err := f(s.token, params[0], params[1], limit)
	if err != nil {
		log.Errorf("error happen %+v", err)
	}

	log.Infof("qq %+v", query["text"])

	repoInfo := gqlResp.Repository
	attachments := make([]Attachment, 0)
	for _, a := range repoInfo.PullRequest.Nodes {
		attachments = append(attachments, Attachment{
			Text: fmt.Sprintf("*%s* - %s %s", a.Author.Login, a.Title, a.Url),
		})
	}
	result := Response{
		Text:        "*Repo Name* : " + repoInfo.Name,
		Attachments: attachments,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func f(token, user, repo string, first int) (*GqlResponse, error) {
	// create a client (safe to share across requests)
	client := graphql.NewClient("https://api.github.com/graphql")

	// make a request
	req := graphql.NewRequest(`
	query($owner: String!, $name: String!, $first: Int!){
		repository (owner: $owner, name: $name) {
			name
			pullRequests(first:$first) {
				nodes {
					title
					url
					author {
						login
					}
				}
			}
		}
    }
`)

	// set any variables
	req.Var("owner", user)
	req.Var("name", repo)
	req.Var("first", first)
	req.Header.Set("Authorization", "Bearer "+token)

	// run it and capture the response
	var respData GqlResponse
	ctx := context.Background()

	if err := client.Run(ctx, req, &respData); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &respData, nil
}
