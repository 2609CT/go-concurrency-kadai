package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type SafeUsrMap struct {
	visited map[string]bool
	mu      sync.Mutex
}

func (s *SafeUsrMap) Visit(url string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.visited[url] {
		s.visited[url] = true
		return true
	}
	return false
}

// URLを受け取り、そのページに含まれるリンク（文字列のスライス）を返します
func fetch(url string) []string {
	// 処理時間をランダムにシミュレーション（0〜1.3秒）
	// ※タイムアウト(1秒)の挙動を確認するため
	time.Sleep(time.Duration(rand.Intn(1300)) * time.Millisecond)

	// 擬似的なリンクデータ
	links := map[string][]string{
		"https://example.com":   {"https://example.com/a", "https://example.com/b"},
		"https://example.com/a": {"https://example.com/c"},
		"https://example.com/b": {"https://example.com/a"}, // 重複リンク
		"https://example.com/c": {},
	}

	return links[url]
}

func Crawl(url string, wg *sync.WaitGroup, smf *SafeUsrMap) {

	defer wg.Done()

	if smf.Visit(url) {

		fmt.Printf("取得開始：%s\n", url)

		ch := make(chan []string, 1)
		//chにURL(スライス)を受信
		go func() {
			ch <- fetch(url)
		}()

		select {
		case val := <-ch:
			fmt.Printf("取得成功：%s -> リンク数：%d\n", url, len(val))
			if len(val) != 0 {
				for _, v := range val {
					wg.Add(1)
					go Crawl(v, wg, smf)
				}
			}

		case <-time.After(1 * time.Second):
			fmt.Printf("タイムアウト：%s\n", url)
		}
	} else {
		fmt.Printf("重複スキップ：%s\n", url)
	}
}

func main() {

	var wg sync.WaitGroup

	sfm := SafeUsrMap{
		visited: map[string]bool{
			"https://example.com":   false,
			"https://example.com/a": false,
			"https://example.com/b": false,
			"https://example.com/c": false,
		},
	}

	startURL := "https://example.com"

	wg.Add(1)

	go Crawl(startURL, &wg, &sfm)

	wg.Wait()

	fmt.Println("すべてのクロールが完了しました")
}
