package web

import "time"

type WebhookEvent struct {
	Id                    string    `json:"id"`
	ActorId               string    `json:"actorId"`
	WebhookSubscriptionId string    `json:"webhookSubscriptionId"`
	CreatedAt             time.Time `json:"createdAt"`
	Event                 string    `json:"event"`
	Payload               struct {
		Id    string `json:"id"`
		Model struct {
			Id           string      `json:"id"`
			Name         string      `json:"name"`
			AvatarUrl    interface{} `json:"avatarUrl"`
			Color        string      `json:"color"`
			Role         string      `json:"role"`
			IsSuspended  bool        `json:"isSuspended"`
			CreatedAt    time.Time   `json:"createdAt"`
			UpdatedAt    time.Time   `json:"updatedAt"`
			DeletedAt    interface{} `json:"deletedAt"`
			LastActiveAt time.Time   `json:"lastActiveAt"`
			Timezone     string      `json:"timezone"`
		} `json:"model"`
	} `json:"payload"`
}
