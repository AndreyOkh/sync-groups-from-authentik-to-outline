package authentik

import (
	"context"
	"goauthentik.io/api/v3"
	"log"
)

type Conf struct {
	Scheme string
	Host   string
	Token  string
}

type Client struct {
	ApiClient   *api.APIClient
	AuthContext context.Context
}

func NewClient(conf *Conf) *Client {
	client := &Client{}
	apiConf := api.NewConfiguration()
	apiConf.Scheme = conf.Scheme
	apiConf.Host = conf.Host
	client.ApiClient = api.NewAPIClient(apiConf)
	client.AuthContext = context.WithValue(context.Background(), api.ContextAccessToken, conf.Token)
	return client
}

func (client *Client) GetGroups(filter, memberUsername string) []api.Group {
	var groups = make([]api.Group, 0)

	var page float32 = 1
	for page != 0 {
		req := client.ApiClient.CoreApi.CoreGroupsList(client.AuthContext).Search(filter).Page(int32(page))
		if memberUsername != "" {
			req = req.MembersByUsername([]string{memberUsername})
		}

		resp, r, err := req.Execute()
		if err != nil {
			log.Printf("Error when calling `CoreApi.CoreGroupsList``: %v\n", err)
			log.Printf("Full HTTP response: %v\n", r)
		}
		if r.StatusCode != 200 {
			break
		}
		for _, g := range resp.Results {
			groups = append(groups, g)
		}
		page = resp.Pagination.Next
	}

	return groups
}

// GetUserByEmail находит пользователя по e-mail и возвращает первое совпадение
func (client *Client) GetUserByEmail(userEmail string) api.User {
	resp, r, err := client.ApiClient.CoreApi.CoreUsersList(client.AuthContext).Email(userEmail).Execute()
	if err != nil {
		log.Printf("Error when calling `CoreApi.CoreUsersList``: %v\n", err)
	}
	if r.StatusCode != 200 {
		log.Printf("Error when calling `CoreApi.CoreUsersList``: %v\n", r.StatusCode)
	}
	return resp.Results[0]
}

// GetUserByName находит пользователя по имени и возвращает первое совпадение
func (client *Client) GetUserByName(name string) api.User {
	resp, r, err := client.ApiClient.CoreApi.CoreUsersList(client.AuthContext).Name(name).Execute()
	if err != nil {
		log.Printf("Error when calling `CoreApi.CoreUsersList``: %v\n", err)
	}
	if r.StatusCode != 200 {
		log.Printf("Error when calling `CoreApi.CoreUsersList``: %v\n", r.StatusCode)
	}
	return resp.Results[0]
}
