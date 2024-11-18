package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"unicode/utf8"
)

const (
	MAX_RUNES = 20
	MAX_INT   = 10000
)

func terminationListener(c <-chan os.Signal) {
	sig := <-c
	fmt.Printf("Signal received: %v\n", sig)
	os.Exit(1)
}

func inputReader(c chan<- string) {

	intro := `任意の値を入力してエンターを押すと、文字種別(string or int) に応じて異なる統計を出力します。
値を入力してください。

`
	fmt.Println(intro)
	buf := bufio.NewReader(os.Stdin)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		line = strings.TrimSpace(line)

		if len(line) != 0 {
			c <- line
		}
	}
}

func inputDispatcher(lineInputChan <-chan string, stringChan chan<- string, stringDone <-chan struct{}, intChan chan<- int, intDone <-chan struct{}) {
	for line := range lineInputChan {
		intVal, err := strconv.Atoi(line)
		if err != nil {
			select {
			case <-stringDone:
				slog.Info("StringCounter is already stopped.")
			default:
				stringChan <- line
			}
		} else {
			select {
			case <-intDone:
				slog.Info("IntCounter is already stopped.")
			default:
				intChan <- intVal
			}
		}
	}
}

func stringCounter(str <-chan string, done chan<- struct{}) {
	var bytes, totalBytes, runes, totalRunes int
	for s := range str {
		bytes = len(s)
		runes = utf8.RuneCountInString(s)
		totalBytes += bytes
		totalRunes += runes
		if totalRunes > MAX_RUNES {
			close(done)
			break
		}
		slog.Info(fmt.Sprintf("Bytes: %v, Runes: %v, Value: %v", bytes, runes, s), slog.Int("totalBytes", totalBytes), slog.Int("totalRunes", totalRunes))
	}
	slog.Info("StringCounter is stopped.")
}

func intCounter(num <-chan int, done chan<- struct{}) {
	var total int
	for n := range num {
		total += n
		if total > MAX_INT {
			close(done)
			break
		}
		slog.Info(fmt.Sprintf("Value: %v", n), slog.Int("total", total))
	}
	slog.Info("IntCounter is stopped.")
}

func main() {
	/*
		TODO:
			◯各カウンタ関数の引数を ctx と value で分ける
			◯関数を命名し main から外だし
			◯slog で構造化ログ出力、専用の goroutine で実装
				標準出力
				ファイル
			◯コマンドライン引数を明示的に指定した際のみ、ログをファイル出力できるようにする
				指定したディレクトリへの読み書きができない場合はエラーを出して終了する
			各カウンタの閾値を超えたら goroutine を終了させる
				共通：回数制限でプログラム全体を終了させる
				◯個別：処理が閉じている旨のエラーログを出力させる
					string: 20 文字
					int: 10000
				各カウンタが終了したら終了する
			テストしやすいコードに修正
			入力受付状態がわかりやすいようプロンプトを表示
			URL に対してリクエストした結果の統計を取るカウンタを実装

	*/

	// コマンドライン引数設定
	var logPath string
	flag.StringVar(&logPath, "p", "", "Log output `path`.")
	flag.Parse()

	// ログ出力
	var writer *os.File
	if len(logPath) != 0 {
		var err error
		writer, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
		if err != nil {
			log.Fatal(err)
		}
		defer writer.Close()
	} else {
		writer = os.Stdout
	}
	handler := slog.NewJSONHandler(writer, nil)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// シグナル待ち受け
	sig := make(chan os.Signal, 1)
	go terminationListener(sig)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// 文字列集計
	stringChan := make(chan string)
	stringDone := make(chan struct{})
	go stringCounter(stringChan, stringDone)

	// 数値集計
	intChan := make(chan int)
	intDone := make(chan struct{})
	go intCounter(intChan, intDone)

	// 入力行の判別と振り分け
	lineInputChan := make(chan string)
	go inputDispatcher(lineInputChan, stringChan, stringDone, intChan, intDone)

	// メイン処理
	inputReader(lineInputChan)
}
