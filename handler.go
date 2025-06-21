package main

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strings"
)

type ipInfo struct {
}

type info struct {
	Country  string `json:"country"`
	Region   string `json:"region"`
	Province string `json:"province"`
	City     string `json:"city"`
	ISP      string `json:"isp"`
}

func newInfo(value string) info {
	i := info{}
	r, err := searcher.SearchByStr(value)
	if err != nil {
		slog.Error("can't search IP", slog.String("err", err.Error()), slog.String("ip", value))
		return i
	}
	arr := strings.Split(r, "|")
	if len(arr) != 5 {
		return i
	}
	if arr[0] != "0" {
		i.Country = arr[0]
	}
	if arr[1] != "0" {
		i.Region = arr[1]
	}
	if arr[2] != "0" {
		i.Province = arr[2]
	}
	if arr[3] != "0" {
		i.City = arr[3]
	}
	if arr[4] != "0" {
		i.ISP = arr[4]
	}
	return i
}

type successRes struct {
	IP   string `json:"ip"`
	Info info   `json:"info"`
}

type errorRes struct {
	Msg string `json:"msg"`
}

func (i ipInfo) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	value := req.URL.Query().Get("ip")
	if value == "" {
		value = req.Header.Get("X-Real-IP")
	} else {
		slog.Info("detected IP from query", slog.String("ip", value))
	}
	if value == "" {
		ips := req.Header.Get("X-Forwarded-For")
		if ips != "" {
			arr := strings.Split(ips, ",")
			value = arr[0]
		}
	} else {
		slog.Info("detected IP from X-Real-IP ", slog.String("ip", value))
	}
	if value == "" {
		values := req.RemoteAddr
		if strings.Contains(values, ":") {
			value = strings.Split(values, ":")[0]
		}
	} else {
		slog.Info("detected IP from X-Forwarded-For", slog.String("ip", value))
	}
	if value == "" {
		slog.Warn("can't detect IP address")
		writeRes(w, http.StatusBadRequest, errorRes{Msg: "can't detect IP address"})
		return
	} else {
		slog.Info("detected IP from RemoteAddr", slog.String("ip", value))
	}

	ip := net.ParseIP(value)
	if ip == nil {
		slog.Warn("invalid IP address", slog.String("ip", value))
		writeRes(w, http.StatusBadRequest, errorRes{Msg: "invalid IP address: " + value})
		return
	}

	if ip.To4() == nil {
		slog.Warn("IPv6 is not supported yet")
		writeRes(w, http.StatusBadRequest, errorRes{Msg: "IPv6 is not supported yet, IP: " + ip.String()})
		return
	}

	r := newInfo(ip.String())

	res := successRes{IP: ip.String(), Info: r}
	slog.Info("success to search", slog.String("info", fmt.Sprintf("%+v", res)))
	writeRes(w, http.StatusOK, res)
}

func notFound(w http.ResponseWriter, req *http.Request) {
	slog.Warn("not found", slog.String("method", req.Method), slog.String("path", req.URL.Path))
	writeRes(w, http.StatusNotFound, errorRes{Msg: "not found"})
}
