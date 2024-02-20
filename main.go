package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	//记录日志
	r.Use(gin.Logger())

	//捕获报错
	r.Use(gin.Recovery())

	store := cookie.NewStore([]byte("secret11111"))
	// 设置session中间件，参数mysession，指的是session的名字，也是cookie的名字
	// store是前面创建的存储引擎，我们可以替换成其他存储引擎
	r.Use(sessions.Sessions("mysession", store))

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "这是登陆页",
		})

	})

	// 验证码图片路由
	r.GET("/captcha.png", func(c *gin.Context) {

		captcha := generateCaptcha()
		c.Writer.Header().Set("Content-Type", "image/png")
		err := png.Encode(c.Writer, captcha)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to generate captcha")
			return
		}

	})

	r.Run(":8083")
}

func generateCaptcha() image.Image {
	// 创建一个 100x40 的图像
	img := image.NewRGBA(image.Rect(0, 0, 100, 40))

	// 填充背景色为白色
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	// 生成随机验证码
	rand.Seed(time.Now().UnixNano())
	captchaStr := strconv.Itoa(rand.Intn(9000) + 1000)

	// 将验证码绘制到图像上
	for i, ch := range captchaStr {
		drawCharacter(img, i*20+10, 20, string(ch))
	}

	return img
}

func drawCharacter(img *image.RGBA, x, y int, ch string) {
	// 使用随机颜色绘制字符
	rand.Seed(time.Now().UnixNano())
	c := color.RGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), 255}

	// 绘制字符到图像上
	for i := 0; i < 20; i++ {
		for j := 0; j < 20; j++ {
			img.Set(x+i, y+j, c)
		}
	}
}
