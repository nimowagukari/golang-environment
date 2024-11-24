package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
)

func parseArgs() []string {
	flag.Parse()
	return flag.Args()
}

func readBytesFromFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func splitByNewLine(b []byte) [][]byte {
	return bytes.Split(b, []byte("\n"))
}

func removeComment(byteSlice [][]byte) ([][]byte, error) {
	re, err := regexp.Compile(`^\s*#`)
	if err != nil {
		return nil, err
	}

	for i, l := range byteSlice {
		if re.Match(l) {
			byteSlice[i] = []byte("")
		}
	}

	return byteSlice, nil
}

func writeOutput(byteSlice [][]byte, writer io.Writer) error {
	b := bytes.Join(byteSlice, []byte("\n"))
	_, err := writer.Write(b)
	if err != nil {
		return err
	}
	return nil
}

// main program
func main() {
	/*
		ロジック整理
			引数のパース
			ファイルの中身を []byte で詠み込む
			改行コードごとに配列で分割
			上から処理
				ラインコメントがあったらその行ごと削除
				ブロックコメントの開始があったら、閉じワードがあるまで消す

		機能拡張
			複数ファイルを受け入れる
			標準入力を受け入れる
			対応するコメントアウト種別を増やす
			拡張子からコメントアウト種別を判別
		周辺環境整備
			go doc 実装
			README, Example 追加

	*/

	filePaths := parseArgs()

	for _, f := range filePaths {
		b, err := readBytesFromFile(f)
		if err != nil {
			log.Fatal(err)
		}

		bs := splitByNewLine(b)
		bs, err = removeComment(bs)
		if err != nil {
			log.Fatal(err)
		}
		if len(filePaths) > 1 {
			fmt.Printf("\n%v: \n", f)
		}
		if err := writeOutput(bs, os.Stdout); err != nil {
			log.Fatal(err)
		}

	}
}
