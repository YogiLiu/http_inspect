package main

import (
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"log/slog"
	"net"
	"net/http"
	"strings"
)

func search(db *geoip2.Reader, ip net.IP) (info, error) {
	i := info{}
	record, err := db.City(ip)
	if err != nil {
		return i, err
	}
	i.City = name(record.City.Names).Filter("en", "zh-CN")
	i.PostalCode = record.Postal.Code
	i.Continent.Names = name(record.Continent.Names).Filter("en", "zh-CN")
	i.Continent.Code = record.Continent.Code
	i.Subdivisions = make([]nameWithIso, 0, 1)
	for _, s := range record.Subdivisions {
		i.Subdivisions = append(i.Subdivisions, nameWithIso{
			Names:   name(s.Names).Filter("en", "zh-CN"),
			IsoCode: s.IsoCode,
		})
	}
	i.RepresentedRegion.Names = name(record.RepresentedCountry.Names).Filter("en", "zh-CN")
	i.RepresentedRegion.IsoCode = record.RepresentedCountry.IsoCode
	i.RepresentedRegion.IsInEuropeanUnion = record.RepresentedCountry.IsInEuropeanUnion
	i.Region.Names = name(record.Country.Names).Filter("en", "zh-CN")
	i.Region.IsoCode = record.Country.IsoCode
	i.Region.IsInEuropeanUnion = record.Country.IsInEuropeanUnion
	i.RegisteredRegion.Names = name(record.RegisteredCountry.Names).Filter("en", "zh-CN")
	i.RegisteredRegion.IsoCode = record.RegisteredCountry.IsoCode
	i.RegisteredRegion.IsInEuropeanUnion = record.RegisteredCountry.IsInEuropeanUnion
	i.Location.Latitude = record.Location.Latitude
	i.Location.Longitude = record.Location.Longitude
	i.Location.MetroCode = record.Location.MetroCode
	i.Location.AccuracyRadius = record.Location.AccuracyRadius
	i.Location.TimeZone = record.Location.TimeZone
	i.Traits.IsAnonymousProxy = record.Traits.IsAnonymousProxy
	i.Traits.IsAnycast = record.Traits.IsAnycast
	i.Traits.IsSatelliteProvider = record.Traits.IsSatelliteProvider
	return i, nil
}

type ipInfo struct {
	db *geoip2.Reader
}

type successRes struct {
	IP            string `json:"ip"`
	Version       int8   `json:"version"`
	IsPrivate     bool   `json:"isPrivate"`
	IsLoopback    bool   `json:"isLoopback"`
	IsMulticast   bool   `json:"isMulticast"`
	IsUnspecified bool   `json:"isUnspecified"`
	GeoInfo       info   `json:"geoInfo"`
}

type errorRes struct {
	Msg string `json:"msg"`
}

func (i ipInfo) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	value := req.URL.Query().Get("ip")
	if value == "" {
		value = req.Header.Get("X-Real-IP")
		if value != "" {
			slog.Info("found IP from X-Real-IP header", slog.String("ip", value))
		}
	} else {
		slog.Info("found IP from query string", slog.String("ip", value))
	}
	if value == "" {
		ips := req.Header.Get("X-Forwarded-For")
		if ips != "" {
			arr := strings.Split(ips, ",")
			value = arr[0]
			if value != "" {
				slog.Info("found IP from X-Forwarded-For header", slog.String("ip", value))
			}
		}
	}
	if value == "" {
		values := req.RemoteAddr
		v, _, err := net.SplitHostPort(values)
		if err == nil {
			slog.Info("found IP from RemoteAddr", slog.String("ip", v))
			value = v
		}
	}
	if value == "" {
		slog.Warn("can't detect IP address")
		writeRes(w, http.StatusBadRequest, errorRes{Msg: "can't detect IP address"})
		return
	}

	ip := net.ParseIP(value)
	if ip == nil {
		slog.Warn("invalid IP address", slog.String("ip", value))
		writeRes(w, http.StatusBadRequest, errorRes{Msg: "invalid IP address: " + value})
		return
	}

	r, err := search(i.db, ip)
	if err != nil {
		slog.Error("can't search IP", slog.String("err", err.Error()), slog.String("ip", value))
		writeRes(w, http.StatusNotFound, errorRes{Msg: "can't find IP: " + value})
		return
	}

	var version int8
	if ip.To4() != nil {
		version = 4
	} else if ip.To16() != nil {
		version = 6
	} else {
		version = 0
	}

	res := successRes{
		IP:            ip.String(),
		Version:       version,
		IsPrivate:     ip.IsPrivate(),
		IsLoopback:    ip.IsLoopback(),
		IsMulticast:   ip.IsMulticast(),
		IsUnspecified: ip.IsUnspecified(),
		GeoInfo:       r,
	}
	slog.Info("success to search", slog.String("info", fmt.Sprintf("%+v", res)))
	writeRes(w, http.StatusOK, res)
}

func notFound(w http.ResponseWriter, req *http.Request) {
	slog.Warn("not found", slog.String("method", req.Method), slog.String("path", req.URL.Path))
	writeRes(w, http.StatusNotFound, errorRes{Msg: "not found"})
}
