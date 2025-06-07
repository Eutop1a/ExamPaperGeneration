CREATE TABLE `QuestionBank` (
    `id` int NOT NULL AUTO_INCREMENT,
    `topic` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
    `topic_material_id` int DEFAULT NULL,
    `answer` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
    `topic_type` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
    `score` decimal(10,1) DEFAULT NULL,
    `difficulty` int DEFAULT NULL,
    `chapter_1` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
    `chapter_2` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
    `label_1` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
    `label_2` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
    `update_time` datetime DEFAULT NULL,
    `topic_image_path` varchar(255) DEFAULT NULL, -- 存储图片路径
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1381 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;

CREATE TABLE `QuestionGenHistory` (
    `id` int NOT NULL AUTO_INCREMENT,
    `test_paper_uid` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
    `test_paper_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
    `question_bank_id` int DEFAULT NULL,
    `topic` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
    `topic_material_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
    `answer` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
    `topic_type` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
    `score` decimal(10,1) DEFAULT NULL,
    `difficulty` int DEFAULT NULL,
    `chapter_1` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
    `chapter_2` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
    `label_1` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
    `label_2` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
    `update_time` datetime DEFAULT NULL,
    `topic_image_path` varchar(255) DEFAULT NULL, -- 存储图片路径
    `topic_table_json` json DEFAULT NULL, -- 存储表格的 JSON 数据
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=147 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;

CREATE TABLE `QuestionLabels` (
  `id` int NOT NULL AUTO_INCREMENT,
  `chapter_1` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `chapter_2` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `label_1` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `label_2` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=51 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;

CREATE TABLE `QuestionMaterial` (
  `id` int NOT NULL AUTO_INCREMENT,
  `material` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
  `update_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;

CREATE TABLE `TestPaperGenHistory` (
  `id` int NOT NULL AUTO_INCREMENT,
  `test_paper_uid` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `test_paper_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `question_count` int DEFAULT NULL,
  `average_difficulty` decimal(10,2) DEFAULT NULL,
  `update_time` datetime DEFAULT NULL,
  `username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `similarity_threshold` decimal(10,2) DEFAULT 0.5, -- 用户特定的相似度阈值
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;

CREATE TABLE `User` (
  `id` int NOT NULL AUTO_INCREMENT,
  `username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `user_role` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `last_login` datetime(6) DEFAULT NULL,
  `enable` int DEFAULT NULL,
  `similarity_threshold` decimal(10,2) DEFAULT 0.5, -- 用户特定的相似度阈值
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_username` (`username`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;

INSERT INTO `QuestionLabels` (`chapter_1`, `chapter_2`, `label_1`, `label_2`) VALUES 
('1','1.1','绪论','微型计算机发展概况'),
('1','1.2','绪论','计算机中数和字符的表示'),
('1','1.3','绪论','微型计算机系统概论'),
('2','2.1','Intel8086微处理器','8086微处理器的内部结构'),
('2','2.2','Intel8086微处理器','8086引脚功能'),
('2','2.3','Intel8086微处理器','8086系统总线时序'),
('2','2.4','Intel8086微处理器','8086寻址方式'),
('2','2.5','Intel8086微处理器','8086指令系统'),
('3','3.1','宏汇编语言程序设计','汇编语言的语句格式'),
('3','3.2','宏汇编语言程序设计','汇编语言的数据项'),
('3','3.3','宏汇编语言程序设计','汇编语言的表达式'),
('3','3.4','宏汇编语言程序设计','伪指令语句'),
('3','3.5','宏汇编语言程序设计','汇编语言程序设计概述'),
('3','3.6','宏汇编语言程序设计','顺序程序设计'),
('3','3.7','宏汇编语言程序设计','分支程序设计'),
('3','3.8','宏汇编语言程序设计','循环程序设计'),
('3','3.9','宏汇编语言程序设计','DOS系统功能调用'),
('3','3.10','宏汇编语言程序设计','子程序设计'),
('3','3.11','宏汇编语言程序设计','宏指令'),
('3','3.12','宏汇编语言程序设计','汇编语言程序的建立、汇编、连接与调试'),
('4','4.1','Intel80486微处理器','80486内部结构'),
('4','4.2','Intel80486微处理器','80486的工作方式'),
('4','4.3','Intel80486微处理器','80486引脚功能'),
('4','4.4','Intel80486微处理器','80486的寻址方式'),
('4','4.5','Intel80486微处理器','80486常用指令介绍'),
('4','4.6','Intel80486微处理器','80486编程举例'),
('5','5.1','半导体存储器','存储器概述'),
('5','5.2','半导体存储器','随机存储器RAM'),
('5','5.3','半导体存储器','只读存储器ROM'),
('5','5.4','半导体存储器','存储器与CPU的连接'),
('5','5.5','半导体存储器','高速缓冲存储器系统'),
('6','6.1','I/O接口技术','I/O接口技术概述'),
('6','6.2','I/O接口技术','程序控制的I/O'),
('6','6.3','I/O接口技术','DMA方式'),
('7','7.1','中断系统','中断系统概述'),
('7','7.2','中断系统','16位微机中断系统'),
('7','7.3','中断系统','32位微处理器的中断'),
('7','7.4','中断系统','中断控制器8259A'),
('8','8.1','常用接口芯片','并行接口芯片8255A'),
('8','8.2','常用接口芯片','定时器/计数器接口芯片8253'),
('8','8.3','常用接口芯片','串行接口芯片8251A'),
('8','8.4','常用接口芯片','模拟接口'),
('8','8.5','常用接口芯片','多功能外围接口芯片82380'),
('9','9.1','总线','总线概述'),
('9','9.2','总线','ISA总线'),
('9','9.3','总线','EISA总线'),
('9','9.4','总线','PCI总线'),
('10','10.1','典型微型计算机系统','1BMPC/XT微型计算机系统'),
('10','10.2','典型微型计算机系统','80486微型计算机系统'),
('10','10.3','典型微型计算机系统','Pentium系列微型计算机系统');



