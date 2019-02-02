package agent

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/pritunl/pritunl-zero/database"
	"github.com/pritunl/pritunl-zero/geo"
	"github.com/pritunl/pritunl-zero/node"
	"github.com/pritunl/pritunl-zero/settings"
	"github.com/ua-parser/uap-go/uaparser"
)

var (
	parser = uaparser.NewFromSaved()
)

const (
	Linux        = "linux"         // Linux = Debian + Linux + Ubuntu
	MacOs1010    = "macos_1010"    // macOS 10.10 = Mac OS X (10/10)
	MacOs1011    = "macos_1011"    // macOS 10.11 = Mac OS X (10/11)
	MacOs1012    = "macos_1012"    // macOS 10.12 = Mac OS X (10/12)
	MacOs1013    = "macos_1013"    // macOS 10.13 = Mac OS X (10/13)
	MacOs1014    = "macos_1014"    // macOS 10.14 = Mac OS X (10/14)
	WindowsXp    = "windows_xp"    // Windows XP = Windows XP
	Windows7     = "windows_7"     // Windows 7 = Windows 7
	WindowsVista = "windows_vista" // Windows Vista = Windows Vista
	Windows8     = "windows_8"     // Windows 8 = Windows 8 + Windows 8.1 + Windows RT 8.1
	Windows10    = "windows_10"    // Windows 10 = Windows 10
	ChromeOs     = "chrome_os"     // Chrome OS = Chrome OS
	Ios8         = "ios_8"         // iOS 8 = iOS (8/x)
	Ios9         = "ios_9"         // iOS 9 = iOS (9/x)
	Ios10        = "ios_10"        // iOS 10 = iOS (10/x)
	Ios11        = "ios_11"        // iOS 11 = iOS (11/x)
	Ios12        = "ios_12"        // iOS 12 = iOS (12/x)
	Android4     = "android_4"     // Android KitKat 4.4 = Android (4/4)
	Android5     = "android_5"     // Android Lollipop 5.0 = Android (5/x)
	Android6     = "android_6"     // Android Marshmallow 6.0 = Android (6/x)
	Android7     = "android_7"     // Android Nougat 7.0 = Android (7/x)
	Android8     = "android_8"     // Android Oreo 8.0 = Android (8/x)
	Android9     = "android_9"     // Android 9.0 = Android (9/x)
	Blackberry10 = "blackberry_10" // Blackerry 10 = BlackBerry OS (10/x)
	WindowsPhone = "windows_phone" // Windows Phone = Windows Phone
	FirefoxOs    = "firefox_os"    // Firefox OS = Firefox OS
	Kindle       = "kindle"        // Kindle = Kindle
)

const (
	Chrome                 = "chrome"                   // Chrome = Chrome + Chromium
	ChromeMobile           = "chrome_mobile"            // Chrome Mobile = Chrome Mobile + Chrome Mobile iOS + Chrome Mobile WebView
	Safari                 = "safari"                   // Safari = Safari
	SafariMobile           = "safari_mobile"            // Safari Mobile = Mobile Safari + Mobile Safari UI/WKWebView
	Firefox                = "firefox"                  // Firefox = Firefox + Firefox Beta
	FirefoxMobile          = "firefox_mobile"           // Firefox Mobile = Firefox Mobile + Firefox iOS
	Edge                   = "edge"                     // Microsoft Edge = Edge
	InternetExplorer       = "internet_explorer"        // Internet Explorer = IE
	InternetExplorerMobile = "internet_explorer_mobile" // Internet Explorer Mobile = IE Mobile
	Opera                  = "opera"                    // Opera = Opera
	OperaMobile            = "opera_mobile"             // Opera Mobile = Opera Mini + Opera Mobile + Opera Tablet + Opera Coast
)

type Agent struct {
	OperatingSystem string  `bson:"operating_system" json:"operating_system"`
	Browser         string  `bson:"browser" json:"browser"`
	Ip              string  `bson:"ip" json:"ip"`
	Isp             string  `bson:"isp" json:"isp"`
	Continent       string  `bson:"continent" json:"continent"`
	ContinentCode   string  `bson:"continent_code" json:"continent_code"`
	Country         string  `bson:"country" json:"country"`
	CountryCode     string  `bson:"country_code" json:"country_code"`
	Region          string  `bson:"region" json:"region"`
	RegionCode      string  `bson:"region_code" json:"region_code"`
	City            string  `bson:"city" json:"city"`
	Latitude        float64 `bson:"latitude" json:"latitude"`
	Longitude       float64 `bson:"longitude" json:"longitude"`
}

