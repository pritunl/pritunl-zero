package policy

import (
	"github.com/ua-parser/uap-go/uaparser"
	"net/http"
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
	Android8     = "android_8"     // Android 8.0 = Android (8/x)
	Blackberry10 = "blackberry_10" // Blackerry 10 = BlackBerry OS (10/x)
	WindowsPhone = "windows_phone" // Windows Phone = Windows Phone
	FirefoxOs    = "firefox_os"    // Firefox OS = Firefox OS
	Kindle       = "kindle"        // Kindle = Kindle
)

func OperatingSystem(r *http.Request) string {
	client := parser.Parse(r.UserAgent())

	switch client.Os.Family {
	case "Android":
		switch client.Os.Major {
		case "4":
			if client.Os.Minor == "4" {
				return Android4
			}
			break
		case "5":
			return Android5
		case "6":
			return Android6
		case "7":
			return Android7
		case "8":
			return Android8
		}
		break
	case "BlackBerry OS":
		if client.Os.Major == "10" {
			return Blackberry10
		}
		break
	case "Firefox OS":
		return FirefoxOs
	case "iOS":
		switch client.Os.Major {
		case "8":
			return Ios8
		case "9":
			return Ios9
		case "10":
			return Ios10
		case "11":
			return Ios11
		case "12":
			return Ios12
		}
		break
	case "Kindle":
		return Kindle
	case "Mac OS X":
		if client.Os.Major == "10" {
			switch client.Os.Minor {
			case "10":
				return MacOs1010
			case "11":
				return MacOs1011
			case "12":
				return MacOs1012
			case "13":
				return MacOs1013
			}
		}
		break
	case "Windows Phone":
		return WindowsPhone
	case "Windows XP":
		return WindowsXp
	case "Windows 7":
		return Windows7
	case "Windows Vista":
		return WindowsVista
	case "Windows 8", "Windows 8.1", "Windows RT 8.1":
		return Windows8
	case "Windows 10":
		return Windows10
	case "Chrome OS":
		return ChromeOs
	case "Linux", "Debian", "Ubuntu":
		return Linux
	}

	return ""
}
