package outline

import "time"

// Groups

type RetrieveGroupRequest struct {
	Id string `json:"id"`
}

type RetrieveGroupResponse struct {
	Data struct {
		Id          string    `json:"id"`
		Name        string    `json:"name"`
		MemberCount int       `json:"memberCount"`
		CreatedAt   time.Time `json:"createdAt"`
		UpdatedAt   time.Time `json:"updatedAt"`
	} `json:"data"`
}

type ListAllGroupsRequest struct {
	Offset     int    `json:"offset,omitempty"`
	Limit      int    `json:"limit,omitempty"`
	Sort       string `json:"sort,omitempty"`
	Direction  string `json:"direction,omitempty"`
	UserId     string `json:"userId,omitempty"`
	ExternalId string `json:"externalId,omitempty"`
	Query      string `json:"query,omitempty"`
}

type ListAllGroupsResponse struct {
	Data struct {
		Groups []struct {
			Id          string    `json:"id"`
			Name        string    `json:"name"`
			MemberCount int       `json:"memberCount"`
			CreatedAt   time.Time `json:"createdAt"`
			UpdatedAt   time.Time `json:"updatedAt"`
		} `json:"groups"`
		GroupMemberships []struct {
			Id      string `json:"id"`
			GroupId string `json:"groupId"`
			UserId  string `json:"userId"`
			User    struct {
				Id           string    `json:"id"`
				Name         string    `json:"name"`
				AvatarUrl    string    `json:"avatarUrl"`
				Email        string    `json:"email"`
				Role         string    `json:"role"`
				IsSuspended  bool      `json:"isSuspended"`
				LastActiveAt time.Time `json:"lastActiveAt"`
				CreatedAt    time.Time `json:"createdAt"`
			} `json:"user"`
		} `json:"groupMemberships"`
	} `json:"data"`
	Pagination struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	} `json:"pagination"`
}

type ListAllGroupMembersRequest struct {
	Offset int    `json:"offset,omitempty"`
	Limit  int    `json:"limit,omitempty"`
	Id     string `json:"id,omitempty"`
	Query  string `json:"query,omitempty"`
}

type ListAllGroupMembersResponse struct {
	Data struct {
		Users []struct {
			Id           string    `json:"id"`
			Name         string    `json:"name"`
			AvatarUrl    string    `json:"avatarUrl"`
			Email        string    `json:"email"`
			Role         string    `json:"role"`
			IsSuspended  bool      `json:"isSuspended"`
			LastActiveAt time.Time `json:"lastActiveAt"`
			CreatedAt    time.Time `json:"createdAt"`
		} `json:"users"`
		GroupMemberships []struct {
			Id      string `json:"id"`
			GroupId string `json:"groupId"`
			UserId  string `json:"userId"`
			User    struct {
				Id           string    `json:"id"`
				Name         string    `json:"name"`
				AvatarUrl    string    `json:"avatarUrl"`
				Email        string    `json:"email"`
				Role         string    `json:"role"`
				IsSuspended  bool      `json:"isSuspended"`
				LastActiveAt time.Time `json:"lastActiveAt"`
				CreatedAt    time.Time `json:"createdAt"`
			} `json:"user"`
		} `json:"groupMemberships"`
	} `json:"data"`
	Pagination struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	} `json:"pagination"`
}

type CreateGroupRequest struct {
	Name string `json:"name"`
}

type CreateGroupResponse struct {
	Data struct {
		Id          string    `json:"id"`
		Name        string    `json:"name"`
		MemberCount int       `json:"memberCount"`
		CreatedAt   time.Time `json:"createdAt"`
		UpdatedAt   time.Time `json:"updatedAt"`
	} `json:"data"`
	Policies []struct {
		Id        string `json:"id"`
		Abilities struct {
			Create              bool `json:"create"`
			Read                bool `json:"read"`
			Update              bool `json:"update"`
			Delete              bool `json:"delete"`
			Restore             bool `json:"restore"`
			Star                bool `json:"star"`
			Unstar              bool `json:"unstar"`
			Share               bool `json:"share"`
			Download            bool `json:"download"`
			Pin                 bool `json:"pin"`
			Unpin               bool `json:"unpin"`
			Move                bool `json:"move"`
			Archive             bool `json:"archive"`
			Unarchive           bool `json:"unarchive"`
			CreateChildDocument bool `json:"createChildDocument"`
		} `json:"abilities"`
	} `json:"policies"`
}

type AddAGroupMemberRequest struct {
	Id     string `json:"id"`
	UserId string `json:"userId"`
}

type AddAGroupMemberResponse struct {
	Data struct {
		Users []struct {
			Id           string    `json:"id"`
			Name         string    `json:"name"`
			AvatarUrl    string    `json:"avatarUrl"`
			Email        string    `json:"email"`
			Role         string    `json:"role"`
			IsSuspended  bool      `json:"isSuspended"`
			LastActiveAt time.Time `json:"lastActiveAt"`
			CreatedAt    time.Time `json:"createdAt"`
		} `json:"users"`
		Groups []struct {
			Id          string    `json:"id"`
			Name        string    `json:"name"`
			MemberCount int       `json:"memberCount"`
			CreatedAt   time.Time `json:"createdAt"`
			UpdatedAt   time.Time `json:"updatedAt"`
		} `json:"groups"`
		GroupMemberships []struct {
			Id           string `json:"id"`
			UserId       string `json:"userId"`
			CollectionId string `json:"collectionId"`
			Permission   string `json:"permission"`
		} `json:"groupMemberships"`
	} `json:"data"`
}

type RemoveAGroupMemberRequest struct {
	Id     string `json:"id"`
	UserId string `json:"userId"`
}

type RemoveAGroupMemberResponse struct {
	Data struct {
		Groups []struct {
			Id          string    `json:"id"`
			Name        string    `json:"name"`
			MemberCount int       `json:"memberCount"`
			CreatedAt   time.Time `json:"createdAt"`
			UpdatedAt   time.Time `json:"updatedAt"`
		} `json:"groups"`
	} `json:"data"`
}

// Users

type RetrieveAUserRequest struct {
	Id string `json:"id"`
}

type RetrieveAUserResponse struct {
	Data struct {
		Id           string    `json:"id"`
		Name         string    `json:"name"`
		AvatarUrl    string    `json:"avatarUrl"`
		Email        string    `json:"email"`
		Role         string    `json:"role"`
		IsSuspended  bool      `json:"isSuspended"`
		LastActiveAt time.Time `json:"lastActiveAt"`
		CreatedAt    time.Time `json:"createdAt"`
	} `json:"data"`
}

type ListAllUsersRequest struct {
	Offset    int      `json:"offset,omitempty"`
	Limit     int      `json:"limit,omitempty"`
	Sort      string   `json:"sort,omitempty"`
	Direction string   `json:"direction,omitempty"`
	Query     string   `json:"query,omitempty"`
	Emails    []string `json:"emails,omitempty"`
	Filter    string   `json:"filter,omitempty"`
	Role      string   `json:"role,omitempty"`
}

type ListAllUsersResponse struct {
	Data []struct {
		Id           string    `json:"id"`
		Name         string    `json:"name"`
		AvatarUrl    string    `json:"avatarUrl"`
		Email        string    `json:"email"`
		Role         string    `json:"role"`
		IsSuspended  bool      `json:"isSuspended"`
		LastActiveAt time.Time `json:"lastActiveAt"`
		CreatedAt    time.Time `json:"createdAt"`
	} `json:"data"`
	Pagination struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	} `json:"pagination"`
}
