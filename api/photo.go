package api

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
	"github.com/sky/api/basic"
	"github.com/sky/util"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"strconv"
	"strings"
)

// 上传响应数据
type UpdateDate struct {
	Size  int64  `json:"size"`  // 大小
	Mime  string `json:"mime"`  // 图片类型
	Imgid string `json:"imgid"` // 图片id
}

// 图像类型
const (
	PNG  = "png"
	JPG  = "jpg"
	JPEG = "jpeg"
	GIF  = "gif"
)

var imgType []string = []string{PNG, JPG, JPEG, GIF}

func GetImgType() []string {
	return imgType
}

func (client *ApiClient) Hello(c *gin.Context) {
	basic.RespWithMsg(c, "hello")
}

func (client *ApiClient) FileUpload(c *gin.Context) {
	file, header, err := c.Request.FormFile("upload")
	if err != nil {
		basic.RespWithErr(c, 400, errors.Wrap(err, "upload error"))
		return
	}
	filename := header.Filename
	fmt.Println(header.Filename)
	out, err := os.Create("/data/service/sky/data/pic/" + filename)
	if err != nil {
		basic.RespWithErr(c, 400, errors.Wrap(err, "upload error"))
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		basic.RespWithErr(c, 400, errors.Wrap(err, "Copy error"))
		return
	}
	basic.RespWithMsg(c, "upload success")
}

func (client *ApiClient) PhotoUpload(c *gin.Context) {
	upfile, header, err := c.Request.FormFile("upload")
	if err != nil {
		basic.RespWithErr(c, 400, errors.Wrap(err, "upload error"))
		return
	}
	defer upfile.Close()
	// 图片解码 --------------------------------------
	// 读入缓存
	bufUpFile := bufio.NewReader(upfile)
	img, imgtype, err := image.Decode(bufUpFile)
	if err != nil {
		basic.RespWithErr(c, 400, errors.Wrap(err, "decode error"))
		return
	}
	if !IsType(imgtype) {
		basic.RespWithErr(c, 400, errors.Wrap(err, "not support this type"))
		return
	}
	// 设置文件读写下标 --------------------------------

	// 设置下次读写位置（移动文件指针位置）
	_, err = upfile.Seek(0, 0)
	if err != nil {
		basic.RespWithErr(c, 400, errors.Wrap(err, "设置文件读写位置失败"))
		return
	}
	// 计算文件的 MD5 值 -----------------------------

	// 初始化 MD5 实例
	md5Hash := md5.New()
	// 读入缓存
	bufFile := bufio.NewReader(upfile)
	_, err = io.Copy(md5Hash, bufFile)
	if err != nil {
		basic.RespWithErr(c, 400, errors.Wrap(err, "计算文件MD5失败"))
		return
	}
	// 进行 MD5 算计，返回 16进制的 byte 数组
	fileMd5FX := md5Hash.Sum(nil)
	fileMd5 := fmt.Sprintf("%x", fileMd5FX)

	// 目录计算 --------------------------------------

	// 组合文件完整路径
	dirPath := util.JoinPath(fileMd5) + "/" // 目录
	filePath := dirPath + fileMd5           // 文件路径
	// 获取目录信息，并创建目录
	dirInfo, err := os.Stat(dirPath)
	if err != nil {
		err = os.MkdirAll(dirPath, 0666)
		if err != nil {
			basic.RespWithErr(c, 400, errors.Wrap(err, "创建目录失败"))
			return
		}
	} else {
		if !dirInfo.IsDir() {
			err = os.MkdirAll(dirPath, 0666)
			if err != nil {
				basic.RespWithErr(c, 400, errors.Wrap(err, "创建目录失败"))
				return
			}
		}
	}
	// 存入文件 --------------------------------------

	_, err = os.Stat(filePath)
	if err != nil {
		// 打开一个文件,文件不存在就会创建
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			basic.RespWithErr(c, 400, errors.Wrap(err, "文件创建失败"))
			return
		}
		defer file.Close()

		if imgtype == PNG {
			err = png.Encode(file, img)
		} else if imgtype == JPG || imgtype == JPEG {
			err = jpeg.Encode(file, img, nil)
		}
		if err != nil {
			basic.RespWithErr(c, 400, errors.Wrap(err, "图片生成失败"))
			return
		}

	}

	respData := UpdateDate{
		Size:  header.Size,
		Mime:  imgtype,
		Imgid: fileMd5,
	}
	basic.RespWithData(c, respData)
}

