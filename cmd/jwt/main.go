package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func run() error {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	index := bytes.IndexByte(input, '.')

	header := make([]byte, base64.RawStdEncoding.DecodedLen(len(input[:index])))
	if _, err := base64.RawStdEncoding.Decode(header, input[:index]); err != nil {
		return fmt.Errorf("decoding header: %w", err)
	}

	fmt.Println("=== HEADER ===")
	var b bytes.Buffer
	if json.Indent(&b, header, "", "  ") == nil {
		_, _ = io.Copy(os.Stdout, &b)
	} else {
		_, _ = os.Stdout.Write(header)
	}
	fmt.Println("\n=== BODY ===")

	input = input[index+1:]
	index = bytes.IndexByte(input, '.')

	body := make([]byte, base64.RawStdEncoding.DecodedLen(len(input[:index])))
	if _, err := base64.RawStdEncoding.Decode(body, input[:index]); err != nil {
		return fmt.Errorf("decoding body: %w", err)
	}

	b.Reset()

	if json.Indent(&b, body, "", "  ") == nil {
		_, _ = io.Copy(os.Stdout, &b)
	} else {
		_, _ = os.Stdout.Write(body)
	}

	fmt.Println("\n=== SIGNATURE ===")
	input = input[index+1:]
	_, _ = os.Stdout.Write(input)
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
