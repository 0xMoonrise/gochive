package handlers

type FilesResponse struct {
	Files any `json:"files"`
	Pages int `json:"pages"`
}

type Success struct {
	Status string `json:"status"`
}

type SuccesUpload struct {
	Success      bool `json:"success"`
	FileResponse `json:"file"`
}

type FileResponse struct {
	Id        int    `json:"id"`
	Filename  string `json:"filename"`
	Editorial string `json:"editorial"`
	Favorite  bool   `json:"favorite"`
}
