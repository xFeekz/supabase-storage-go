package test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"testing"

	storage_go "github.com/xFeekz/supabase-storage-go"
)

var (
	rawUrl = "https://abc.supabase.co/storage/v1"
	token  = ""
	apiKey = ""
)

func TestUpload(t *testing.T) {
	file, err := os.Open("dummy.txt")
	if err != nil {
		panic(err)
	}
	c := storage_go.NewClient(rawUrl, token, apiKey, map[string]string{})
	resp, err := c.UploadFile("test", "test.txt", file)
	fmt.Println(resp, err)

	// resp, err = c.UploadFile("test", "hola.txt", []byte("hello world"))
	// fmt.Println(resp, err)
}

func TestConcurrentUploadRace(t *testing.T) {
	file, err := os.Open("dummy.txt")
	if err != nil {
		t.Fatal(err)
	}
	bf, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	c := storage_go.NewClient(rawUrl, token, apiKey, map[string]string{})

	ct := "application/text"

	var wg sync.WaitGroup
	n := 2
	x := true

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			resp, err := c.UploadFile("test", fmt.Sprintf("race_%v.txt", i), bytes.NewReader(bf), storage_go.FileOptions{
				ContentType: &ct,
				Upsert: &x,
			})
			if err != nil {
				fmt.Printf("err: %v\n", err)
			} else {
				fmt.Printf("resp: %+v\n", resp)
			}
		}(i)
	}

	wg.Wait()
}

func TestUpdate(t *testing.T) {
	file, err := os.Open("dummy.txt")
	if err != nil {
		panic(err)
	}
	c := storage_go.NewClient(rawUrl, token, apiKey, map[string]string{})
	resp, err := c.UpdateFile("test", "test.txt", file)

	fmt.Println(resp, err)
}

func TestMoveFile(t *testing.T) {
	c := storage_go.NewClient(rawUrl, token, apiKey, map[string]string{})
	resp, err := c.MoveFile("test", "test.txt", "random/test.txt")

	fmt.Println(resp, err)
}

func TestSignedUrl(t *testing.T) {
	c := storage_go.NewClient(rawUrl, token, apiKey, map[string]string{})
	resp, err := c.CreateSignedUrl("test", "test.txt", 120)

	fmt.Println(resp, err)
}

func TestPublicUrl(t *testing.T) {
	c := storage_go.NewClient(rawUrl, token, apiKey, map[string]string{})
	resp := c.GetPublicUrl("shield", "book.pdf")

	fmt.Println(resp)
}

func TestDeleteFile(t *testing.T) {
	c := storage_go.NewClient(rawUrl, token, apiKey, map[string]string{})
	resp, err := c.RemoveFile("shield", []string{"book.pdf"})

	fmt.Println(resp, err)
}

func TestListFile(t *testing.T) {
	c := storage_go.NewClient(rawUrl, token, apiKey, map[string]string{})
	resp, err := c.ListFiles("shield", "", storage_go.FileSearchOptions{
		Limit:  10,
		Offset: 0,
		SortByOptions: storage_go.SortBy{
			Column: "",
			Order:  "",
		},
	})

	fmt.Println(resp, err)
}

func TestCreateUploadSignedUrl(t *testing.T) {
	c := storage_go.NewClient(rawUrl, token, apiKey, map[string]string{"apiKey": apiKey})
	resp, err := c.CreateSignedUploadUrl("your-bucket-id", "book.pdf")

	fmt.Println(resp, err)
}

func TestUploadToSignedUrl(t *testing.T) {
	c := storage_go.NewClient(rawUrl, token, apiKey, map[string]string{"apiKey": apiKey})
	file, err := os.Open("dummy.txt")
	if err != nil {
		panic(err)
	}
	// resp, err := c.CreateSignedUploadUrl("test", "vu.txt")
	res, err := c.UploadToSignedUrl("your-response-url", file)
	fmt.Println(res, err)
}

func TestDownloadFile(t *testing.T) {
	c := storage_go.NewClient(rawUrl, token, apiKey, map[string]string{})
	resp, err := c.DownloadFile("test", "book.pdf")
	if err != nil {
		t.Fatalf("DownloadFile failed: %v", err)
	}

	err = os.WriteFile("book.pdf", resp, 0644)
	if err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
}