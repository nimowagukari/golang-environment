package main

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"reflect"
	"testing"
)

func TestReadBytesFromFile_NoSuchFileOrDirectory(t *testing.T) {
	path := "/tmp/non_exists.txt"
	_, got := readBytesFromFile(path)

	var want *fs.PathError
	if !errors.As(got, &want) {
		t.Errorf("got: %#v, want: %#v\n", got, want)
	}
}

func TestReadBytesFromFile_PermissionDenied(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "coff_permission_denied_*.txt")
	if err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	if err := tmpfile.Chmod(0o000); err != nil {
		t.Error(err)
		return
	}

	_, got := readBytesFromFile(tmpfile.Name())
	var want *fs.PathError
	if !errors.As(got, &want) {
		t.Errorf("got: %#v, want: %#v\n", got, want)
	}
}

func TestSplitByNewLine(t *testing.T) {
	testcases := []struct {
		name string
		arg  []byte
		want [][]byte
	}{
		{
			name: "Positive1",
			arg:  []byte("hoge\nfuga\ngeso\n"),
			want: [][]byte{
				[]byte("hoge"),
				[]byte("fuga"),
				[]byte("geso"),
				[]byte(""),
			},
		},
	}

	for _, tc := range testcases {
		t.Run(
			tc.name,
			func(t *testing.T) {
				got := splitByNewLine(tc.arg)
				want := tc.want
				if !reflect.DeepEqual(got, want) {
					t.Errorf("\ngot: \n%s\nwant: \n%s\n", bytes.Join(got, []byte("\n")), bytes.Join(want, []byte("\n")))
				}
			},
		)
	}

}

func TestRemoveComment(t *testing.T) {
	testcases := []struct {
		name string
		arg  [][]byte
		want [][]byte
	}{
		{
			name: "Positive1",
			arg: [][]byte{
				[]byte("# removed line"),
				[]byte("keeped line"),
				[]byte("    # removed line"),
				[]byte("		# removed line"),
				[]byte("    keeped line"),
				[]byte("　　　# keeped line"),
				[]byte(""),
			},
			want: [][]byte{
				[]byte(""),
				[]byte("keeped line"),
				[]byte(""),
				[]byte(""),
				[]byte("    keeped line"),
				[]byte("　　　# keeped line"),
				[]byte(""),
			},
		},
	}

	for _, tc := range testcases {
		t.Run(
			tc.name,
			func(t *testing.T) {
				got, err := removeComment(tc.arg)
				if err != nil {
					t.Error(err)
				}
				want := tc.want
				if !reflect.DeepEqual(got, want) {
					t.Errorf("\ngot: \n%s\nwant: \n%s\n", bytes.Join(got, []byte("\n")), bytes.Join(want, []byte("\n")))
				}
			},
		)
	}

}
