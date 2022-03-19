package pipe

import "time"

type (
	GLApiSuccessfulStepsResponse []struct {
		ID             int         `json:"id"`
		Status         string      `json:"status"`
		Stage          string      `json:"stage"`
		Name           string      `json:"name"`
		Ref            string      `json:"ref"`
		Tag            bool        `json:"tag"`
		Coverage       interface{} `json:"coverage"`
		AllowFailure   bool        `json:"allow_failure"`
		CreatedAt      time.Time   `json:"created_at"`
		StartedAt      time.Time   `json:"started_at"`
		FinishedAt     time.Time   `json:"finished_at"`
		Duration       float64     `json:"duration"`
		QueuedDuration float64     `json:"queued_duration"`
		User           struct {
			ID              int         `json:"id"`
			Name            string      `json:"name"`
			Username        string      `json:"username"`
			State           string      `json:"state"`
			AvatarURL       string      `json:"avatar_url"`
			WebURL          string      `json:"web_url"`
			CreatedAt       time.Time   `json:"created_at"`
			Bio             string      `json:"bio"`
			BioHTML         string      `json:"bio_html"`
			Location        string      `json:"location"`
			PublicEmail     string      `json:"public_email"`
			Skype           string      `json:"skype"`
			Linkedin        string      `json:"linkedin"`
			Twitter         string      `json:"twitter"`
			WebsiteURL      string      `json:"website_url"`
			Organization    string      `json:"organization"`
			JobTitle        string      `json:"job_title"`
			Bot             bool        `json:"bot"`
			WorkInformation interface{} `json:"work_information"`
			Followers       int         `json:"followers"`
			Following       int         `json:"following"`
		} `json:"user"`
		Commit struct {
			ID             string    `json:"id"`
			ShortID        string    `json:"short_id"`
			CreatedAt      time.Time `json:"created_at"`
			ParentIds      []string  `json:"parent_ids"`
			Title          string    `json:"title"`
			Message        string    `json:"message"`
			AuthorName     string    `json:"author_name"`
			AuthorEmail    string    `json:"author_email"`
			AuthoredDate   time.Time `json:"authored_date"`
			CommitterName  string    `json:"committer_name"`
			CommitterEmail string    `json:"committer_email"`
			CommittedDate  time.Time `json:"committed_date"`
			WebURL         string    `json:"web_url"`
		} `json:"commit"`
		Pipeline struct {
			ID        int       `json:"id"`
			ProjectID int       `json:"project_id"`
			Sha       string    `json:"sha"`
			Ref       string    `json:"ref"`
			Status    string    `json:"status"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			WebURL    string    `json:"web_url"`
		} `json:"pipeline"`
		WebURL        string `json:"web_url"`
		ArtifactsFile struct {
			Filename string `json:"filename"`
			Size     int    `json:"size"`
		} `json:"artifacts_file,omitempty"`
		Artifacts []struct {
			FileType   string `json:"file_type"`
			Size       int    `json:"size"`
			Filename   string `json:"filename"`
			FileFormat string `json:"file_format"`
		} `json:"artifacts"`
		Runner struct {
			ID          int    `json:"id"`
			Description string `json:"description"`
			IPAddress   string `json:"ip_address"`
			Active      bool   `json:"active"`
			IsShared    bool   `json:"is_shared"`
			Name        string `json:"name"`
			Online      bool   `json:"online"`
			Status      string `json:"status"`
		} `json:"runner"`
		ArtifactsExpireAt time.Time `json:"artifacts_expire_at"`
		TagList           []string  `json:"tag_list"`
	}

	StepId struct {
		id   int
		name string
	}

	DownloadedArtifact struct {
		path string
		name string
	}
)
