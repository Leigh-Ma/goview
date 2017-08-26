package conf

import (
	"strings"
	"os"
	"github.com/astaxie/beego"
	"time"
	"net/url"
	"github.com/Unknwon/goconfig"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/beego/i18n"
	"path/filepath"
	"github.com/beego/compress"
	"github.com/astaxie/beego/cache"
	"github.com/howeyc/fsnotify"
)

const(
	DB_HOST = "dbhost"
	DB_PORT = "dbport"
	DB_NAME = "bdname"
	DB_USER = "dbuser"
	DB_PASSWORD = "dbpassword"
)


const (
	APP_VER = "0.1.0.1114"
)

var (
	AppName             string
	AppVer              string
	AppHost             string
	AppUrl              string
	AppLogo             string
	EnforceRedirect     bool
	AvatarURL           string
	SecretKey           string
	IsProMode           bool
	ActiveCodeLives     int
	ResetPwdCodeLives   int
	DateFormat          string
	DateTimeFormat      string
	DateTimeShortFormat string
	TimeZone            string
	RealtimeRenderMD    bool
	ImageSizeSmall      int
	ImageSizeMiddle     int
	ImageLinkAlphabets  []byte
	ImageXSend          bool
	ImageXSendHeader    string
	Langs               []string

	LoginRememberDays int
	LoginMaxRetries   int
	LoginFailedBlocks int

	CookieRememberName string
	CookieUserName     string
	TimeLimitCodeLength int

)


const (
	LangEnUS = iota
	LangZhCN
)


var (
	Cfg     *goconfig.ConfigFile
	Cache   cache.Cache
)

var (
	GlobalConfPath   = "conf/global/app.ini"
	AppConfPath      = "conf/app.ini"
	CompressConfPath = "conf/compress.json"
)

// LoadConfig loads configuration file.
func LoadConfig() *goconfig.ConfigFile {
	var err error

	if fh, _ := os.OpenFile(AppConfPath, os.O_RDONLY|os.O_CREATE, 0600); fh != nil {
		fh.Close()
	}

	// Load configuration, set app version and log level.
	Cfg, err = goconfig.LoadConfigFile(GlobalConfPath)

	if Cfg == nil {
		Cfg, err = goconfig.LoadConfigFile(AppConfPath)
		if err != nil {
			fmt.Println("Fail to load configuration file: " + err.Error())
			os.Exit(2)
		}

	} else {
		Cfg.AppendFiles(AppConfPath)
	}

	Cfg.BlockMode = false

	// set time zone of wetalk system
	TimeZone = Cfg.MustValue("app", "time_zone", "UTC")
	if _, err := time.LoadLocation(TimeZone); err == nil {
		os.Setenv("TZ", TimeZone)
	} else {
		fmt.Println("Wrong time_zone: " + TimeZone + " " + err.Error())
		os.Exit(2)
	}

	// Trim 4th part.
	AppVer = strings.Join(strings.Split(APP_VER, ".")[:3], ".")

	beego.BConfig.RunMode = Cfg.MustValue("app", "run_mode")
	beego.BConfig.Listen.HTTPPort = Cfg.MustInt("app", "http_port")

	IsProMode = beego.BConfig.RunMode == "pro"
	if IsProMode {
		beego.SetLevel(beego.LevelInformational)
	}

	// cache system
	Cache, err = cache.NewCache("memory", `{"interval":360}`)

	// session settings
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.Session.SessionProvider = Cfg.MustValue("session", "session_provider", "file")
	//beego.BConfig.WebConfig.Session.SessionSavePath = Cfg.MustValue("session", "session_path", "sessions")
	beego.BConfig.WebConfig.Session.SessionName = Cfg.MustValue("session", "session_name", "wetalk_sess")
	beego.BConfig.WebConfig.Session.SessionCookieLifeTime = Cfg.MustInt("session", "session_life_time", 0)
	beego.BConfig.WebConfig.Session.SessionGCMaxLifetime = Cfg.MustInt64("session", "session_gc_time", 86400)

	beego.BConfig.WebConfig.EnableXSRF = true
	// xsrf token expire time
	beego.BConfig.WebConfig.XSRFExpire = 86400 * 365

	driverName := Cfg.MustValue("orm", "driver_name", "mysql")
	dataSource := Cfg.MustValue("orm", "data_source", "root:123456@tcp(127.0.0.1:3306)/artist_work_dev?charset=utf8&loc=UTC")
	maxIdle := Cfg.MustInt("orm", "max_idle_conn", 30)
	maxOpen := Cfg.MustInt("orm", "max_open_conn", 50)

	// set default database
	err = orm.RegisterDataBase("default", driverName, dataSource, maxIdle, maxOpen)
	if err != nil {
		beego.Error(err)
	}

	orm.RunCommand()
	orm.Debug = true
	orm.DebugLog = orm.NewLog(os.Stdout)

	err = orm.RunSyncdb("default", false, false)
	if err != nil {
		beego.Error(err)
	}

	reloadConfig()

	settingLocales()
	settingCompress()

	configWatcher()

	return Cfg
}

