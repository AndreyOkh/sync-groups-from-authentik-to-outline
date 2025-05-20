package web

import (
	"encoding/json"
	"fmt"
	"group-sync-from-authentik-to-outline/authentik"
	"group-sync-from-authentik-to-outline/outline"
	"io"
	"log"
	"net/http"
)

type Handler struct {
	OClient       *outline.Client
	AClient       *authentik.Client
	GroupPrefix   string
	GroupSelector string
}

func NewHandler(router *http.ServeMux, handler Handler) {
	router.HandleFunc("GET /health", handler.Health())
	router.HandleFunc("POST /", handler.Webhook())
}

func (h *Handler) Health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if isShuttingDown.Load() {
			http.Error(w, "Shutting down", http.StatusServiceUnavailable)
			return
		}
		if _, err := fmt.Fprintln(w, "OK"); err != nil {
			log.Println(err)
		}
	}
}

func (h *Handler) Webhook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if isShuttingDown.Load() {
			http.Error(w, "Shutting down", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Printf("error closing body: %v\n", err)
			}
		}(r.Body)

		event := WebhookEvent{}
		err := json.NewDecoder(r.Body).Decode(&event)
		if err != nil {
			log.Printf("error decoding json: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		userInfo := h.AClient.GetUserByName(event.Payload.Model.Name)
		fmt.Println(userInfo.Username)
		groups := h.AClient.GetGroups(h.GroupPrefix, userInfo.Username)
		for _, group := range groups {
			fmt.Println(group.Name)
			gName := group.Attributes[h.GroupSelector]
			oGroups, ok, err := h.OClient.ListGroups(gName.(string))
			if err != nil || !ok {
				log.Printf("error listing groups: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			groupId := oGroups.Data.Groups[0].Id
			if !ok {
				oGroup, err := h.OClient.CreateGroup(gName.(string))
				if err != nil {
					log.Printf("error creating group: %v\n", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				groupId = oGroup
			}
			if err := h.OClient.AddGroupMember(groupId, event.Payload.Model.Id); err != nil {
				log.Printf("error adding group member: %v\n", err)
			}
		}
	}
}
