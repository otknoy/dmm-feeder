package model

type SolrDocument struct {
	Add add `json:"add"`
}

type add struct {
	Doc Item `json:"doc"`
}

func NewSolrDocument(item Item) SolrDocument {
	d := SolrDocument{
		add{item},
	}
	return d
}
