// Copyright 2015 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tango

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestStatic(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	tg := New()
	tg.Use(Static())

	req, err := http.NewRequest("GET", "http://localhost:8000/test.html", nil)
	if err != nil {
		t.Error(err)
	}

	tg.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "hello tango")

	buff.Reset()

	req, err = http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	tg.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "this is index.html")
}

func TestStatic2(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	tg := New()
	tg.Use(Static())

	req, err := http.NewRequest("GET", "http://localhost:8000/test.png", nil)
	if err != nil {
		t.Error(err)
	}

	tg.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusNotFound)
}

func TestStatic3(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	tg := New()
	tg.Use(Static(StaticOptions{
		Prefix:   "/public",
		RootPath: "./public",
	}))

	req, err := http.NewRequest("GET", "http://localhost:8000/public/test.html", nil)
	if err != nil {
		t.Error(err)
	}

	tg.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	expect(t, buff.String(), "hello tango")
}

func TestStatic4(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	tg := New()
	tg.Use(Static(StaticOptions{
		Prefix:   "/public",
		RootPath: "./public",
	}))

	req, err := http.NewRequest("GET", "http://localhost:8000/public/t.html", nil)
	if err != nil {
		t.Error(err)
	}

	tg.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusNotFound)
	expect(t, buff.String(), NotFound().Error())
}

func TestStatic5(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	tg := New()
	tg.Use(Static(StaticOptions{
		Prefix:     "/public",
		RootPath:   "./public",
		ListDir:    true,
		IndexFiles: []string{"a.html"},
	}))

	req, err := http.NewRequest("GET", "http://localhost:8000/public/", nil)
	if err != nil {
		t.Error(err)
	}

	tg.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
}

type MemoryFileSystem map[string][]byte

type MemoryFile struct {
	Name  string
	isDir bool
	*bytes.Reader
}

func (m *MemoryFile) Close() error {
	return nil
}

func (m *MemoryFile) Readdir(count int) ([]os.FileInfo, error) {
	var infos []os.FileInfo
	err := json.NewDecoder(m.Reader).Decode(&infos)
	if err != nil {
		return nil, err
	}
	if count > 0 && count < len(infos) {
		return infos[:count], nil
	}
	return infos, nil
}

type MemoryFileInfo struct {
	name string
	size int64
	time.Time
	isDir bool
}

func (m *MemoryFileInfo) Name() string {
	return m.name
}

func (m *MemoryFileInfo) Size() int64 {
	return m.size
}

func (m *MemoryFileInfo) Mode() os.FileMode {
	return os.ModePerm
}

func (m *MemoryFileInfo) ModTime() time.Time {
	return m.Time
}

func (m *MemoryFileInfo) IsDir() bool {
	return m.isDir
}

func (m *MemoryFileInfo) Sys() interface{} {
	return nil
}

func (m *MemoryFile) Stat() (os.FileInfo, error) {
	return &MemoryFileInfo{
		name:  m.Name,
		size:  int64(m.Len()),
		Time:  time.Now(),
		isDir: m.isDir,
	}, nil
}

var (
	_ http.FileSystem = &MemoryFileSystem{}
	_ http.File       = &MemoryFile{}
)

func (m MemoryFileSystem) Open(name string) (http.File, error) {
	if name == "/" || name == "" {
		var finfos []os.FileInfo
		for k, v := range m {
			finfos = append(finfos, &MemoryFileInfo{
				name:  k,
				size:  int64(len(v)),
				Time:  time.Now(),
				isDir: v[0] == '[',
			})
		}

		bs, err := json.Marshal(finfos)
		if err != nil {
			return nil, err
		}

		return &MemoryFile{
			Name:   "/",
			isDir:  true,
			Reader: bytes.NewReader(bs),
		}, nil
	}

	bs, ok := m[name]
	if !ok {
		return nil, os.ErrNotExist
	}

	return &MemoryFile{
		Name:   name,
		isDir:  bs[0] == '[',
		Reader: bytes.NewReader(bs),
	}, nil
}

func TestStatic6(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	var myFileSystem = MemoryFileSystem{
		"a.html": []byte("<html></html>"),
	}
	tg := New()
	tg.Use(Static(StaticOptions{
		Prefix:     "/public",
		RootPath:   "./public",
		ListDir:    false,
		IndexFiles: []string{"a.html"},
		FileSystem: myFileSystem,
	}))

	req, err := http.NewRequest("GET", "http://localhost:8000/public/", nil)
	if err != nil {
		t.Error(err)
	}

	tg.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
}

func TestStatic7(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	tg := New()
	tg.Use(Static(StaticOptions{
		RootPath: "./public",
	}))

	req, err := http.NewRequest("GET", "http://localhost:8000/test.html", nil)
	if err != nil {
		t.Error(err)
	}

	tg.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	expect(t, buff.String(), "hello tango")
}

func TestStatic8(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	tg := New()
	tg.Use(Return())
	tg.Use(Static(StaticOptions{
		RootPath: "./public",
	}))
	tg.Get("/b", func() string {
		return "hello"
	})

	req, err := http.NewRequest("GET", "http://localhost:8000/test.html", nil)
	if err != nil {
		t.Error(err)
	}

	tg.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	expect(t, buff.String(), "hello tango")

	buff.Reset()

	req, err = http.NewRequest("GET", "http://localhost:8000/b", nil)
	if err != nil {
		t.Error(err)
	}

	tg.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	expect(t, buff.String(), "hello")
}
