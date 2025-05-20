package main

import (
	"flag"
	"goauthentik.io/api/v3"
	"log"
	"os"
	"sync-groups-from-authentik-to-outline/authentik"
	"sync-groups-from-authentik-to-outline/config"
	"sync-groups-from-authentik-to-outline/outline"
	"sync-groups-from-authentik-to-outline/web"
)

type Group struct {
	Name         string
	UsersEmail   []string
	FriendlyName string
}

func main() {
	conf, err := config.Init()
	if err != nil {
		log.Printf("error initializing config: %v", err)
		os.Exit(1)
	}

	reSync := flag.Bool("force-resync", false, "Принудительное обновление всех групп и пользователей")
	flag.Parse()

	outlineClient := outline.NewClient(outline.Conf{
		URL:   conf.Outline.URL,
		Token: conf.Outline.Token,
	})

	authentikClient := authentik.NewClient(&authentik.Conf{
		Scheme: conf.Authentik.Scheme,
		Host:   conf.Authentik.Host,
		Token:  conf.Authentik.Token,
	})

	if *reSync {
		ForceResync(conf, outlineClient, authentikClient)
		os.Exit(0)
	}

	web.Webserver(&web.WebserverDeps{
		OClient: outlineClient,
		AClient: authentikClient,
		Config:  conf,
	})

}

func Difference(src, dst []string) (diff []string) {
	m := make(map[string]bool)
	for _, item := range src {
		m[item] = true
	}
	for _, item := range dst {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return
}

func GetGroupMembers(outlineClient *outline.Client, groupId string) ([]string, error) {
	var outlineGroupMembers = make([]string, 0)

	groupMembers, err := outlineClient.ListAllGroupMembers(groupId)
	if err != nil {
		return nil, err
	}

	for _, member := range groupMembers.Data.Users {
		email := GetUserEmail(outlineClient, member.Id)
		outlineGroupMembers = append(outlineGroupMembers, email)
	}
	return outlineGroupMembers, nil
}

func GetUserEmail(outlineClient *outline.Client, id string) string {
	user, err := outlineClient.RetrieveUser(id)
	if err != nil {
		log.Printf("error retrieving user: %v\n", err)
	}
	return user.Data.Email
}

func CreateGroup(outlineClient *outline.Client, groupName string) (string, error) {
	var groupId string

	groupResp, ok, err := outlineClient.ListGroups(groupName)
	if err != nil {
		return "", err
	}
	if !ok {
		groupId, err = outlineClient.CreateGroup(groupName)
		if err != nil {
			return "", err
		}
	} else {
		groupId = groupResp.Data.Groups[0].Id
	}
	return groupId, nil
}

// SeparateGroups функция очищает структуру от не нужных данных
func SeparateGroups(groups []api.Group, atr string) []Group {
	var result = make([]Group, 0)
	for _, g := range groups {
		if len(g.UsersObj) != 0 {
			var group Group
			group.Name = g.Name
			for _, user := range g.UsersObj {
				if *user.Email != "" {
					group.UsersEmail = append(group.UsersEmail, *user.Email)
				}
			}
			group.FriendlyName = g.Attributes[atr].(string)
			if len(group.UsersEmail) != 0 && group.FriendlyName != "" {
				result = append(result, group)
			}
		}
	}
	return result
}

// ForceResync функция обновляет все группы и их членов
func ForceResync(conf *config.Conf, outlineClient *outline.Client, authentikClient *authentik.Client) {
	groups := authentikClient.GetGroups(conf.App.GroupPrefix, "")
	g := SeparateGroups(groups, conf.App.GroupNameSelector)
	for _, group := range g {
		groupId, err := CreateGroup(outlineClient, group.FriendlyName)
		if err != nil {
			log.Printf("error creating group: %v\n", err)
		}

		outlineGroupMembers, err := GetGroupMembers(outlineClient, groupId)
		if err != nil {
			log.Printf("error getting group members: %v\n", err)
		}

		// Поиск новых членов группы отсутствующих в outline и добавление
		diff := Difference(outlineGroupMembers, group.UsersEmail)
		for _, user := range diff {
			id, ok, err := outlineClient.ListUsers(user)
			if err != nil {
				log.Printf("error listing users: %v\n", err)
				continue
			}
			if !ok {
				log.Printf("error listing user: %s not found\n", user)
				continue
			}
			err = outlineClient.AddGroupMember(groupId, id.Data[0].Id)
			if err != nil {
				log.Printf("error adding group member: %v\n", err)
			}
		}

		// Поиск отсутствующих членов группы и удаление
		diff = Difference(group.UsersEmail, outlineGroupMembers)
		for _, user := range diff {
			id, _, _ := outlineClient.ListUsers(user)
			err := outlineClient.RemoveGroupMember(groupId, id.Data[0].Id)
			if err != nil {
				log.Printf("error removing group member: %v\n", err)
			}
		}
	}
}