// 判断是否有这个图片类型
func IsType(str string) bool {

	// 转小写
	str = strings.ToLower(str)

	for _, v := range imgType {
		if v == str {
			return true
		}
	}

	return false
}

func (client *ApiClient) GetImg(c *gin.Context) {
	urlParse, ok := c.GetQuery("img")
	if !ok {
		basic.RespWithError(c, 400, "img id empty")
	}
	filePath := util.UrlParse(urlParse)
	if filePath == "" {
		basic.RespWithMsg(c, "img not found")
		return
	}
	heightStr := c.DefaultQuery("h", "0")
	widthStrr := c.DefaultQuery("w", "0")
	height, _ := strconv.Atoi(heightStr)
	width, _ := strconv.Atoi(widthStrr)
	CutImage(c, filePath, width, height)

}

// 裁剪图像
func CutImage(w *gin.Context, path string, width, height int) {

	// 没有宽高，就是在加载原图像
	if width == 0 && height == 0 {

		file, err := os.Open(path)
		if err != nil {
			NoImage(w)
			return
		}
		defer file.Close()

		_, err = io.Copy(w.Writer, file)
		if err != nil {
			basic.RespWithErr(w, 400, errors.Wrap(err, "jpeg Encode error"))
			return
		}
		return
	}

	// 裁剪图像 --------------------------------------

	// 裁剪图像的组合路径
	CutPath := fmt.Sprintf("%s_%d_%d", path, width, height)

	// 判断是否存在裁剪图像
	_, err := os.Stat(CutPath)
	if err == nil {

		file, err := os.Open(CutPath)
		if err != nil {
			NoImage(w)
			return
		}
		defer file.Close()

		_, err = io.Copy(w.Writer, file)
		if err != nil {
			basic.RespWithErr(w, 400, errors.Wrap(err, "jpeg Encode error"))
			return
		}
		return
	}

	// 判断是否存在原图像
	_, err = os.Stat(path)
	if err != nil {
		NoImage(w)
		return
	}

	// 原图像
	file, err := os.Open(path)
	if err != nil {
		NoImage(w)
		return
	}
	defer file.Close()

	// 图片解码 --------------------------------------

	bufFile := bufio.NewReader(file)
	img, imgtype, err := image.Decode(bufFile)
	if err != nil {
		NoImage(w)
		return
	}

	// gif 图就不处理了
	if imgtype == GIF {

		_, err = file.Seek(0, 0)
		if err != nil {
			NoImage(w)
			return
		}

		_, err = io.Copy(w.Writer, file)
		if err != nil {
			basic.RespWithErr(w, 400, errors.Wrap(err, "jpeg Encode error"))
			return
		}
		return
	}

	if width > img.Bounds().Max.X {
		width = img.Bounds().Max.X
	}

	if height > img.Bounds().Max.Y {
		height = img.Bounds().Max.Y
	}

	// 进行裁剪
	reimg := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
	// 裁剪的存储
	out, err := os.Create(CutPath)
	if err != nil {
		NoImage(w)
		return
	}
	defer out.Close()

	if imgtype == JPEG || imgtype == JPG {
		err = jpeg.Encode(out, reimg, nil)
		if err != nil {
			basic.RespWithErr(w, 400, errors.Wrap(err, "jpeg Encode error"))
			return
		}
	} else if imgtype == PNG {
		err = png.Encode(out, reimg)
		if err != nil {
			basic.RespWithErr(w, 400, errors.Wrap(err, "png Encode error"))
			return
		}
	}

	_, err = out.Seek(0, 0)
	if err != nil {
		NoImage(w)
		return
	}

	_, err = io.Copy(w.Writer, out)
	if err != nil {
		basic.RespWithErr(w, 400, errors.Wrap(err, "Copy error"))
		return
	}

}

func NoImage(c *gin.Context) {

	// 图片流方式输出
	c.Header("Content-Type", "image/png")
	// 进行图片的编码
	err := png.Encode(c.Writer, util.NoImg)
	if err != nil {
		basic.RespWithErr(c, 400, err)
	}
}
