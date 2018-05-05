package slack

type Payload struct {
	Text        string       `json:"text"`
	UserName    string       `json:"username"`
	IconEmoji   string       `json:"icon_emoji"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	Title     string            `json:"title"`
	TitleLink string            `json:"title_link"`
	Text      string            `json:"text"`
	Fallback  string            `json:"fallback"`
	Fields    []AttachmentField `json:"fields"`
}

type AttachmentField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}
