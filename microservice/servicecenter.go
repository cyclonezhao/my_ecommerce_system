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
	nameServicePrefix = "my-ecommerce-system"
	// 默认的租赁时间
	defLeaseSecond = 30
)

// 保存已经注册的服务
var serviceMap map[string]*EndpointUnit

// 服务元数据信息
type Endpoint struct {
	Name    string `json:"name"` // 二级服务名称，比如（my_system）中的user，作为endpoint的末级标识
	Addr    string `json:"addr"`
	Port    int    `json:"port"`
	Version string `json:"version"`
}

// 服务元数据信息扩展，主要加上etcd的租约id
type EndpointUnit struct {
	Name    string
	LeaseID clientv3.LeaseID
	E       Endpoint
}

// 服务注册、发现和管理的核心结构体
type NamingService struct {
	EtcdUrl  string // etcd的url地址，如：http://localhost:2379
	EtcdHost string
	EtcdPort int
	Name     string            // 一级服务名称，比如 my_system
	MTarget  string            // nameServicePrefix/Name，它作为manager的标识
	Client   *clientv3.Client  // etcd 客户端，用于与 etcd 通信
	manager  endpoints.Manager // etcd 的 endpoints 管理器，用于处理服务的注册和注销
}

// 创建一个新的命名服务
//
// 这个注册器参考了(https://blog.csdn.net/small_to_large/article/details/130656230)，按本项目情况做了适配性改动
func NewNamingService(serviceName string) (*NamingService, error) {
	etcdClientWrapper := my_client.EtcdClientWrapper
	etcdClient := etcdClientWrapper.EtcdClient
	etcdUrl := etcdClientWrapper.EtcdUrl
	etcdHost := etcdClientWrapper.EtcdHost
	etcdPort := etcdClientWrapper.EtcdPort

	target := fmt.Sprintf("%s/%s", nameServicePrefix, serviceName)
	// etcd的endpoints管理
	manager, err := endpoints.NewManager(etcdClient, target)
	if err != nil {
		return nil, err
	}
	ns := NamingService{
		EtcdUrl:  etcdUrl,
		EtcdHost: etcdHost,
		EtcdPort: etcdPort,
		Name:     serviceName,
		MTarget:  target,
		manager:  manager,
		Client:   etcdClient,
	}
	serviceMap = make(map[string]*EndpointUnit)
	return &ns, nil
}

func (naming *NamingService) GetFullServerName(name string, leaseID clientv3.LeaseID) string {
	return fmt.Sprintf("%s/%s/%d", naming.MTarget, name, leaseID)
}

func (naming *NamingService) GetPathServerName(name string) string {
	return fmt.Sprintf("etcd://%s:%d/%s/%s", naming.EtcdHost, naming.EtcdPort, naming.MTarget, name)
}

// 添加/注册新的服务
func (naming *NamingService) AddEndpoint(e Endpoint) error {
	b, _ := json.Marshal(e)
	ep := endpoints.Endpoint{
		Addr:     fmt.Sprintf("%s:%d", e.Addr, e.Port),
		Metadata: string(b), // 这里有个坑，必须传入字符串，没来得及看原因
	}
	// 在etcd创建一个续期的lease对象
	lease, err := naming.Client.Grant(context.TODO(), defLeaseSecond)
	if err != nil {
		return err
	}
	key := naming.GetFullServerName(e.Name, lease.ID)
	// 向etcd注册一个Endpoint并绑定续期
	err = naming.manager.AddEndpoint(context.TODO(), key, ep, clientv3.WithLease(lease.ID))
	if err != nil {
		return err
	}
	log.Println("AddEndpoint success:", key)
	// 开启自动续期KeepAlive
	ch, err := naming.Client.KeepAlive(context.TODO(), lease.ID)
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

	// 用于本地记录服务注册信息，便于后续删除或更新。
	serviceMap[e.Name] = &EndpointUnit{Name: e.Name, LeaseID: lease.ID, E: e}
	return nil
}

// 移除一个服务
func (naming *NamingService) DelEndpoint(name string) error {
	eu := serviceMap[name]
	if eu == nil {
		return nil
	}

	// 从 etcd 中删除指定的 Endpoint
	err := naming.manager.DeleteEndpoint(context.TODO(), naming.GetFullServerName(name, eu.LeaseID))
	if err != nil {
		log.Fatalf("Delete endpoint error %v", err)
		return err
	}

	// 撤销租约，确保本客户端不会继续往 etcd 中续约
	_, err = naming.Client.Revoke(context.TODO(), eu.LeaseID)
	if err != nil {
		log.Fatalf("Revoke lease error %v", err)
		return err
	}

	// 删除本地记录
	delete(serviceMap, name)
	log.Printf("DeleteEndpoint [%s] success\n", name)
	return nil
}

// 移除所有服务
func (naming *NamingService) DelAllEndpoint() {
	for k := range serviceMap {
		err := naming.DelEndpoint(k)
		if err != nil {
			log.Fatalln("Ignore Failure Continue...")
		}
	}
}

// 创建一个 etcd 的 gRPC 解析器（resolver）。通过它，客户端能够自动发现和连接到注册的 gRPC 服务
func (naming *NamingService) NewEtcdResolver() (gresolver.Builder, error) {
	etcdResolver, err := resolver.NewBuilder(naming.Client)
	if err != nil {
		log.Fatalf("Etcd resolver error %v", err)
		return nil, err
	}
	return etcdResolver, nil
}
