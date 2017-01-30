package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
	"syscall"
)

// echo "funayose" | go run main.go FUNAYOSE
func main() {
	tr(os.Stdin, os.Stdout, os.Stderr)
}

func tr(src io.Reader, dst io.Writer, errDst io.Writer) error {
	// 実行すコマンド tr a-z A-Z
	cmd := exec.Command("tr", "a-z", "A-Z")

	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	err := cmd.Start() // コマンド実行
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		// コマンドの標準入力にsrcからコピーする
		_, err := io.Copy(stdin, src)
		if e, ok := err.(*os.PathError); ok && e.Err == syscall.EPIPE {
		} else if err != nil {
			log.Println("failed to wite to STDIN", err)
		}
		stdin.Close()
		wg.Done()
	}()
	go func() {
		// コマンドの標準出力をdstにコピーする
		io.Copy(dst, stdout)
		stdout.Close()
		wg.Done()
	}()
	go func() {
		// コマンドの標準エラー出力をerrDstにコピーする
		io.Copy(errDst, stderr)
		stderr.Close()
		wg.Done()
	}()

	wg.Wait()
	return cmd.Wait()
}
