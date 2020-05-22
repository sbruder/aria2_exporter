package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/ybbus/jsonrpc"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	aria2TorrentPeers = prometheus.NewDesc(
		"aria2_torrent_peers",
		"Number of peers",
		[]string{"hash", "torrent"}, nil,
	)
	aria2TorrentSeeders = prometheus.NewDesc(
		"aria2_torrent_seeders",
		"Number of seeders",
		[]string{"hash", "torrent"}, nil,
	)
	aria2TorrentSize = prometheus.NewDesc(
		"aria2_torrent_size_bytes",
		"Size of torrent data",
		[]string{"hash", "torrent"}, nil,
	)
	aria2TorrentDownloaded = prometheus.NewDesc(
		"aria2_torrent_downloaded_bytes_total",
		"Amount of data downloaded",
		[]string{"hash", "torrent"}, nil,
	)
	aria2TorrentUploaded = prometheus.NewDesc(
		"aria2_torrent_uploaded_bytes_total",
		"Amount of data uploaded",
		[]string{"hash", "torrent"}, nil,
	)
)

type BittorrentInfo struct {
	Name string `json:"name"`
}

type Bittorrent struct {
	Info BittorrentInfo `json:"info"`
}

type Download struct {
	Bittorrent Bittorrent `json:"bittorrent"`
	Downloaded string     `json:"completedLength"`
	InfoHash   string     `json:"infoHash"`
	Peers      string     `json:"connections"`
	Seeders    string     `json:"numSeeders"`
	Size       string     `json:"totalLength"`
	Uploaded   string     `json:"uploadLength"`
}

type Downloads []Download

type Exporter struct {
	URL    string
	Secret string
}

func stringToFloat64(input string) (output float64) {
	output, err := strconv.ParseFloat(input, 64)
	if err != nil {
		log.Print(err)
	}

	return
}

func (e Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- aria2TorrentPeers
	ch <- aria2TorrentSeeders
	ch <- aria2TorrentSize
	ch <- aria2TorrentUploaded
	ch <- aria2TorrentDownloaded
}

func (e Exporter) Collect(ch chan<- prometheus.Metric) {
	rpcClient := jsonrpc.NewClient(fmt.Sprintf("%s/jsonrpc", e.URL))
	response, err := rpcClient.Call("aria2.tellActive", "token:"+e.Secret)
	if err != nil {
		log.Print(err)
		return
	}

	if response.Error != nil {
		log.Printf("JSON-RPC error has occurred: %s", response.Error)
		return
	}

	var downloads Downloads
	err = response.GetObject(&downloads)
	if err != nil {
		log.Print(err)
		return
	}

	for _, download := range downloads {
		if download.Bittorrent.Info.Name == "" {
			continue
		}
		hash := download.InfoHash
		name := download.Bittorrent.Info.Name
		peers := stringToFloat64(download.Peers)
		seeders := stringToFloat64(download.Seeders)
		size := stringToFloat64(download.Size)
		uploaded := stringToFloat64(download.Uploaded)
		downloaded := stringToFloat64(download.Downloaded)

		ch <- prometheus.MustNewConstMetric(aria2TorrentPeers, prometheus.GaugeValue, peers, hash, name)
		ch <- prometheus.MustNewConstMetric(aria2TorrentSeeders, prometheus.GaugeValue, seeders, hash, name)
		ch <- prometheus.MustNewConstMetric(aria2TorrentSize, prometheus.CounterValue, size, hash, name)
		ch <- prometheus.MustNewConstMetric(aria2TorrentUploaded, prometheus.CounterValue, uploaded, hash, name)
		ch <- prometheus.MustNewConstMetric(aria2TorrentDownloaded, prometheus.CounterValue, downloaded, hash, name)
	}
}

func main() {
	e := Exporter{URL: os.Getenv("ARIA2_URL"), Secret: os.Getenv("ARIA2_RPC_SECRET")}
	if e.URL == "" {
		log.Fatal("Please specify the environment variable ARIA2_URL")
	}
	if e.Secret == "" {
		log.Print("Specify RPC secret token using the environment variable ARIA2_RPC_SECRET if needed")
	}
	prometheus.MustRegister(e)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9578", nil))
}
