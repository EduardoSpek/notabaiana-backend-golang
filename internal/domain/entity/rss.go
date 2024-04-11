package entity

import "encoding/xml"

type MediaContent struct {
	URL string `xml:"url,attr"`
}

type Item struct {
	Title       string       `xml:"title"`
	Description string       `xml:"description"`
	Link        string       `xml:"link"`
	Media       MediaContent `xml:"content"`
}

type Channel struct {
	Title string `xml:"title"`
	Items []Item `xml:"item"`
}

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}