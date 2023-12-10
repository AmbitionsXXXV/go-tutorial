//go:build ignore

package main

import (
	"html/template" // 引入HTML模板库
	"log"           // 引入日志库，用于记录日志
	"net/http"      // 引入HTTP库，用于处理HTTP请求
	"os"            // 引入操作系统库，用于文件操作
	"regexp"        // 引入正则表达式库
)

// Page 定义了一个页面的结构，包括标题和正文
type Page struct {
	Title string
	Body  []byte
}

// save 方法用于将页面保存到文件系统
func (p *Page) save() error {
	filename := p.Title + ".txt"

	return os.WriteFile(filename, p.Body, 0600) // 保存文件，文件权限设置为0600
}

// loadPage 从文件系统加载页面
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename) // 读取文件

	if err != nil {
		return nil, err // 读取出错，返回错误
	}

	return &Page{Title: title, Body: body}, nil // 返回加载的页面
}

// viewHandler 处理/view/路径的HTTP请求
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)

	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound) // 页面不存在，重定向到编辑页面

		return
	}

	renderTemplate(w, "view", p) // 渲染页面
}

// editHandler 处理/edit/路径的HTTP请求
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)

	if err != nil {
		p = &Page{Title: title} // 页面不存在，创建新页面
	}

	renderTemplate(w, "edit", p) // 渲染编辑页面
}

// saveHandler 处理/save/路径的HTTP请求
func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // 保存出错，返回错误

		return
	}

	http.Redirect(w, r, "/view/"+title, http.StatusFound) // 保存成功，重定向到查看页面
}

// templates 缓存了HTML模板
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

// renderTemplate 使用HTML模板渲染页面
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // 渲染出错，返回错误
	}
}

// validPath 是一个正则表达式，用于匹配有效的路径
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

// makeHandler 包装了HTTP处理函数，用于提取路径中的标题
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)

		if m == nil {
			http.NotFound(w, r) // 路径无效，返回404

			return
		}

		fn(w, r, m[2]) // 调用处理函数
	}
}

// main 函数设置了路由并启动HTTP服务器
func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	log.Fatal(http.ListenAndServe(":8080", nil)) // 在8080端口启动服务器
}
