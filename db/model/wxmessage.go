package model

type WXMessage struct {
	ToUserName   string       `json:"ToUserName"`
	FromUserName string       `json:"FromUserName"`
	CreateTime   int64        `json:"CreateTime"`
	MsgType      string       `json:"MsgType"`
	Content      string       `json:"Content"`
	Image        *WXImage     `json:"Image"`
	Voice        *WXVoice     `json:"Voice"`
	Video        *WXVideo     `json:"Video"`
	Music        *WXMusic     `json:"music"`
	Artical      []*WXArticle `json:"artical"`
}

type WXImage struct {
	MediaId string `json:"MediaId"`
}

type WXVoice struct {
	MediaId string `json:"MediaId"`
}

type WXVideo struct {
	MediaId     string `json:"MediaId"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}

type WXMusic struct {
	Title        string `json:"Title"`
	Description  string `json:"Description"`
	MusicUrl     string `json:"MusicUrl"`
	HQMusicUrl   string `json:"HQMusicUrl"`
	ThumbMediaId string `json:"ThumbMediaId"`
}

type WXArticle struct {
	Title       string `json:"Title"`
	Description string `json:"Description"`
	PicUrl      string `json:"PicUrl"`
	Url         string `json:"Url"`
}
