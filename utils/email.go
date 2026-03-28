package utils

import (
	"bytes"
	"fmt"
	"gblog/config"
	"html/template"
	"path/filepath"

	"gopkg.in/gomail.v2"
)

// Mailer 结构体，方便以后扩展（比如切换到 Resend）
type Mailer struct {
	Host string
	Port int
	User string
	Pass string
}

// NewMailer 初始化发信器
func NewMailer() *Mailer {
	return &Mailer{
		Host: config.AppConfig.Mail.Host,
		Port: config.AppConfig.Mail.Port,
		User: config.AppConfig.Mail.User,
		Pass: config.AppConfig.Mail.Password,
	}
}

// ArticleEmailData 文章邮件数据结构
type ArticleEmailData struct {
	Title    string
	URL      string
	Username string
}

// RenderArticleEmailTemplate 渲染文章邮件模板
func RenderArticleEmailTemplate(title, articleID, username string) (string, error) {
	data := ArticleEmailData{
		Title:    title,
		URL:      "https://orange2006.online/articles/" + articleID,
		Username: username,
	}

	templatePath := filepath.Join("templates", "article_email.html")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to parse email template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute email template: %w", err)
	}

	return buf.String(), nil
}

// Send 发送邮件的具体实现
func (m *Mailer) Send(to string, subject string, body string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.User)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	d := gomail.NewDialer(m.Host, m.Port, m.User, m.Pass)

	if err := d.DialAndSend(msg); err != nil {
		return fmt.Errorf("failed to send email to %s: %w", to, err)
	}
	return nil
}

// SendArticleEmail 发送文章通知邮件
func SendArticleEmail(to, title, articleID, username string) error {
	body, err := RenderArticleEmailTemplate(title, articleID, username)
	if err != nil {
		return fmt.Errorf("failed to render email template: %w", err)
	}

	subject := "🍊 橘子发布了新文章：" + title
	return NewMailer().Send(to, subject, body)
}
