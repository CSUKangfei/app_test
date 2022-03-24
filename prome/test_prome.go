package prome

import (
	"context"
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"time"
)

var (
	NodeCpuCoreTotal = "count(node_cpu_seconds_total{instance=~\"%s\", mode='system'})"
	IDCCpuCoreTotal  = "count(node_cpu_seconds_total{idc_id=~\"%s\", mode='system'})"
	NodeCpuUsage     = "(1 - avg(irate(node_cpu_seconds_total{instance=~\"%s\",mode=\"idle\"}[5m])) by (instance))*100"
	IDCCpuUsage      = "(1 - avg(irate(node_cpu_seconds_total{idc_id=~\"%s\",mode=\"idle\"}[5m])) by (instance))*100" //todo

	NodeMemTotal = "sum(node_memory_MemTotal_bytes{instance=~\"%s\"})"
	IDCMemTotal  = "sum(node_memory_MemTotal_bytes{idc_id=~\"%s\"})"
	NodeMemUsed  = "node_memory_MemTotal_bytes{instance=~\"%s\"} - node_memory_MemAvailable_bytes{instance=~\"%s\"}"
	IdcMemUsed   = "node_memory_MemTotal_bytes{idc_id=~\"%s\"} - node_memory_MemAvailable_bytes{idc_id=~\"%s\"}"

	NodeRootDiskTotal = "node_filesystem_files{instance=~\"%s\",fstype=~\"ext.?|xfs\",mountpoint=\"/\"}"
	IdcRootDiskTotal  = "sum(node_filesystem_files{idc_id=~\"%s\",fstype=~\"ext.?|xfs\",mountpoint=\"/\"})"
	NodeRootDiskFree  = "node_filesystem_files_free{instance=~\"%s\",fstype=~\"ext.?|xfs\",mountpoint=\"/\"}"
	IdcRootDiskFree   = "sum(node_filesystem_files_free{idc_id=~\"%s\",fstype=~\"ext.?|xfs\",mountpoint=\"/\"})"

	PodTest = "container_memory_usage_bytes{name!=\"\"}"
)

