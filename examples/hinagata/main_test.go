package main

import (
	"flag"
	"io"
	"os"
	"reflect"
	"testing"
)

func TestParseFlags(t *testing.T) {
	testcases := []struct {
		name    string
		args    []string
		want    Config
		isError bool
	}{
		{
			name: "Simple",
			args: []string{"cmd", "-host", "host1", "-hostname", "1.2.3.4"},
			want: Config{
				ConfPath: "/home/app/.ssh/config",
				SshConfig: SshConfig{
					Host:     "host1",
					Hostname: "1.2.3.4",
					User:     "root",
				},
			},
			isError: false,
		},
		{
			name: "SetCustomConfigPath",
			args: []string{"cmd", "-c", "/tmp/ssh_config", "-host", "host2", "-hostname", "5.6.7.8"},
			want: Config{
				ConfPath: "/tmp/ssh_config",
				SshConfig: SshConfig{
					Host:     "host2",
					Hostname: "5.6.7.8",
					User:     "root",
				},
			},
			isError: false,
		},
		{
			name: "SetCustomUser",
			args: []string{"cmd", "-host", "host3", "-hostname", "9.10.11.12", "-user", "hoge"},
			want: Config{
				ConfPath: "/home/app/.ssh/config",
				SshConfig: SshConfig{
					Host:     "host3",
					Hostname: "9.10.11.12",
					User:     "hoge",
				},
			},
			isError: false,
		},
	}

	for _, tc := range testcases {
		t.Run(
			tc.name,
			func(t *testing.T) {
				os.Args = tc.args

				flag.CommandLine = flag.NewFlagSet("t1", flag.ContinueOnError)
				got, err := parseFlags()
				if (err != nil) != tc.isError {
					t.Errorf("error: %v, but isError: %v\n", err, tc.isError)
				}
				want := tc.want
				if !reflect.DeepEqual(got, want) {
					t.Errorf("got: %v, want: %v\n", got, want)
				}
			},
		)
	}
}

func TestRenderTemplate(t *testing.T) {
	testcases := []struct {
		name    string
		conf    SshConfig
		want    string
		isError bool
	}{
		{
			name: "Simple",
			conf: SshConfig{"host1", "1.2.3.4", "root"},
			want: "Host host1\n    Hostname 1.2.3.4\n    User root\n\n",
		},
		{
			name: "SetCustomUser",
			conf: SshConfig{"host2", "5.6.7.8", "hoge"},
			want: "Host host2\n    Hostname 5.6.7.8\n    User hoge\n\n",
		},
		{
			name: "UseDNSHostnameAsHostname",
			conf: SshConfig{"host3", "fuga.example.com", "fuga"},
			want: "Host host3\n    Hostname fuga.example.com\n    User fuga\n\n",
		},
	}

	for _, tc := range testcases {
		t.Run(
			tc.name,
			func(t *testing.T) {
				got, err := renderTemplate(tc.conf)
				if (err != nil) != tc.isError {
					t.Errorf("error: %v, but isError: %v\n", err, tc.isError)
				}
				want := tc.want
				if got != want {
					t.Errorf("got: %v, want: %v\n", got, want)
				}
			},
		)
	}
}

func TestWriteFile_Success(t *testing.T) {
	testcases := []struct {
		name       string
		content    string
		outputPath string
		want       string
		isError    bool
	}{
		{
			name:       "Simple",
			content:    "Host host1\n    Hostname 1.2.3.4\n    User root\n\n",
			outputPath: "/home/app/.ssh/config",
			want:       "Host host1\n    Hostname 1.2.3.4\n    User root\n\n",
			isError:    false,
		},
	}

	for _, tc := range testcases {
		t.Run(
			tc.name,
			func(t *testing.T) {
				if err := writeFile(tc.content, tc.outputPath); (err != nil) != tc.isError {
					t.Errorf("error: %v, but isError: %v\n", err, tc.isError)
				}
				file, err := os.Open(tc.outputPath)
				if (err != nil) != tc.isError {
					t.Errorf("error: %v, but isError: %v\n", err, tc.isError)
				}
				defer func() {
					file.Close()
					os.Remove(tc.outputPath)
				}()

				fi, err := file.Stat()
				if (err != nil) != tc.isError {
					t.Errorf("error: %v, but isError: %v\n", err, tc.isError)
				}
				if fi.Mode() != os.FileMode(0o600) {
					t.Errorf("got: %v, want: %v\n", fi.Mode(), os.FileMode(0o600))
				}
			},
		)
	}
}

func TestWriteFile_Failure(t *testing.T) {
	t.Run(
		"FileAlreadyExists",
		func(t *testing.T) {
			file, err := os.CreateTemp("/tmp", "TestWriteFile_Failure_")
			if err != nil {
				t.Error(err)
			}
			defer func() {
				file.Close()
				os.Remove(file.Name())
			}()

			content := "hogehoge\n"
			_, err = file.WriteString(content)
			if err != nil {
				t.Error(err)
			}
			_, err = file.Seek(0, 0)
			if err != nil {
				t.Error(err)
			}

			if err := writeFile("Host host1\n    Hostname 1.2.3.4\n    User root\n\n", file.Name()); err == nil {
				t.Error("File exists, but no errors\n")
			}
			b, err := io.ReadAll(file)
			if err != nil {
				t.Error(err)
			}
			if content != string(b) {
				t.Error("File Changed!")
			}
		},
	)

}
