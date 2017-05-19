package process

import (
	"fmt"

	"os"
	"path"
	"strconv"
	"strings"

	pinyin "github.com/mozillazg/go-pinyin"
	"github.com/nipeone/imgprocessor/file"
)

var args pinyin.Args
var limit int

const detail = "./detail.txt"
const summary = "./summary.txt"
const vocPath = "./voc.txt"

// Processor 预处理
type Processor interface {
	ParseLabel(labelPath string) map[string][]string
	Process(d map[string][]string, v ...interface{})
}

// ImgProcessor 图片预处理
type ImgProcessor struct {
	Fr string
	To string
}

func init() {
	args = pinyin.NewArgs()
	args.Separator = ""
	args.Fallback = func(r rune, a pinyin.Args) []string {
		return []string{string(r)}
	}
}

// New 默认构建函数
func New(from, to string, l int) *ImgProcessor {
	limit = l
	return &ImgProcessor{
		Fr: from,
		To: to,
	}
}

// ParseLabel 解析标签文件,格式如下
// col1	col2
// :	:
// img	label1,label2,label3
func (p *ImgProcessor) ParseLabel(labelPath string) (map[string][]string, map[string]string) {
	lines := file.ReadLines(labelPath)
	d := make(map[string][]string)
	voc := make(map[string]string)
	for _, line := range lines {
		parts := strings.Split(line, "\t")
		if len(parts) != 2 {
			panic(fmt.Errorf("%s can not be split into 2 parts", line))
		} else {
			img := string(parts[0][strings.LastIndex(parts[0], "/")+1:])
			labels := strings.Split(strings.TrimSpace(parts[1]), ",")
			for _, l := range labels {
				//将中文转为拼音
				k := pinyin.Slug(strings.TrimSpace(l), args)
				voc[k] = strings.TrimSpace(l)
				if v, ok := d[k]; ok {
					d[k] = append(v, img)
				} else {
					d[k] = []string{img}
				}
			}
		}
	}
	for k, v := range d {
		if !filter(v) {
			delete(voc, k)
			delete(d, k)
		}
	}
	return d, voc
}

// Process 根据标签解析的键值对，将图片移动到对应分类
func (p *ImgProcessor) Process(d map[string][]string) {
	if file.IsExist(p.To) {
		os.RemoveAll(p.To)
	}
	os.Mkdir(p.To, os.ModeDir)
	count := 0
	for k, v := range d {
		fmt.Println(k, " ", len(v))
		count++
		if !file.IsExist(path.Join(p.To, k)) {
			os.Mkdir(path.Join(p.To, k), os.ModeDir)
		}
		for _, f := range v {
			file.CopyFile(path.Join(p.Fr, f), path.Join(p.To, k, f))
		}
	}
	fmt.Println("the filtered photo limit is ", limit)
	fmt.Println("total classes:", len(d), "/", "filtered classes:", count)
}

//Voc 生成中英文对照表
func (p *ImgProcessor) Voc(voc map[string]string) {
	path := path.Join(p.To, vocPath)
	lines := []string{}
	for k, v := range voc {
		lines = append(lines, k+":"+v)
	}
	file.WriteLines(path, lines)
}

func filter(v []string) bool {
	return len(v) > limit
}

func Detail(d map[string][]string) {
	if file.IsExist(detail) {
		os.Remove(detail)
	}
	f, _ := os.Create(detail)
	// f, _ := os.OpenFile("total.txt", os.O_APPEND, 0666)
	// 遍历map
	for k, v := range d {
		f.WriteString(k + " " + strconv.Itoa(len(v)) + "\n")
		f.WriteString(strings.Join(v, "\n") + "\n")
		f.WriteString("\n")

	}

	defer f.Close()
}

func Summary(d map[string][]string) {
	if file.IsExist(summary) {
		os.Remove(summary)
	}
	f, _ := os.Create(summary)

	// 遍历map
	for k, v := range d {
		f.WriteString(k + "\t" + strconv.Itoa(len(v)) + "\n")
	}

	defer f.Close()
}

func Test() {
	py := pinyin.Slug("2015中国人abc", args)
	fmt.Println(py)

}
