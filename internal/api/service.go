package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
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
	Text interface{} `json:"text"`
}

type GqlResponse struct {
	Repository struct {
		Name        string `json:"name"`
		PullRequest struct {
			Nodes []struct {
				Title string `json:"title"`
				Url   string `json:"url"`
			} `json:"nodes"`
		} `json:"pullRequest"`
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
	gqlResp, err := f(s.token)
	if err != nil {
		log.Errorf("error happen %+v", err)
	}

	log.Infof("qq %+v", query["text"])

	result := Response{
		Text: gqlResp,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func f(token string) (*GqlResponse, error) {
	// create a client (safe to share across requests)
	client := graphql.NewClient("https://api.github.com/graphql")

	// make a request
	req := graphql.NewRequest(`
	query($owner: String!, $name: String!){
		repository (owner: $owner, name: $name) {
			name
		}
    }
`)

	// set any variables
	req.Var("owner", "octocat")
	req.Var("name", "Hello-World")
	req.Var("first", 5)
	req.Header.Set("Authorization", "Bearer "+token)

	// run it and capture the response
	var respData GqlResponse
	ctx := context.Background()

	if err := client.Run(ctx, req, &respData); err != nil {
		log.Info("req %+v\n", req)
		log.Fatal(err)
		return nil, err
	}

	return &respData, nil
}