func Parse(db *database.Database, r *http.Request) (agnt *Agent, err error) {
	if settings.System.Demo {
		return
	}

	client := parser.Parse(r.UserAgent())

	ip := node.Self.GetRemoteAddr(r)

	ge, err := geo.Get(db, ip)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("agent: Failed to get geo IP information")
		err = nil
		return
	}

	agnt = &Agent{
		Ip:            ip,
		Isp:           ge.Isp,
		Continent:     ge.Continent,
		ContinentCode: ge.ContinentCode,
		Country:       ge.Country,
		CountryCode:   ge.CountryCode,
		Region:        ge.Region,
		RegionCode:    ge.RegionCode,
		City:          ge.City,
		Longitude:     ge.Longitude,
		Latitude:      ge.Latitude,
	}

	switch client.Os.Family {
	case "Android":
		switch client.Os.Major {
		case "4":
			if client.Os.Minor == "4" {
				agnt.OperatingSystem = Android4
				break
			}
			break
		case "5":
			agnt.OperatingSystem = Android5
			break
		case "6":
			agnt.OperatingSystem = Android6
			break
		case "7":
			agnt.OperatingSystem = Android7
			break
		case "8":
			agnt.OperatingSystem = Android8
			break
		case "9":
			agnt.OperatingSystem = Android9
			break
		}
		break
	case "BlackBerry OS":
		if client.Os.Major == "10" {
			agnt.OperatingSystem = Blackberry10
			break
		}
		break
	case "Firefox OS":
		agnt.OperatingSystem = FirefoxOs
		break
	case "iOS":
		switch client.Os.Major {
		case "8":
			agnt.OperatingSystem = Ios8
			break
		case "9":
			agnt.OperatingSystem = Ios9
			break
		case "10":
			agnt.OperatingSystem = Ios10
			break
		case "11":
			agnt.OperatingSystem = Ios11
			break
		case "12":
			agnt.OperatingSystem = Ios12
			break
		}
		break
	case "Kindle":
		agnt.OperatingSystem = Kindle
		break
	case "Mac OS X":
		if client.Os.Major == "10" {
			switch client.Os.Minor {
			case "10":
				agnt.OperatingSystem = MacOs1010
				break
			case "11":
				agnt.OperatingSystem = MacOs1011
				break
			case "12":
				agnt.OperatingSystem = MacOs1012
				break
			case "13":
				agnt.OperatingSystem = MacOs1013
				break
			case "14":
				agnt.OperatingSystem = MacOs1014
				break
			}
		}
		break
	case "Windows Phone":
		agnt.OperatingSystem = WindowsPhone
		break
	case "Windows XP":
		agnt.OperatingSystem = WindowsXp
		break
	case "Windows 7":
		agnt.OperatingSystem = Windows7
		break
	case "Windows Vista":
		agnt.OperatingSystem = WindowsVista
		break
	case "Windows 8", "Windows 8.1", "Windows RT 8.1":
		agnt.OperatingSystem = Windows8
		break
	case "Windows 10":
		agnt.OperatingSystem = Windows10
		break
	case "Chrome OS":
		agnt.OperatingSystem = ChromeOs
		break
	case "Linux", "Debian", "Ubuntu":
		agnt.OperatingSystem = Linux
		break
	}

	switch client.UserAgent.Family {
	case "Chrome", "Chromium":
		agnt.Browser = Chrome
		break
	case "Chrome Mobile", "Chrome Mobile iOS", "Chrome Mobile WebView":
		agnt.Browser = ChromeMobile
		break
	case "Safari":
		agnt.Browser = Safari
		break
	case "Mobile Safari", "Mobile Safari UI/WKWebView":
		agnt.Browser = SafariMobile
		break
	case "Firefox", "Firefox Beta":
		agnt.Browser = Firefox
		break
	case "Firefox Mobile", "Firefox iOS":
		agnt.Browser = FirefoxMobile
		break
	case "Edge":
		agnt.Browser = Edge
		break
	case "IE":
		agnt.Browser = InternetExplorer
		break
	case "IE Mobile":
		agnt.Browser = InternetExplorerMobile
		break
	case "Opera":
		agnt.Browser = Opera
		break
	case "Opera Mini", "Opera Mobile", "Opera Tablet", "Opera Coast":
		agnt.Browser = OperaMobile
		break
	}

	return
}

func (a *Agent) Diff(agnt *Agent) bool {
	if a.OperatingSystem != agnt.OperatingSystem ||
		a.Browser != agnt.Browser ||
		a.Ip != agnt.Ip ||
		a.Isp != agnt.Isp ||
		a.Continent != agnt.Continent ||
		a.ContinentCode != agnt.ContinentCode ||
		a.Country != agnt.Country ||
		a.CountryCode != agnt.CountryCode ||
		a.Region != agnt.Region ||
		a.RegionCode != agnt.RegionCode ||
		a.City != agnt.City ||
		a.Longitude != agnt.Longitude ||
		a.Latitude != agnt.Latitude {

		return true
	}

	return false
}
