package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/term"
)

const (
	KeyEnter     byte = 13
	KeyEscape    byte = 27
	KeySpace     byte = 32
	KeyBackspace byte = 127

	CRLF             string = "\r\n"
	ClearScreen      string = "\033[2J"
	RequestCursorPos string = "\033[6n"
	SaveCursorPos           = "\033[s"
	RestoreCursorPos        = "\033[u"

	MoveCursorLeft string = "\033[D"
	MoveCursorUp   string = "\033[A"
)

type ExitBuf struct {
	buf [3]byte
	idx int
}

func NewExitBuf() *ExitBuf {
	return &ExitBuf{
		buf: [3]byte{},
		idx: 0,
	}
}

func (eb *ExitBuf) Insert(input byte) {
	eb.buf[eb.idx] = input
	eb.idx++
	if eb.idx > 2 {
		eb.idx = 0
	}
}

func (eb *ExitBuf) ShouldExit() bool {
	return eb.buf[0] == KeyEscape && eb.buf[1] == ':' && eb.buf[2] == 'q'
}

type CursorPos struct {
	row int
	col int
}

func getCursorPos() (CursorPos, error) {
	fmt.Print(RequestCursorPos)
	buf := make([]byte, 32)
	n, err := os.Stdin.Read(buf)
	if err != nil {
		return CursorPos{}, err
	}
	response := string(buf[:n])
	if !strings.HasPrefix(response, "\033[") || !strings.HasSuffix(response, "R") {
		return CursorPos{}, fmt.Errorf("unexpected response format: %s", response)
	}

	trimmed := strings.TrimSuffix(strings.TrimPrefix(response, "\033["), "R")
	parts := strings.Split(trimmed, ";")
	if len(parts) != 2 {
		return CursorPos{}, fmt.Errorf("unexpected response format: %s", response)
	}

	row, err := strconv.Atoi(parts[0])
	if err != nil {
		return CursorPos{}, err
	}

	col, err := strconv.Atoi(parts[1])
	if err != nil {
		return CursorPos{}, err
	}

	return CursorPos{
		row: row,
		col: col,
	}, nil
}

func handleBackspace() error {
	pos, err := getCursorPos()
	if err != nil {
		return err
	}
	if pos.col == 1 && pos.row > 1 {
		fmt.Print(MoveCursorUp)

		// because we don't have a way to know the current line length,
		// we can't move the cursor to the correct position, which is usually
		// the last character of the line
	}
	fmt.Print(MoveCursorLeft)
	fmt.Printf("%c", KeySpace)
	fmt.Print(MoveCursorLeft)
	return nil

}

func main() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	fmt.Printf(ClearScreen)
	fmt.Printf("\033[0;0H")

	exitBuf := NewExitBuf()

	for {
		buf := make([]byte, 1)
		_, err = os.Stdin.Read(buf)
		if err != nil {
			fmt.Printf("stdin.Read(): %v\n", err)
			return
		}

		input := buf[0]

		exitBuf.Insert(input)
		if exitBuf.ShouldExit() {
			return
		}

		if input == KeyEnter {
			fmt.Print(CRLF)
			continue
		}

		if input == KeyBackspace {
			if err = handleBackspace(); err != nil {
				return
			}
			continue
		}

		fmt.Printf("%c", input)
	}
}
