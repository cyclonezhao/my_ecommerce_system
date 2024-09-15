package microservice

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	my_client "my_ecommerce_system/pkg/client"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	gresolver "google.golang.org/grpc/resolver"
)

const (
	nameServicePrefix = "my_ecommerce_system/mservice"
	// 默认的租赁时间
	defLeaseSecond = 30
)

var etcdUrl string // etcd的url地址，如：http://localhost:2379
var etcdHost string
var etcdPort int
var etcdClient *clientv3.Client           // etcd 客户端，用于与 etcd 通信
var etcdEndpointManager endpoints.Manager // etcd 的 endpoints 管理器，用于处理服务的注册和注销
var serviceInstance *ServiceInstance
var leaseId clientv3.LeaseID // 租约id，也是服务实例id

// 服务元数据信息
type ServiceInstance struct {
	ServiceName string `json:"service_name"` // 服务名称，比如 my_system，作为endpoint的倒数第二级标识（末级是租约id，也作为实例id）
	Addr        string `json:"addr"`
	Port        int    `json:"port"`
}

// 创建一个新的命名服务
//
// 这个注册器参考了(https://blog.csdn.net/small_to_large/article/details/130656230)，按本项目情况做了适配性改动
func RegisterSelf(serviceInstance *ServiceInstance) error {
	etcdClientWrapper := my_client.EtcdClientWrapper
	etcdClient = etcdClientWrapper.EtcdClient
	etcdUrl = etcdClientWrapper.EtcdUrl
	etcdHost = etcdClientWrapper.EtcdHost
	etcdPort = etcdClientWrapper.EtcdPort

	// etcd的endpoints管理
	var err error
	etcdEndpointManager, err = endpoints.NewManager(etcdClient, nameServicePrefix)
	if err != nil {
		return err
	}

	err = addEndpoint(serviceInstance)
	if err != nil {
		return err
	}

	return nil
}

func GetFullServerName(name string, leaseID clientv3.LeaseID) string {
	return fmt.Sprintf("%s/%s/%d", nameServicePrefix, name, leaseID)
}

func GetPathServerName(name string) string {
	return fmt.Sprintf("etcd://%s:%d/%s/%s", etcdHost, etcdPort, nameServicePrefix, name)
}

// 添加实例到Etcd
func addEndpoint(e *ServiceInstance) error {
	b, _ := json.Marshal(e)
	ep := endpoints.Endpoint{
		Addr:     fmt.Sprintf("%s:%d", e.Addr, e.Port),
		Metadata: string(b), // 这里有个坑，必须传入字符串，没来得及看原因
	}
	// 在etcd创建一个续期的lease对象
	lease, err := etcdClient.Grant(context.TODO(), defLeaseSecond)
	if err != nil {
		return err
	}

	serviceInstance = e
	leaseId = lease.ID

	key := GetFullServerName(e.ServiceName, leaseId)
	// 向etcd注册一个Endpoint并绑定续期
	err = etcdEndpointManager.AddEndpoint(context.TODO(), key, ep, clientv3.WithLease(leaseId))
	if err != nil {
		return err
	}
	log.Println("AddEndpoint success:", key)
	// 开启自动续期KeepAlive
	ch, err := etcdClient.KeepAlive(context.TODO(), leaseId)
	if err != nil {
		return err
	}

	// 这个方法会异步打印出每次续期调用的日志
	go func() {
		for ka := range ch {
			log.Printf("%v 自动续期: %d", ka.ID, ka.TTL)
		}
		log.Println("终止自动续期: " + key)
	}()

	return nil
}

// 从Etcd中移除自身
func DelEndpoint() error { // 从 etcd 中删除指定的 Endpoint
	err := etcdEndpointManager.DeleteEndpoint(context.TODO(), GetFullServerName(serviceInstance.ServiceName, leaseId))
	if err != nil {
		log.Fatalf("Delete endpoint error %v", err)
		return err
	}

	// 撤销租约，确保本客户端不会继续往 etcd 中续约
	_, err = etcdClient.Revoke(context.TODO(), leaseId)
	if err != nil {
		log.Fatalf("Revoke lease error %v", err)
		return err
	}

	log.Println("从Etcd中移除自身成功！")
	return nil
}

// 创建一个 etcd 的 gRPC 解析器（resolver）。通过它，客户端能够自动发现和连接到注册的 gRPC 服务
func NewEtcdResolver() (gresolver.Builder, error) {
	etcdClientWrapper := my_client.EtcdClientWrapper
	etcdClient = etcdClientWrapper.EtcdClient
	etcdResolver, err := resolver.NewBuilder(etcdClient)
	if err != nil {
		log.Fatalf("Etcd resolver error %v", err)
		return nil, err
	}
	return etcdResolver, nil
}
