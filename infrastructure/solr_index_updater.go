package infrastructure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/otknoy/dmm-feeder/interfaces"
	"github.com/otknoy/dmm-feeder/model"
)

type SolrIndexUpdater struct {
	solrHost       string
	solrPort       int
	solrCollection string
}

func NewSolrIndexUpdater(host string, port int, collection string) interfaces.SolrIndexUpdater {
	return &SolrIndexUpdater{host, port, collection}
}

func (siu *SolrIndexUpdater) AddItem(item model.Item) error {
	jsonStr, err := json.Marshal([]model.Item{item})
	if err != nil {
		return err
	}

	url := fmt.Sprintf(
		"http://%s:%d/solr/%s/update",
		siu.solrHost,
		siu.solrPort,
		siu.solrCollection,
	)

	res, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("err, http-status-code: %d", res.StatusCode)
	}

	return nil
}

func (siu *SolrIndexUpdater) Commit() error {
	url := fmt.Sprintf(
		"http://%s:%d/solr/%s/update?commit=true",
		siu.solrHost,
		siu.solrPort,
		siu.solrCollection,
	)

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("err, http-status-code: %d", res.StatusCode)
	}

	return nil
}
