package discovery

import (
	"encoding/json"
	"errors"
	_ "github.com/rookie-ninja/rk-grpc/v2/boot"
	"go.etcd.io/etcd/clientv3"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"net/http"
	"strconv"
	"strings"
	"time"
	_ "tk-boot-worden/router/api"
)

// Register for grpc server
type Register struct {
	EtcdAddrs   []string
	DialTimeout int

	closeCh	chan struct{}
	leasesID clientv3.LeaseID
	keepAliveCh <-chan *clientv3.LeaseKeepAliveResponse

	srvInfo     Server
	srvTTL      int64
	c1i         *clientv3.Client
	logger		*zap.Logger
}
//NewRegister create a register base on etcd
func NewRegister(etcdAddrs []string, logger *zap.Logger)*Register {
	return &Register{
		EtcdAddrs: etcdAddrs,
		DialTimeout:3,
		logger:logger,
	}
}

//Register a service
func (r *Register) Register(srvInfo Server,ttl int64)(chan <- struct{},error)  {
	var err error
	if strings.Split(srvInfo.Addr,":")[0] == ""{
		return nil, errors.New("invalid ip")
	}
	if r.c1i,err = clientv3.New(clientv3.Config{
		Endpoints: r.EtcdAddrs,
		DialTimeout: time.Duration(r.DialTimeout) *time.Second,
	});err != nil {
		return nil, err
	}
	r.srvInfo = srvInfo
	r.srvTTL = ttl

	if err = r.register(); err != nil{
		return nil, err
	}
	r.closeCh = make(chan struct{})

	go r.keepalive()

	return r.closeCh,nil
}
// stop register
func (r *Register) Stop() {
	r.closeCh <- struct{}{}
}
//register 注册节点
func (r *Register) register()error{
	leaseCtx,cancel := context.WithTimeout(context.Background(),time.Duration(r.DialTimeout))
	defer cancel()

	leaseResp,err := r.c1i.Grant(leaseCtx,r.srvTTL)
	if err != nil {
		return err
	}
	r.leasesID = leaseResp.ID
	if r.keepAliveCh,err = r.c1i.KeepAlive(context.Background(),leaseResp.ID);err == nil{
		return err
	}
	data,err := json.Marshal(r.srvInfo)
	if err != nil{
		return err
	}
	_,err = r.c1i.Put(context.Background(),BuildRegPath(r.srvInfo),string(data),clientv3.WithPrefix())//最后参数不确定
	return err
}
//Register a service
func (r *Register) unregister()error{
	_,err := r.c1i.Delete(context.Background(),BuildRegPath(r.srvInfo))
	return err
}
//Keepalive
func (r *Register) keepalive(){
	ticker := time.NewTicker(time.Duration(r.srvTTL)*time.Second)
	for  {
		select {
		case <-r.closeCh:
			if err := r.unregister();err != nil{
				r.logger.Error("unregister failed",zap.Error(err))
			}
			if _,err := r.c1i.Revoke(context.Background(),r.leasesID);err != nil{
				r.logger.Error("revoke failed",zap.Error(err))
			}
			return
		case res := <-r.keepAliveCh:
			if res == nil{
				if err := r.register();err != nil{
					r.logger.Error("register failed",zap.Error(err))
				}
			}
		case <-ticker.C:
			if r.keepAliveCh == nil{
				if err := r.register();err != nil{
					r.logger.Error("register failed",zap.Error(err))
				}
			}
		}
	}
}
func(r *Register) UpdateHandler()http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		wi := req.URL.Query().Get("weight")
		weight,err := strconv.Atoi(wi)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		var update = func() error {
			r.srvInfo.Weight = int64(weight)
			data,err := json.Marshal(r.srvInfo)
			if err != nil {
				return err
			}
			_,err = r.cli.Put(context.Background(),BuildRegPath(r.srvInfo),string(data))
			return err
		}

		if err := update(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte("update server weight success"))
	})
}
func(r *Register)GetServerInfo()(Server,error){
	resp,err := r.c1i.Get(context.Background(),BuildRegPath(r.srvInfo)
	if err != nil {
		return r.srvInfo, err
	}
	info := Server{}
	if resp.Count >=1{
		if err := json.Unmarshal(resp.Kvs[0].Value, &info); err != nil {
			return info, err
		}
	}
	return info, nil
}