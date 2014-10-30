package handlers

import (
	"fmt"
)

func CollectionNameQnA(appid string) string {
	return fmt.Sprintf("qna:%s", appid)
}

func CollectionNameFAQCategory(appid string) string {
	return fmt.Sprintf("faq:%s:category", appid)
}

func CollectionNameFAQ(appid string) string {
	return fmt.Sprintf("faq:%s", appid)
}

func CollectionNameNotice(appid string) string {
	return fmt.Sprintf("notice:%s", appid)
}
