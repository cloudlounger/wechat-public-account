package model

type WXMessage struct {
	ToUserName   string       `json:"ToUserName" xml:"ToUserName"`
	FromUserName string       `json:"FromUserName" xml:"FromUserName"`
	CreateTime   int64        `json:"CreateTime" xml:"CreateTime"`
	MsgType      string       `json:"MsgType" xml:"MsgType"`
	Content      string       `json:"Content" xml:"Content"`
	Image        *WXImage     `json:"Image" xml:"Image"`
	Voice        *WXVoice     `json:"Voice" xml:"Voice"`
	Video        *WXVideo     `json:"Video" xml:"Video"`
	Music        *WXMusic     `json:"music" xml:"music"`
	Artical      []*WXArticle `json:"artical" xml:"artical"`
}

type WXImage struct {
	MediaId string `json:"MediaId" xml:"MediaId"`
}

type WXVoice struct {
	MediaId string `json:"MediaId" xml:"MediaId"`
}

type WXVideo struct {
	MediaId     string `json:"MediaId" xml:"MediaId"`
	Title       string `json:"Title" xml:"Title"`
	Description string `json:"Description" xml:"Description"`
}

type WXMusic struct {
	Title        string `json:"Title" xml:"Title"`
	Description  string `json:"Description" xml:"Description"`
	MusicUrl     string `json:"MusicUrl" xml:"MusicUrl"`
	HQMusicUrl   string `json:"HQMusicUrl" xml:"HQMusicUrl"`
	ThumbMediaId string `json:"ThumbMediaId" xml:"ThumbMediaId"`
}

type WXArticle struct {
	Title       string `json:"Title" xml:"Title"`
	Description string `json:"Description" xml:"Description"`
	PicUrl      string `json:"PicUrl" xml:"PicUrl"`
	Url         string `json:"Url" xml:"Url"`
}
