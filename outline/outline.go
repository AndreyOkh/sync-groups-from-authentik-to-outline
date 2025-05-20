package outline

// ol_api_ckcvWSsppU10MnXfYVbfR1EC5dXwTWKzOOzxGY

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Conf struct {
	URL   string
	Token string
}

type Client struct {
	URL   string
	Token string
}

func NewClient(conf Conf) *Client {
	return &Client{
		URL:   conf.URL,
		Token: conf.Token,
	}
}

func (client *Client) ListGroups(query string) (ListAllGroupsResponse, bool, error) {
	url := fmt.Sprintf("%s/api/%s", client.URL, "groups.list")

	req := ListAllGroupsRequest{
		Offset:    0,
		Limit:     100,
		Sort:      "updatedAt",
		Direction: "DESC",
	}
	if query != "" {
		req.Query = query
	}

	resp, err := sendRequest[ListAllGroupsResponse](url, client.Token, req)
	if err != nil || len(resp.Data.Groups) == 0 {
		return resp, false, err
	}
	return resp, true, nil
}

func (client *Client) ListAllGroupMembers(id string) (ListAllGroupMembersResponse, error) {
	url := fmt.Sprintf("%s/api/%s", client.URL, "groups.memberships")
	resp, err := sendRequest[ListAllGroupMembersResponse](url, client.Token, ListAllGroupMembersRequest{
		Id: id,
	})
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (client *Client) CreateGroup(groupName string) (string, error) {
	url := fmt.Sprintf("%s/api/%s", client.URL, "groups.create")
	resp, err := sendRequest[CreateGroupResponse](url, client.Token, CreateGroupRequest{
		Name: groupName,
	})
	if err != nil {
		return "", err
	}
	return resp.Data.Id, nil
}

func (client *Client) AddGroupMember(groupId string, userId string) error {
	url := fmt.Sprintf("%s/api/%s", client.URL, "groups.add_user")
	_, err := sendRequest[AddAGroupMemberResponse](url, client.Token, AddAGroupMemberRequest{
		Id:     groupId,
		UserId: userId,
	})
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) RemoveGroupMember(groupId string, userId string) error {
	url := fmt.Sprintf("%s/api/%s", client.URL, "groups.remove_user")
	_, err := sendRequest[RemoveAGroupMemberResponse](url, client.Token, RemoveAGroupMemberRequest{
		Id:     groupId,
		UserId: userId,
	})
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) ListUsers(query string) (ListAllUsersResponse, bool, error) {
	url := fmt.Sprintf("%s/api/%s", client.URL, "users.list")

	req := ListAllUsersRequest{
		Offset:    0,
		Limit:     100,
		Sort:      "updatedAt",
		Direction: "DESC",
	}
	if query != "" {
		req.Query = query
	}

	resp, err := sendRequest[ListAllUsersResponse](url, client.Token, req)
	if err != nil || len(resp.Data) == 0 {
		return resp, false, err
	}
	return resp, true, nil
}

func (client *Client) RetrieveUser(id string) (RetrieveAUserResponse, error) {
	url := fmt.Sprintf("%s/api/%s", client.URL, "users.info")
	resp, err := sendRequest[RetrieveAUserResponse](url, client.Token, RetrieveAUserRequest{
		Id: id,
	})
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func sendRequest[T any](baseUrl, token string, payload any) (T, error) {
	var body T
	var data bytes.Buffer
	err := json.NewEncoder(&data).Encode(payload)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("POST", baseUrl, &data)
	if err != nil {
		return body, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return body, err
	}

	if res.StatusCode != http.StatusOK {
		return body, errors.New(res.Status)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(res.Body)
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		log.Println(err)
	}
	return body, nil
}
