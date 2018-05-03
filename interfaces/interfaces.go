package interfaces

import "github.com/otknoy/dmm-feeder/model"

type SolrIndexUpdater interface {
	AddItem(model.Item) error
	Commit() error
}
