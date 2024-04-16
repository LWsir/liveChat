package main

import (
	"bytes"
	"fmt"
	"liveChat/Utils"
	"liveChat/model"
	"net/http"
	"time"

	"github.com/dchest/captcha"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// 中间件，处理session
func Session(keyPairs string) gin.HandlerFunc {
	store := SessionConfig()
	return sessions.Sessions(keyPairs, store)
}
func SessionConfig() sessions.Store {
	sessionMaxAge := 3600
	sessionSecret := "topgoer"
	var store sessions.Store
	store = cookie.NewStore([]byte(sessionSecret))
	store.Options(sessions.Options{
		MaxAge: sessionMaxAge, //seconds
		Path:   "/",
	})
	return store
}

func Captcha(c *gin.Context, length ...int) {
	l := captcha.DefaultLen
	w, h := 107, 36
	if len(length) == 1 {
		l = length[0]
	}
	if len(length) == 2 {
		w = length[1]
	}
	if len(length) == 3 {
		h = length[2]
	}
	captchaId := captcha.NewLen(l)
	fmt.Println("验证码id", captchaId)
	session := sessions.Default(c)
	session.Set("captcha", captchaId)
	_ = session.Save()
	_ = Serve(c.Writer, c.Request, captchaId, ".png", "zh", false, w, h)
}
func CaptchaVerify(c *gin.Context, code string) bool {
	session := sessions.Default(c)
	if captchaId := session.Get("captcha"); captchaId != nil {
		session.Delete("captcha")
		_ = session.Save()
		if captcha.VerifyString(captchaId.(string), code) {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}
func Serve(w http.ResponseWriter, r *http.Request, id, ext, lang string, download bool, width, height int) error {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	var content bytes.Buffer
	switch ext {
	case ".png":
		w.Header().Set("Content-Type", "image/png")
		_ = captcha.WriteImage(&content, id, width, height)
	case ".wav":
		w.Header().Set("Content-Type", "audio/x-wav")
		_ = captcha.WriteAudio(&content, id, lang)
	default:
		return captcha.ErrNotFound
	}

	if download {
		w.Header().Set("Content-Type", "application/octet-stream")
	}
	http.ServeContent(w, r, id+ext, time.Time{}, bytes.NewReader(content.Bytes()))
	return nil
}
func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Use(Session("topgoer"))
	router.GET("/captcha.png", func(c *gin.Context) {
		Captcha(c, 4)
	})
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "这是登陆页",
		})

	})

	router.POST("/login", func(c *gin.Context) {
		userName := c.PostForm("username")
		password := c.PostForm("password")
		code := c.PostForm("code")
		isLogin := false

		if userName == "abner" && password == "123" {
			if CaptchaVerify(c, code) {
				isLogin = true
			}
		}

		if isLogin {
			c.HTML(http.StatusOK, "im.html", gin.H{
				"title": "这是聊天页",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": 1, "msg": "failed"})
		}
	})

	router.POST("/register", func(c *gin.Context) {
		//写一个注册
		type User struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		var newUser User

		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "参数校验失败",
				"data": "",
			})
			return
		}

		db_abner, err := Utils.InitDB("abner", "root", "123456", "127.0.0.1", 3306)

		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -2,
				"msg":  err.Error(),
				"data": "",
			})
			return
		}

		defer db_abner.Close()

		var userdb []model.UserDB

		// 转换 sql ，doris 比较特殊，用 gorm 的查询语句会报错Unsupported command(COM_STMT_PREPARE)
		sql := fmt.Sprintf("select * from t_users where user_name = '%s' ", newUser.Username)

		fmt.Printf(sql)

		err = db_abner.Raw(sql).Scan(&userdb).Error

		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -3,
				"msg":  err.Error(),
				"data": "",
			})
			return
		}

		if len(userdb) < 1 {
			insertUserDB := model.UserDB{
				UserName:    newUser.Username,
				Password:    newUser.Password,
				Email:       "",
				PhoneNumber: "",
				CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
				UpdatedAt:   time.Now().Format("2006-01-02 15:04:05"),
				State:       0,
			}
			result := db_abner.Table("t_users").Create(&insertUserDB)

			if result.Error != nil {
				c.JSON(http.StatusOK, gin.H{
					"code": -4,
					"msg":  result.Error.Error(),
					"data": "",
				})
				return
			}
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": -5,
				"msg":  newUser.Username + "已被注册",
				"data": userdb[0].PhoneNumber,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  newUser.Username + "注册成功",
			"data": "",
		})
		return

	})
	router.Run(":8083")
}
