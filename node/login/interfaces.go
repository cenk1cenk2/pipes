package login

type (
	NpmLoginJson struct {
		Username string `json:"username"           validate:"required"`
		Token    string `json:"token"              validate:"required"`
		Registry string `json:"registry,omitempty"                     default:"registry.npmjs.org"`
		UseHttps bool   `json:"useHttps,omitempty"                     default:"true"`
	}
)
