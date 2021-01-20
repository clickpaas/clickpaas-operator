# clickpaas-operator

# 部署controller到kubernetes
1. 下载源代码构建镜像
```bash
git clone https://github.com/clickpaas/clickpaas-operator.git
make all
```


2. 开始部署中间件
> middle 基本yml文件位于artifacts/middleware 目录下
```bash
# 安装mysql
kubectl apply -f artifacts/middleware/mysql
# 安装redis
kubectl apply -f artifacts/middleware/gcache
```


3. 自动化部分
> diamond 部署自动往mysql里面注册表/数据信息 <br>
> lts 自动在mysql里面创建数据库 <br>
> redis/gcache 自动完成集群初始化(不依赖其他配置)

4. issue
> zookeeper 批量创建数据  --- 暂时未提供 <br>
> mysql 数据初始化未提供 <br>
> mongo 数据自动初始化未提供 <br>
> redis 数据导入未提供 <br>
> 样例中间件编排文件里面服务名称还未适配ipaas/apaas