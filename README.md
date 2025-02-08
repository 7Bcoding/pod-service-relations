# 项目名称
一个获取异常状态pod及影响的services的巡检程序

## 快速开始
go mod tidy 安装部署所有依赖

启动server: go run main.go  // 目前不需要

建表：go run test/create_db_tables/main.go

初始化config配置：

bash env_init.sh online // 线上环境配置
bash env_init.sh test   // 测试环境配置

获取pod-service关系: 
go run cmd/get-pod-services/main.go

## 测试

测试单个文件：
go test -v -cover test.go

测试单个方法：
go test -v -cover test.go -run TestFunc

测试所有正则匹配TestFunc的方法：
go test -v -cover test.go -run "TestFunc"

## 如何贡献
贡献patch流程、质量要求

## 讨论


## 链接

