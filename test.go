package main

// import (
// 	"fmt"
// 	"strings"
// )

// func extractRepoName(url string) (string, error) {
// 	const prefix = "repos/"
// 	startIndex := strings.Index(url, prefix)
// 	if startIndex == -1 {
// 		return "", errors.New("prefix not found in URL")
// 	}
// 	parts := strings.SplitN(url[startIndex+len(prefix):], "/", 3)
// 	if len(parts) < 2 {
// 		return "", errors.New("invalid URL format")
// 	}
// 	return parts[1], nil
// }

// // func main() {
// // 	url := "https://api.github.com/repos/octocat/Hello-World/git/commits/6dcb09b5b57875f334f61aebed695e2e4193db5e"
// // 	repoName, err := extractRepoName(url)
// // 	if err != nil {
// // 		fmt.Println("Error:", err)
// // 	} else {
// // 		fmt.Println("Repository Name:", repoName)
// // 	}
// // }