func PromeData() {
	config := api.Config{
		//Address: "http://10.90.138.203:9090",
		Address: "http://localhost:9090",
	}
	client, err := api.NewClient(config)
	if err != nil {
		logger.Error("new client with error: %v", err.Error())
		return
	}
	clientApi := v1.NewAPI(client)

	result, _, err := clientApi.Query(context.Background(), fmt.Sprintf(PodTest), time.Now())
	if err != nil {
		logger.Error("PodTest: ", err.Error())
		return
	}
	v, _ := result.(model.Vector)
	logger.Info("PodTest: ", v.String())

	result, _, err = clientApi.Query(context.Background(), fmt.Sprintf(NodeCpuCoreTotal, "10.90.139.192:9100"), time.Now())
	if err != nil {
		logger.Error("NodeCpuCoreTotal: ", err.Error())
		return
	}
	v, _ = result.(model.Vector)
	logger.Info("NodeCpuCoreTotal: ", int(v[0].Value))

	result, _, err = clientApi.Query(context.Background(), fmt.Sprintf(IDCCpuCoreTotal, "9"), time.Now())
	if err != nil {
		logger.Error("IDCCpuCoreTotal: ", err.Error())
		return
	}
	v, _ = result.(model.Vector)
	logger.Info("IDCCpuCoreTotal: ", int(v[0].Value))

	result, _, err = clientApi.Query(context.Background(), fmt.Sprintf(NodeCpuUsage, "10.90.138.149:9100"), time.Now())
	if err != nil {
		logger.Error("query with error: %s", err.Error())
		return
	}
	v, _ = result.(model.Vector)
	logger.Info("NodeCpuUsage: ", int(v[0].Value))

	result, _, err = clientApi.Query(context.Background(), fmt.Sprintf(IDCCpuUsage, "9"), time.Now())
	if err != nil {
		logger.Error("query with error: %s", err.Error())
		return
	}
	v, _ = result.(model.Vector)
	logger.Info("IDCCpuUsage: ", int(v[0].Value))

	result, _, err = clientApi.Query(context.Background(), fmt.Sprintf(NodeMemTotal, "10.90.139.192:9100"), time.Now())
	if err != nil {
		logger.Error("query with error: %s", err.Error())
		return
	}
	v, _ = result.(model.Vector)
	logger.Info("NodeMemTotal: ", int(v[0].Value))

	result, _, err = clientApi.Query(context.Background(), fmt.Sprintf(IDCMemTotal, "9"), time.Now())
	if err != nil {
		logger.Error("query with error: %s", err.Error())
		return
	}
	v, _ = result.(model.Vector)
	logger.Info("IDCMemTotal: ", int(v[0].Value))

	result, _, err = clientApi.Query(context.Background(), fmt.Sprintf(NodeMemUsed, "10.90.139.192:9100", "10.90.139.192:9100"), time.Now())
	if err != nil {
		logger.Error("query with error: %s", err.Error())
		return
	}
	v, _ = result.(model.Vector)
	logger.Info("node mem used: ", int(v[0].Value))

	result, _, err = clientApi.Query(context.Background(), fmt.Sprintf(IdcMemUsed, "9", "9"), time.Now())
	if err != nil {
		logger.Error("query with error: %s", err.Error())
		return
	}
	v, _ = result.(model.Vector)
	logger.Info("idc_id mem used: ", int(v[0].Value))

	result, _, err = clientApi.Query(context.Background(), fmt.Sprintf(NodeRootDiskTotal, "10.90.139.192:9100"), time.Now())
	if err != nil {
		logger.Error("query with error: %s", err.Error())
		return
	}
	v, _ = result.(model.Vector)
	logger.Info("NodeRootDiskTotal: ", int(v[0].Value))
	result, _, err = clientApi.Query(context.Background(), fmt.Sprintf(IdcRootDiskTotal, "9"), time.Now())
	if err != nil {
		logger.Error("query with error: %s", err.Error())
		return
	}
	v, _ = result.(model.Vector)
	logger.Info("IdcRootDiskTotal: ", int(v[0].Value))
	result, _, err = clientApi.Query(context.Background(), fmt.Sprintf(NodeRootDiskFree, "10.90.139.192:9100"), time.Now())
	if err != nil {
		logger.Error("query with error: %s", err.Error())
		return
	}
	v, _ = result.(model.Vector)
	logger.Info("NodeRootDiskFree: ", int(v[0].Value))
	result, _, err = clientApi.Query(context.Background(), fmt.Sprintf(IdcRootDiskFree, "9"), time.Now())
	if err != nil {
		logger.Error("query with error: %s", err.Error())
		return
	}
	v, _ = result.(model.Vector)
	logger.Info("IdcRootDiskFree: ", int(v[0].Value))
	//switch result.Type() {
	//case model.ValNone:
	//	fmt.Println("None Type")
	//case model.ValScalar:
	//	fmt.Println("Scalar Type")
	//	v, _ := result.(*model.Scalar)
	//	displayScalar(v)
	//case model.ValVector:
	//	fmt.Println("Vector Type")
	//	v, _ := result.(model.Vector)
	//	displayVector(v)
	//case model.ValMatrix:
	//	fmt.Println("Matrix Type")
	//	v, _ := result.(model.Matrix)
	//	displayMatrix(v)
	//case model.ValString:
	//	fmt.Println("String Type")
	//	v, _ := result.(*model.String)
	//	displayString(v)
	//default:
	//	fmt.Printf("Unknow Type")
	//}
	//rangeTime := v1.Range{
	//	Start: time.Unix(1643186520, 0),
	//	End:   time.Unix(1643186530, 0),
	//	Step:  1,
	//}
}