package params

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Params struct {
	LogLevel  string          `yaml:"loglevel"`
	Redis     string          `yaml:"redis"`
	Publisher PublisherParams `yaml:"publisher"`
	Worker    WorkerParams    `yaml:"worker"`
	Testing   Testing         `yaml:"testing"`
}
type PublisherParams struct {
}
type WorkerParams struct {
	Concurrency int `yaml:"concurrency"`
}

type Testing struct {
	JobMaxTasks int `yaml:"job_max_tasks"`
	NumJobs     int `yaml:"num_jobs"`
}

func (p *Params) ParseArgs() {

	config := flag.String("config", "config.yaml", "Path to config.yaml file")
	loglevel := flag.String("log", "", "Available debug levels: panic, fatal, error, warn, info, debug, trace")
	flag.Parse()

	p.getConf(*config)
	if *loglevel != "" {
		p.LogLevel = *loglevel
	}

	level, err := logrus.ParseLevel(p.LogLevel)
	if err != nil {
		level = logrus.WarnLevel
		logrus.Warnf("Unknown log level %s. Defaulting to warn", p.LogLevel)
	}

	logrus.SetLevel(level)
}

func (p *Params) getConf(conf string) {

	yamlFile, err := ioutil.ReadFile(conf)
	if err != nil {
		logrus.Fatalf("yamlFile.Get err   #%v ", err)
		panic(err)
	}
	logrus.Infof("Loading config \n%s", yamlFile)
	err = yaml.Unmarshal(yamlFile, p)
	if err != nil {
		logrus.Fatalf("Unmarshal: %v", err)
		panic(err)
	}
	fmt.Println(p)
}
