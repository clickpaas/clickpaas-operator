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