# 概述

该项目来自于一个 GitHub 的开源项目，使用了 Go 语言重写了后端。

项目链接：https://github.com/inferno0303/TestPapaerGen-WebApp

如果上述链接失效，使用这个：https://github.com/Eutop1a/TestPapaerGen-WebApp

## 介绍

名称：在线组卷系统 - ExamPapaerGeneration

简介：自动组卷系统，遗传算法、贪心算法，支持导入题库，手动选择、自动组卷，生成排版美观的Word文档，前后端分离WebApp，

技术栈：后端 Go + Gin，前端 React Umi.js

类型：WebApp



## 安装

### 目录结构

component：拦截器相关

config：配置文件

controller：控制器，包含接口处理函数

entity：实体，包含数据库中实体的定义

mapper：包含各种实体的操作函数，其中指定数据库连接信息的文件为 repository.go

resource：资源文件，包括模板文件，数据库文件，试卷模板，试题

services：服务，接口处理函数调用这里的服务

templates：试卷模板文件

test：测试相关

utils：通用工具

### 如何运行



后端：项目基于 go1.24 开发，首先需要本地安装 go 语言，[Go语言中文网](https://studygolang.com/dl)，拉取依赖成功后运行

```
go mod tidy 
go run main.go
```

数据库：记得导入数据库表结构，默认utf8mb4，数据库表结构sql文件已包含建库、建表语句。

```
mysql -u root -h host -p < xxx.sql
```



## 功能



1. 登录功能，支持注册账号，登录，基于拦截器实现的权限认证；
2. 题库管理，支持填空题、选择题、判断题、设计题、阅读题多种类型，所有题目可自由增删改查；
3. 手动组卷，支持手动从试题库选择题目，加入组卷列表中，作为题目输出；
4. 自动组卷，支持按照难易度，题目类型数，分值，章节等多个维度按需自动组卷，带随机算法，并非简单查库，可按照相同设置自动出A/B卷。
5. 输出试卷，排版美观，输出openxml格式的文档，可office打开，可以直接打印，效果如下图；
6. 出题历史，如字面意思，可查看出卷历史，统计出卷难度，复盘试卷题型；
7. 完善的可视化统计，各种炫酷的图表，可视化汇报数据状态，基于Echarts。
