package iternal

import (
	"fmt"
	"strings"
)

func normalizeURL(url string) (string,error) {
    if strings.HasPrefix(url, "https://") {
        url = strings.TrimPrefix(url, "https://")
    } else if strings.HasPrefix(url, "http://") {
        url = strings.TrimPrefix(url, "http://")
    }else{
		return "",fmt.Errorf("error occured durning normalizeURL")
	}
    url = strings.TrimSuffix(url, "/")
    return url,nil
}
