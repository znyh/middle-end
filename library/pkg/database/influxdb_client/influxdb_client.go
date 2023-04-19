package Influxdb_Client

import (
	"errors"
	"fmt"
	"strings"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
)

type Influxdb_Client struct {
	Conf client.HTTPConfig
	BP   client.BatchPointsConfig
}

func (c Influxdb_Client) NewClient() (client.Client, error) {
	db, err := client.NewHTTPClient(c.Conf)

	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	return db, err
}

func (c Influxdb_Client) NewBatchPoints() (client.BatchPoints, error) {
	bp, err := client.NewBatchPoints(c.BP)

	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	return bp, err
}

//插入序列
func (c Influxdb_Client) Insert(measurements string, tags map[string]string, fields map[string]interface{}, t ...time.Time) error {
	db, err := c.NewClient()
	if err != nil {
		return err
	}
	bp, err := c.NewBatchPoints()
	if err != nil {
		return err
	}
	defer db.Close()

	pt, err := client.NewPoint(measurements, tags, fields, t...)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return err
	}

	bp.AddPoint(pt)

	// Write the batch
	if err := db.Write(bp); err != nil {
		fmt.Println("Error: ", err.Error())
		return err
	}

	return nil
}

//查询
func (c Influxdb_Client) Query(command, database, precision string) ([]client.Result, error) {
	db, err := c.NewClient()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	if strings.Contains(command, "DELETE") {
		return nil, errors.New("can't execute 'DELETE' command")
	}

	q := client.NewQuery(command, database, precision)
	if response, err := db.Query(q); err == nil && response.Error() == nil {
		return response.Results, nil
	} else {
		if err == nil {
			return nil, response.Error()
		} else {
			return nil, err
		}
	}
}

//内部方法 暂不提供
func (c Influxdb_Client) Delete(measurements string) error {
	db, err := c.NewClient()
	if err != nil {
		return err
	}
	defer db.Close()

	return nil
}
