package i18n

import (
	"encoding/json"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/share-group/share-go/provider/config"
	loggerFactory "github.com/share-group/share-go/provider/logger"
	"golang.org/x/text/language"
	"os"
	"path"
)

// https://github.com/nicksnyder/go-i18n/blob/main/.github/README.zh-Hans.md

var logger = loggerFactory.GetLogger()
var bundle = i18n.NewBundle(language.English)

func init() {
	i18nBasePath := path.Join(config.GetRootDir(), "i18n")
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	files, err := os.ReadDir(i18nBasePath)
	if err != nil {
		// 开发者没有定义多语言文件
		return
	}

	for _, file := range files {
		bundle.MustLoadMessageFile(path.Join(i18nBasePath, file.Name()))
		logger.Info(path.Join(i18nBasePath, file.Name()))
	}
}

// 执行翻译
//
// i18nKey-多语言key;locale-期望输出的语言;templateData-如果是字符串模板的，需要传入模板参数
func T(i18nKey string, locale language.Tag, templateData ...map[string]any) string {
	localize := i18n.NewLocalizer(bundle, locale.String())
	localizeConfig := &i18n.LocalizeConfig{MessageID: i18nKey}
	if len(templateData) > 0 {
		localizeConfig.TemplateData = templateData[0]
	}
	return localize.MustLocalize(localizeConfig)
}
