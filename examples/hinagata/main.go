package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
)

type Options struct {
	Host     string
	Hostname string
	User     string
}

func main() {
	conf, err := parseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}

	content, err := renderTemplate(conf.SshConfig)
	if err != nil {
		log.Fatal(err)
	}

	if err := writeFile(content, conf.ConfPath); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ssh_config created.\npath: %v\n\n", conf.ConfPath)
}

type SshConfig struct {
	Host     string
	Hostname string
	User     string
}
type Config struct {
	ConfPath  string
	SshConfig SshConfig
}

func parseFlags() (Config, error) {
	var confPath string
	var host, hostname, user string
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}
	flag.StringVar(&confPath, "c", homeDir+"/.ssh/config", "Set Config file `Path`.")
	flag.StringVar(&host, "host", "", "Set `Host`.")
	flag.StringVar(&hostname, "hostname", "", "Set `Hostname`.")
	flag.StringVar(&user, "user", "root", "Set `User`.")
	flag.Parse()

	if host == "" {
		return Config{}, errors.New("-host option is requred")
	}
	if hostname == "" {
		return Config{}, errors.New("-hostname option is requred")
	}

	conf := Config{
		ConfPath: confPath,
		SshConfig: SshConfig{
			Host:     host,
			Hostname: hostname,
			User:     user,
		},
	}

	return conf, nil

}

func renderTemplate(conf SshConfig) (string, error) {
	tmpl, err := template.New("ssh_config").Parse(`Host {{ .Host }}
    Hostname {{ .Hostname }}
    User {{ .User }}

`)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, conf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func writeFile(content, outputPath string) error {
	perm := os.FileMode(0o600)
	file, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, perm)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}
