package handlers

import (
	"fmt"
)

func CollectionNameFAQCategory(appid string) string {
	return fmt.Sprintf("faq:%s:category", appid)
}

func CollectionNameFAQ(appid string) string {
	return fmt.Sprintf("faq:%s", appid)
}
