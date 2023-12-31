# meta_launchpad

元初发射台，SAAS运营活动、包括普通购，优先购，秒杀，空投

```
# HTTP Server
[server]
    Address     = ":8126"
[nacos]
    # 应用标识
    AppName                 = "launchpad"
    # 运行模式
    Mode                    = "dev"
    # 配置文件类型
    FileExtension           = "toml"
    # 是否启用nacos作为配置中心
    EnableConfig            = true
    # 是否启用nacos用为服务发现注册中心
    EnableDiscovery         = true
    # 应用的监听端口
    AppPort                 = 8126
    # 配置文件的Group属性，默认为public
    ConfigGroup             = "public"
    # 服务注册中的Group属性，默认为DEFAULT_GROUP
    DiscoveryGroup          = "DEFAULT_GROUP"
    # 服务器的IP，在gateway中将通过该ip对外提供服务
    AppIp                   = "39.107.72.102"

# 服务注册时的meta参数，全部自定义
[nacos.meta]
    private                 = "false"

# nacos配置中心的配置参数，配置来自：https://github.com/nacos-group/nacos-sdk-go/blob/master/common/constant/config.go的ClientConfig
[nacos.config]
    TimeoutMs = 5000
    NotLoadCacheAtStart = true
    RotateTime = "1h"
    MaxAge = 3
    LogLevel = "debug"
    #TimeoutMs            uint64 //timeout for requesting Nacos server, default value is 10000ms
    #ListenInterval       uint64 //Deprecated
    #BeatInterval         int64  //the time interval for sending beat to server,default value is 5000ms
    #NamespaceId          string //the namespaceId of Nacos
    #Endpoint             string //the endpoint for get Nacos server addresses
    #RegionId             string //the regionId for kms
    #AccessKey            string //the AccessKey for kms
    #SecretKey            string //the SecretKey for kms
    #OpenKMS              bool   //it's to open kms,default is false. https://help.aliyun.com/product/28933.html
    #CacheDir             string //the directory for persist nacos service info,default value is current path
    #UpdateThreadNum      int    //the number of gorutine for update nacos service info,default value is 20
    #NotLoadCacheAtStart  bool   //not to load persistent nacos service info in CacheDir at start time
    #UpdateCacheWhenEmpty bool   //update cache when get empty service instance from server
    #Username             string //the username for nacos auth
    #Password             string //the password for nacos auth
    #LogDir               string //the directory for log, default is current path
    #RotateTime           string //the rotate time for log, eg: 30m, 1h, 24h, default is 24h
    #MaxAge               int64  //the max age of a log file, default value is 3
    #LogLevel             string //the level of log, it's must be debug,info,warn,error, default value is info
    #ContextPath          string //the nacos server contextpath

# nacos服务发现配置参数，来自：https://github.com/nacos-group/nacos-sdk-go/blob/master/common/constant/config.go的ServerConfig
[nacos.discovery]
    Scheme = "http"
    ContextPath = "/nacos"
    IpAddr = "39.107.72.102"
    #IpAddr = "nacos.gfanx.com"
    # IpAddr = "nacos.gfanx.com"
    Port  = 8848
    #Scheme      string //the nacos server scheme
    #ContextPath string //the nacos server contextpath
    #IpAddr      string //the nacos server address
    #Port        uint64 //the nacos server port
[database]
    [[database.default]]
        link = "mysql:root:Gfanx#root#2022@tcp(39.107.72.102:33061)/meta_launchpad?charset=utf8mb4&parseTime=True&loc=Local"
        charset   = "utf8mb4" #数据库编码
        debug = true
        dryRun = false #空跑
        maxIdle      = "10" #连接池最大闲置的连接数
        maxOpen     = "10" #连接池最大打开的连接数
        maxLifetime  = "30" #(单位秒)连接对象可重复使用的时间长度

[redis]
    open = true #是否开启 redis 缓存 若不开启使用gchache缓存方式
    #default = "127.0.0.1:6379,9?idleTimeout=20&maxActive=100"
    default = "39.107.72.102:36379,0,redis_123456?idleTimeout=20&maxActive=100"

#开发者后台的配置
[developer]
    host = "http://39.107.72.102:1573"

# rabbitmq 配置
[rabbitmq]
    [rabbitmq.default]
        link = "amqp://chat:ZHlLORFUeUcv@39.107.72.102:5672/"
        

#sms配置
[sms]
    url = "https://market.juncdt.com/smartmarket/msgService/sendMessage"
    accessKey = "d0ad7a519058866cb5"
    accessSecret = "f27ee1886f933516ffb01d4bb6407f"
    classificationSecret = "2atbqWLP"
    signCode = "jHq74VEp"
    templateCode = "AqY0bsaM"          
```