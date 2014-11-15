package handlers

import (
	"../../src"
	"fmt"
)

func CollectionNameQnA(appid string) string {
	return fmt.Sprintf("%s:qna:%s", config.NAMESPACE, appid)
}

func CollectionNameFAQCategory(appid string) string {
	return fmt.Sprintf("%s:faq:%s:category", config.NAMESPACE, appid)
}

func CollectionNameFAQ(appid string) string {
	return fmt.Sprintf("%s:faq:%s", config.NAMESPACE, appid)
}

func CollectionNameNotice(appid string) string {
	return fmt.Sprintf("%s:notice:%s", config.NAMESPACE, appid)
}

func CollectionTable(classesName, appid string) string {
	return fmt.Sprintf("%s:%s:%s", config.NAMESPACE, classesName, appid)
}

func HttpErr(code int, err string) (int, map[string]interface{}) {

	//var str string
	//fmt.Fprintf(c, "%s", `{"page": 1, "fruits": ["apple", "peach"]}`)

	return code, map[string]interface{}{"code": code, "err": err}

}
