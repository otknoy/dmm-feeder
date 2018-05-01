package interfaces

import "github.com/otknoy/dmm-feeder/model"

type SolrAdder interface {
	Add(model.Item) error
}
