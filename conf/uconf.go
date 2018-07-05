package conf

type Configure struct {
	Gitpath      string `json:"Gitpath"`
	Pubdir       string `json:"Pubdir"`
	Prvdir       string `json:"Prvdir"`
	Prvpass      string `json:"Prvpass"`
	Tmpdir       string `json:"Tmpdir"`
	UploadPrefix string `json:"UploadPrefix"`
}
