package pipe

import "time"

type DockerHubLoginResponse struct {
	Token string `json:"token"`
}

type DockerHubUpdateReadmeRequest struct {
	Description string `json:"description"`
	Readme      string `json:"full_description" validate:"required"`
}

type DockerHubUpdateReadmeResponse struct {
	User              string    `json:"user"`
	Name              string    `json:"name"`
	Namespace         string    `json:"namespace"`
	RepositoryType    string    `json:"repository_type"`
	Status            int       `json:"status"`
	Description       string    `json:"description"`
	IsPrivate         bool      `json:"is_private"`
	IsAutomated       bool      `json:"is_automated"`
	CanEdit           bool      `json:"can_edit"`
	StarCount         int       `json:"star_count"`
	PullCount         int       `json:"pull_count"`
	LastUpdated       time.Time `json:"last_updated"`
	IsMigrated        bool      `json:"is_migrated"`
	CollaboratorCount int       `json:"collaborator_count"`
	Affiliation       string    `json:"affiliation"`
	HubUser           string    `json:"hub_user"`
	HasStarred        bool      `json:"has_starred"`
	FullDescription   string    `json:"full_description"`
	Permissions       struct {
		Read  bool `json:"read"`
		Write bool `json:"write"`
		Admin bool `json:"admin"`
	} `json:"permissions"`
}
