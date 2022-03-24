package k8s

import (
	"context"
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/golang/glog"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

type PrometheusConfig struct {
	Global struct {
		ScrapeInterval     string `yaml:"scrape_interval"`
		EvaluationInterval string `yaml:"evaluation_interval"`
	} `yaml:"global"`
	RuleFiles     []string `yaml:"rule_files"`
	ScrapeConfigs []struct {
		JobName       string `yaml:"job_name"`
		StaticConfigs []struct {
			Targets []string `yaml:"targets"`
		} `yaml:"static_configs"`
		HttpSdConfigs []struct {
			Url string `yaml:"url"`
		} `yaml:"http_sd_configs"`
	} `yaml:"scrape_configs"`
}

type RulesConfig struct {
	Groups []struct {
		Name  string `yaml:"name"`
		Rules []Rule `yaml:"rules"`
	} `yaml:"groups"`
}

type Rule struct {
	Alert  string `yaml:"alert"`
	Expr   string `yaml:"expr"`
	For    string `yaml:"for"`
	Labels struct {
		User string `yaml:"user"`
	} `yaml:"labels"`
	Annotations struct {
		Summary     string `yaml:"summary"`
		Description string `yaml:"description"`
	} `yaml:"annotations"`
}

func PodsTest() {
	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)
	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		glog.Fatalf("Error building kubeconfig: %fileData", err.Error())
	}
	mclient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}
	configProm(mclient)
	configPromRules(mclient)

}
func configPromRules(clientset *kubernetes.Clientset) {
	configmap, err := clientset.CoreV1().ConfigMaps("default").Get(context.Background(), "prometheus-rules-config", metav1.GetOptions{})
	if err != nil {
		logger.Error("get config map with error: ", err.Error())
		return
	}

	fileData := configmap.Data["prometheus-rules-config.yml"]
	fmt.Println("file data:\n ", fileData)

	structData := &RulesConfig{}
	err = yaml.Unmarshal([]byte(fileData), structData)
	if err != nil {
		logger.Error("unmarshal with error: ", err.Error())
		return
	}

	fmt.Println("struct data:\n ", structData)
}
func configProm(clientset *kubernetes.Clientset) {
	configmap, err := clientset.CoreV1().ConfigMaps("default").Get(context.Background(), "prometheus-config", metav1.GetOptions{})
	if err != nil {
		logger.Error("get config map with error: ", err.Error())
		return
	}

	fileData := configmap.Data["prometheus.yml"]
	fmt.Println("file data:\n ", fileData)

	structData := &PrometheusConfig{}
	err = yaml.Unmarshal([]byte(fileData), structData)
	if err != nil {
		logger.Error("unmarshal with error: ", err.Error())
		return
	}

	for index, c := range structData.ScrapeConfigs {
		if c.JobName == "prometheus" {
			structData.ScrapeConfigs = append(structData.ScrapeConfigs[:index], structData.ScrapeConfigs[index+1:]...)
			break
		}
	}

	byteConfig, err := yaml.Marshal(structData)
	if err != nil {
		fmt.Println("marshal struct data to byte with error: ", err.Error())
		return
	}
	configmap.Data["prometheus.yml"] = string(byteConfig)

	configmap, err = clientset.CoreV1().ConfigMaps("default").Update(context.Background(), configmap, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("update config map with error: ", err.Error())
		return
	}
	fmt.Println("struct data:\n ", structData)
}
