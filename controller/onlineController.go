package controller

import (
	"GFV/download"
	"GFV/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type OnlineController struct {
}

type NowFile struct {
	Md5            string
	Ext            string
	LastActiveTime int64
}

var (
	Pattern      string
	Address      string
	AllFile      map[string]*NowFile
	ExpireTime   int64
	AllOfficeEtx = []string{".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".txt"}
	AllImageEtx  = []string{".jpg", ".png", ".gif"}
)

type preview struct {
	Url  string `json:"url"`
	Type string `json:"type"`
}

func (oc OnlineController) OnlinePreview(c *gin.Context) {
	requestUrl := c.Query("url")
	requestType := c.Query("type")
	filePath := requestUrl[len(Pattern):]
	if filePath, err := download.DownloadFile(filePath, "cache/download/"+path.Base(filePath)); err == nil {
		if requestType == "pdf" && (path.Ext(filePath) == ".pdf" || utils.IsInArr(path.Ext(filePath), AllOfficeEtx)) { //预留的PDF预览接口
			if path.Ext(filePath) == ".pdf" {
				dataByte := pdfPageDownload("cache/download/" + path.Base(filePath))
				c.Writer.Header().Set("content-length", strconv.Itoa(len(dataByte)))
				c.Writer.Header().Set("content-type", "text/html;charset=UTF-8")
				c.Writer.Write([]byte(dataByte))
				setFileMap(path.Base(filePath))
			} else {
				//若传入文件不是pdf，则进行转换
				//若文件已存在，则获取缓存的文件
				if isHavePdf(path.Base(filePath)) {
					pdfPath := "cache/pdf/" + strings.Split(path.Base(filePath), ".")[0] + ".pdf"
					dataByte := pdfPage(pdfPath)
					c.Writer.Header().Set("content-length", strconv.Itoa(len(dataByte)))
					c.Writer.Header().Set("content-type", "text/html;charset=UTF-8")
					c.Writer.Write([]byte(dataByte))
					setFileMap(path.Base(filePath))
					return
				}
				if pdfPath := utils.ConvertToPDF(filePath); pdfPath != "" {
					dataByte := pdfPage("cache/pdf/" + path.Base(pdfPath))
					c.Writer.Header().Set("content-length", strconv.Itoa(len(dataByte)))
					c.Writer.Header().Set("content-type", "text/html;charset=UTF-8")
					c.Writer.Write([]byte(dataByte))
					setFileMap(path.Base(filePath))
				} else {
					c.Writer.Write([]byte("转换为PDF时出现错误!"))
				}
			}
		} else if utils.IsInArr(path.Ext(filePath), AllImageEtx) {
			dataByte := imagePage(filePath)
			c.Writer.Header().Set("content-length", strconv.Itoa(len(dataByte)))
			c.Writer.Header().Set("content-type", "text/html;charset=UTF-8")
			c.Writer.Write([]byte(dataByte))
		} else if utils.IsInArr(path.Ext(filePath), AllOfficeEtx) {
			if isHave(path.Base(filePath)) {
				dataByte := officePage("cache/convert/" + strings.Split(path.Base(filePath), ".")[0])
				c.Writer.Header().Set("content-length", strconv.Itoa(len(dataByte)))
				c.Writer.Header().Set("content-type", "text/html;charset=UTF-8")
				c.Writer.Write([]byte(dataByte))
				return
			}
			if pdfPath := utils.ConvertToPDF(filePath); pdfPath != "" {
				if imgPath := utils.ConvertToImg(pdfPath); imgPath != "" {
					dataByte := officePage(imgPath)
					c.Writer.Header().Set("content-length", strconv.Itoa(len(dataByte)))
					c.Writer.Header().Set("content-type", "text/html;charset=UTF-8")
					c.Writer.Write([]byte(dataByte))
					setFileMap(path.Base(filePath))
				} else {
					c.Writer.Write([]byte("转换为图片时出现错误!"))
				}
			} else {
				c.Writer.Write([]byte("转换为PDF时出现错误!"))
			}
		} else {
			log.Println("Error: <", err, "> when download file")
			c.Writer.Write([]byte("获取文件失败...请检查你的路径是否正确!"))
		}
	} else if path.Ext(filePath) == ".pdf" {
		if isHave(path.Base(filePath)) {
			dataByte := officePage("cache/convert/" + strings.Split(path.Base(filePath), ".")[0])
			c.Writer.Header().Set("content-length", strconv.Itoa(len(dataByte)))
			c.Writer.Header().Set("content-type", "text/html;charset=UTF-8")
			c.Writer.Write([]byte(dataByte))
			return
		}
		if imgPath := utils.ConvertToImg(filePath); imgPath != "" {
			dataByte := officePage(imgPath)
			c.Writer.Header().Set("content-length", strconv.Itoa(len(dataByte)))
			c.Writer.Header().Set("content-type", "text/html;charset=UTF-8")
			c.Writer.Write([]byte(dataByte))
			setFileMap(path.Base(filePath))
		} else {
			c.Writer.Write([]byte("转换为图片时出现错误!"))
		}
	} else {
		c.Writer.Write([]byte("获取文件失败...请检查你的路径是否正确!"))
	}
}

func (oc OnlineController) OfflinePreview(c *gin.Context) {
	filePath := c.Param("filePath")
	filePath = c.Query("filePath")
	fmt.Println(filePath)
	if isHaveImg(filePath) {
		dataByte := officePage("cache/convert/" + strings.Split(path.Base(filePath), ".")[0])
		c.Writer.Header().Set("content-length", strconv.Itoa(len(dataByte)))
		c.Writer.Header().Set("content-type", "text/html;charset=UTF-8")
		c.Writer.Write([]byte(dataByte))
		return
	}
}

func (oc OnlineController) Static(c *gin.Context) {
	url := c.Request.URL.String()
	DataByte, err := ioutil.ReadFile("html" + url)
	if err != nil {
		c.Writer.Header().Set("content-length", strconv.Itoa(len("404")))
		c.Writer.Header().Set("content-type", "text/html;charset=UTF-8")
		c.Writer.Write([]byte("出现了一些问题,导致File View无法获取您的数据!"))
		return
	}
	c.Writer.Header().Set("content-length", strconv.Itoa(len(DataByte)))
	if path.Ext(url) == ".css" {
		c.Writer.Header().Set("content-type", "text/css;charset=UTF-8")
	} else if path.Ext(url) == ".js" {
		c.Writer.Header().Set("content-type", "application/x-javascript;charset=UTF-8")
	}
	c.Writer.Write(DataByte)
}

func pdfPageDownload(filePath string) []byte {
	dataByte, _ := ioutil.ReadFile("html/pdf.html")
	dataStr := string(dataByte)
	pdfUrl := "img_asset/" + path.Base(filePath)
	dataStr = strings.Replace(dataStr, "{{url}}", pdfUrl, -1)
	dataByte = []byte(dataStr)
	return dataByte
}
func isHavePdf(fileName string) bool {
	fileName = strings.Split(fileName, ".")[0] + ".pdf"
	// 指定目录和文件名
	directory := "cache/pdf/"
	// 使用 filepath 包来生成文件的完整路径
	filePath := filepath.Join(directory, fileName)
	// 使用 os.Stat 来获取文件信息
	_, err := os.Stat(filePath)

	if err == nil {
		return true
	} else {
		if os.IsNotExist(err) {
			return false
		} else {
			log.Println("ERROR: 查看文件缓存 ", err)
			return false
		}
	}
}
func isHave(fileName string) bool {
	fileName = strings.Split(fileName, ".")[0]
	if _, ok := AllFile[fileName]; ok {
		AllFile[fileName].LastActiveTime = time.Now().Unix()
		return true
	} else {
		return false
	}
}

func officePage(imgPath string) []byte {
	rd, _ := ioutil.ReadDir(imgPath)
	sort.Slice(rd, func(i, j int) bool {
		numI, _ := strconv.Atoi(strings.TrimSuffix(rd[i].Name(), filepath.Ext(rd[i].Name())))
		numJ, _ := strconv.Atoi(strings.TrimSuffix(rd[j].Name(), filepath.Ext(rd[j].Name())))
		return numI < numJ
	})
	dataByte, _ := ioutil.ReadFile("html/office.html")
	dataStr := string(dataByte)
	htmlCode := ""
	for _, fi := range rd {
		if !fi.IsDir() {
			htmlCode = htmlCode + `<img class="my-photo" alt="loading" title="查看大图" style="cursor: pointer;"
									data-src="office_asset/` + path.Base(imgPath) + "/" + fi.Name() + `" src="images/loading.gif"
									">`
		}
	}
	dataStr = strings.Replace(dataStr, "{{AllImages}}", htmlCode, -1)
	dataByte = []byte(dataStr)
	return dataByte
}

func imagePage(filePath string) []byte {
	dataByte, _ := ioutil.ReadFile("html/image.html")
	dataStr := string(dataByte)
	imageUrl := "img_asset/" + path.Base(filePath)
	htmlCode := `<li>
					<img id="` + imageUrl + `" url="` + imageUrl + `"
						src="` + imageUrl + `" width="1px" height="1px">
				 </li>`
	dataStr = strings.Replace(dataStr, "{{AllImages}}", htmlCode, -1)
	dataStr = strings.Replace(dataStr, "{{FirstPath}}", imageUrl, -1)
	dataByte = []byte(dataStr)
	return dataByte
}

func pdfPage(filePath string) []byte {
	dataByte, _ := ioutil.ReadFile("html/pdf.html")
	dataStr := string(dataByte)
	pdfUrl := "pdf_asset/" + path.Base(filePath)
	dataStr = strings.Replace(dataStr, "{{url}}", pdfUrl, -1)
	dataByte = []byte(dataStr)
	return dataByte
}
func setFileMap(fileName string) {
	ext := path.Ext(fileName)
	fileName = strings.Split(fileName, ".")[0]
	if _, ok := AllFile[fileName]; ok {
		AllFile[fileName].LastActiveTime = time.Now().Unix()
		return
	} else {
		temp := &NowFile{
			Md5:            fileName,
			Ext:            ext,
			LastActiveTime: time.Now().Unix(),
		}
		AllFile[fileName] = temp
	}
}
func isHaveImg(fileName string) bool {
	// 指定目录和文件名
	directory := "cache/convert/"
	// 使用 filepath 包来生成文件的完整路径
	filePath := filepath.Join(directory, fileName)
	// 使用 os.Stat 来获取文件信息
	_, err := os.Stat(filePath)

	if err == nil {
		return true
	} else {
		if os.IsNotExist(err) {
			return false
		} else {
			log.Println("ERROR: 查看文件缓存 ", err)
			return false
		}
	}
}
