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
		[]string{"torrent"}, nil,
	)
	aria2TorrentSeeders = prometheus.NewDesc(
		"aria2_torrent_seeders",
		"Number of seeders",
		[]string{"torrent"}, nil,
	)
	aria2TorrentSize = prometheus.NewDesc(
		"aria2_torrent_size_bytes",
		"Size of torrent data",
		[]string{"torrent"}, nil,
	)
	aria2TorrentDownloaded = prometheus.NewDesc(
		"aria2_torrent_downloaded_bytes_total",
		"Amount of data downloaded",
		[]string{"torrent"}, nil,
	)
	aria2TorrentUploaded = prometheus.NewDesc(
		"aria2_torrent_uploaded_bytes_total",
		"Amount of data uploaded",
		[]string{"torrent"}, nil,
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
	Peers      string     `json:"connections"`
	Seeders    string     `json:"numSeeders"`
	Size       string     `json:"totalLength"`
	Downloaded string     `json:"completedLength"`
	Uploaded   string     `json:"uploadLength"`
}

type Downloads []Download

type Exporter struct {
	URL string
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
	response, err := rpcClient.Call("aria2.tellActive")
	if err != nil {
		log.Print(err)
	}
	var downloads Downloads
	err = response.GetObject(&downloads)
	if err != nil {
		log.Print(err)
	}

	for _, download := range downloads {
		name := download.Bittorrent.Info.Name
		peers := stringToFloat64(download.Peers)
		seeders := stringToFloat64(download.Seeders)
		size := stringToFloat64(download.Size)
		uploaded := stringToFloat64(download.Uploaded)
		downloaded := stringToFloat64(download.Downloaded)

		ch <- prometheus.MustNewConstMetric(aria2TorrentPeers, prometheus.GaugeValue, peers, name)
		ch <- prometheus.MustNewConstMetric(aria2TorrentSeeders, prometheus.GaugeValue, seeders, name)
		ch <- prometheus.MustNewConstMetric(aria2TorrentSize, prometheus.CounterValue, size, name)
		ch <- prometheus.MustNewConstMetric(aria2TorrentUploaded, prometheus.CounterValue, uploaded, name)
		ch <- prometheus.MustNewConstMetric(aria2TorrentDownloaded, prometheus.CounterValue, downloaded, name)
	}
}

func main() {
	e := Exporter{URL: os.Getenv("ARIA2_URL")}
	if e.URL == "" {
		log.Fatal("Please specify the environment variable ARIA2_URL")
	}
	prometheus.MustRegister(e)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9578", nil))
}