func reloadConfig() {
	AppName = Cfg.MustValue("app", "app_name", "WeTalk Community")
	beego.BConfig.AppName = AppName

	AppHost = Cfg.MustValue("app", "app_host", "127.0.0.1:8080")
	AppUrl = Cfg.MustValue("app", "app_url", "http://127.0.0.1:8080/")
	AppLogo = Cfg.MustValue("app", "app_logo", "/static/img/logo.gif")
	AvatarURL = Cfg.MustValue("app", "avatar_url")

	EnforceRedirect = Cfg.MustBool("app", "enforce_redirect")

	DateFormat = Cfg.MustValue("app", "date_format")
	DateTimeFormat = Cfg.MustValue("app", "datetime_format")
	DateTimeShortFormat = Cfg.MustValue("app", "datetime_short_format")

	SecretKey = Cfg.MustValue("app", "secret_key")
	if len(SecretKey) == 0 {
		fmt.Println("Please set your secret_key in app.ini file")
	}

	ActiveCodeLives = Cfg.MustInt("app", "acitve_code_live_minutes", 180)
	ResetPwdCodeLives = Cfg.MustInt("app", "resetpwd_code_live_minutes", 180)

	LoginRememberDays = Cfg.MustInt("app", "login_remember_days", 7)
	LoginMaxRetries = Cfg.MustInt("app", "login_max_retries", 5)
	LoginFailedBlocks = Cfg.MustInt("app", "login_failed_blocks", 10)

	CookieRememberName = Cfg.MustValue("app", "cookie_remember_name", "wetalk_magic")
	CookieUserName = Cfg.MustValue("app", "cookie_user_name", "wetalk_powerful")

	RealtimeRenderMD = Cfg.MustBool("app", "realtime_render_markdown")

	ImageSizeSmall = Cfg.MustInt("image", "image_size_small")
	ImageSizeMiddle = Cfg.MustInt("image", "image_size_middle")

	if ImageSizeSmall <= 0 {
		ImageSizeSmall = 300
	}

	if ImageSizeMiddle <= ImageSizeSmall {
		ImageSizeMiddle = ImageSizeSmall + 400
	}

	str := Cfg.MustValue("image", "image_link_alphabets")
	if len(str) == 0 {
		str = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
	ImageLinkAlphabets = []byte(str)

	ImageXSend = Cfg.MustBool("image", "image_xsend", false)
	ImageXSendHeader = Cfg.MustValue("image", "image_xsend_header", "X-Accel-Redirect")

	orm.Debug = Cfg.MustBool("orm", "debug_log")

	// search setting
}



func settingCompress() {
	setting, err := compress.LoadJsonConf(CompressConfPath, IsProMode, AppUrl)
	if err != nil {
		beego.Error(err)
		return
	}

	setting.RunCommand()

	if IsProMode {
		setting.RunCompress(true, false, true)
	}

	beego.AddFuncMap("CompressJs", setting.Js.CompressJs)
	beego.AddFuncMap("CompressCss", setting.Css.CompressCss)
}

var eventTime = make(map[string]int64)

func configWatcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic("Failed start app watcher: " + err.Error())
	}

	go func() {
		for {
			select {
			case event := <-watcher.Event:
				switch filepath.Ext(event.Name) {
				case ".ini":
					if checkEventTime(event.Name) {
						continue
					}
					beego.Info(event)

					if err := Cfg.Reload(); err != nil {
						beego.Error("Conf Reload: ", err)
					}

					if err := i18n.ReloadLangs(); err != nil {
						beego.Error("Conf Reload: ", err)
					}

					reloadConfig()
					beego.Info("Config Reloaded")

				case ".json":
					if checkEventTime(event.Name) {
						continue
					}
					if event.Name == CompressConfPath {
						settingCompress()
						beego.Info("Beego Compress Reloaded")
					}
				}
			}
		}
	}()

	if err := watcher.WatchFlags("conf", fsnotify.FSN_MODIFY); err != nil {
		beego.Error(err)
	}
}

// checkEventTime returns true if FileModTime does not change.
func checkEventTime(name string) bool {
	mt := getFileModTime(name)
	if eventTime[name] == mt {
		return true
	}

	eventTime[name] = mt
	return false
}

// getFileModTime retuens unix timestamp of `os.File.ModTime` by given path.
func getFileModTime(path string) int64 {
	path = strings.Replace(path, "\\", "/", -1)
	f, err := os.Open(path)
	if err != nil {
		beego.Error("Fail to open file[ %s ]\n", err)
		return time.Now().Unix()
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		beego.Error("Fail to get file information[ %s ]\n", err)
		return time.Now().Unix()
	}

	return fi.ModTime().Unix()
}

func IsMatchHost(uri string) bool {
	if len(uri) == 0 {
		return false
	}

	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return false
	}

	if u.Host != AppHost {
		return false
	}

	return true
}