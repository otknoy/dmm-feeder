package service

import (
	"github.com/otknoy/dmm-feeder/model"
)

type ItemProcessService struct {
}

func NewItemProcessService() ItemProcessService {
	return ItemProcessService{}
}

func (ips *ItemProcessService) Process(dmmItem model.DmmItem) model.Item {
	actresses := parseActress(dmmItem)
	genres := parseGenre(dmmItem)
	makers := parseMaker(dmmItem)

	item := model.Item{
		dmmItem.ContentID,
		dmmItem.Title,
		dmmItem.URL,
		dmmItem.ImageURL.Large,
		actresses,
		genres,
		makers,
	}

	return item
}

func parseActress(dmmItem model.DmmItem) []string {
	var actresses []string
	for i, actress := range dmmItem.Iteminfo.Actress {
		if i%3 != 0 {
			continue
		}
		actresses = append(actresses, actress.Name)
	}
	return actresses
}

func parseGenre(dmmItem model.DmmItem) []string {
	var genres []string
	for _, genre := range dmmItem.Iteminfo.Genre {
		genres = append(genres, genre.Name)
	}
	return genres
}

func parseMaker(dmmItem model.DmmItem) []string {
	var makers []string
	for _, maker := range dmmItem.Iteminfo.Maker {
		makers = append(makers, maker.Name)
	}
	return makers
}
