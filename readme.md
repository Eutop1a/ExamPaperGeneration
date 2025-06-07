# 概述

本项目基于 GitHub 的开源项目 TestPaperGen-WebAPP，使用 Go 语言重写和二次开发。

项目链接：https://github.com/inferno0303/TestPapaerGen-WebApp

如果上述链接失效，使用这个：https://github.com/Eutop1a/TestPapaerGen-WebApp



## 介绍

名称：在线组卷系统 - ExamPapaerGeneration

简介：自动组卷系统，遗传算法、贪心算法，支持导入题库，手动选择、自动组卷，生成排版美观的Word文档，前后端分离WebApp，

技术栈：后端 Go + Gin，前端 React Umi.js

类型：WebApp



## 安装

### 目录结构

#### ExamPaperGenerationBe

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

#### ExamPaperGenerationFe

src：源代码



### 如何运行

后端：项目基于 go1.24 开发，首先需要本地安装 go 语言，拉取依赖成功后运行

```
go mod tidy 
go run main.go
```

数据库：记得导入数据库表结构，默认utf8mb4，数据库表结构sql文件已包含建库、建表语句。

tabel.sql 在 ExamPaperGenerationBe 的 resource 文件夹下

```
mysql -u root -h host -p < tabel.sql
```

前端：标准 webpack 工程，在 package.json 目录下执行 npm install 拉取依赖，npm start 运行工程，npm build 构建工程。

需要注意，该项目依赖的 nodejs 版本较低，建议使用 nodejs v16.16.0

```
npm install
npm start
npm build
```

## 功能



1. 登录功能，支持注册账号，登录，基于拦截器实现的权限认证；
2. 题库管理，支持填空题、选择题、判断题、简答题多种类型，所有题目可自由增删改查；
3. 手动组卷，支持手动从试题库选择题目，加入组卷列表中，作为题目输出；
4. 自动组卷，支持按照难易度，题目类型数，分值，章节等多个维度按需自动组卷，带随机算法，并非简单查库，可按照相同设置自动出A/B卷。
5. 输出试卷，排版美观，输出openxml格式的文档，可office打开，可以直接打印，效果如下图；
6. 出题历史，如字面意思，可查看出卷历史，统计出卷难度，复盘试卷题型；
7. 完善的可视化统计，各种炫酷的图表，可视化汇报数据状态，基于Echarts。



## 运行



### 视频演示

视频演示：https://www.bilibili.com/video/BV1mcTLzpETo

### 欢迎首页

![image-20250607161606885](https://renovice-1311449499.cos.ap-chongqing.myqcloud.com/image-20250607161606885.png)

### 题库管理

![image-20250607161626336](https://renovice-1311449499.cos.ap-chongqing.myqcloud.com/image-20250607161626336.png)

新增了对图片的支持

![image-20250607161658445](https://renovice-1311449499.cos.ap-chongqing.myqcloud.com/image-20250607161658445.png)

添加题目处，新增了导入图片的接口：

![image-20250607161807678](https://renovice-1311449499.cos.ap-chongqing.myqcloud.com/image-20250607161807678.png)

### 知识点标签管理

支持对以往的知识点进行更新和删除，同时可以添加新的知识点

![image-20250607162921468](https://renovice-1311449499.cos.ap-chongqing.myqcloud.com/image-20250607162921468.png)

### 导出试卷文档

- 支持导出word格式文档
- 支持导出参考答案

![image-20250607162347845](https://renovice-1311449499.cos.ap-chongqing.myqcloud.com/image-20250607162347845.png)

### 导入题库

- 支持导入Excel格式的题库

![image-20250607162411080](https://renovice-1311449499.cos.ap-chongqing.myqcloud.com/image-20250607162411080.png)

excel 的格式如ExamPaperGenerationBe\resources\QuestionBanl.xlsx所示

![image-20250607162552551](https://renovice-1311449499.cos.ap-chongqing.myqcloud.com/image-20250607162552551.png)

### 题库概览

![image-20250607161831163](https://renovice-1311449499.cos.ap-chongqing.myqcloud.com/image-20250607161831163.png)

### 自动组卷

组卷部分支持根据知识点权重出题：

![image-20250607161944041](https://renovice-1311449499.cos.ap-chongqing.myqcloud.com/image-20250607161944041.png)

使用贪心算法出题：

![image-20250607162046780](https://renovice-1311449499.cos.ap-chongqing.myqcloud.com/image-20250607162046780.png)

![image-20250607162059522](https://renovice-1311449499.cos.ap-chongqing.myqcloud.com/image-20250607162059522.png)



使用遗传算法出题：

![image-20250607162121373](https://renovice-1311449499.cos.ap-chongqing.myqcloud.com/image-20250607162121373.png)

![image-20250607162127750](https://renovice-1311449499.cos.ap-chongqing.myqcloud.com/image-20250607162127750.png)

### 出题历史

![image-20250607162611129](https://renovice-1311449499.cos.ap-chongqing.myqcloud.com/image-20250607162611129.png)

### 重新编辑历史组卷

![image-20250607162632415](https://renovice-1311449499.cos.ap-chongqing.myqcloud.com/image-20250607162632415.png)

### 注册账户、管理员账户

![image-20250607162647170](https://renovice-1311449499.cos.ap-chongqing.myqcloud.com/image-20250607162647170.png)
