-- 任务表
CREATE TABLE `job` (
    `id` int NOT NULL AUTO_INCREMENT,
    `name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '任务名称',
    `exec_type` smallint NOT NULL COMMENT '执行类型 1: shell; 2: http; 3:file',
    `active` smallint DEFAULT '1' COMMENT '启用状态 1启用；2停用',
    `cron_expr` varchar(128) DEFAULT NULL COMMENT 'cron 表达式',
    `node_id` int NOT NULL COMMENT '节点id',
    `notify_status` tinyint(1) DEFAULT '2' COMMENT '通知状态 1启用；2停用',
    `notify_type` smallint DEFAULT NULL COMMENT '通知方式 1邮件',
    `notify_strategy` smallint DEFAULT NULL COMMENT '通知策略 1成功后通知；2失败后通知；3总是通知',
    `internal` json DEFAULT NULL,
    `created_time` datetime DEFAULT NULL,
    `updated_time` datetime DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 任务记录表
CREATE TABLE `job_record` (
    `id` int NOT NULL AUTO_INCREMENT,
    `job_id` int NOT NULL,
    `status` smallint COMMENT '执行状态 0待执行；1运行中；2成功；3失败',
    `start_time` datetime DEFAULT NULL COMMENT '开始执行时间',
    `end_time` datetime DEFAULT NULL COMMENT '结束执行时间',
    `duration` float DEFAULT NULL COMMENT '运行耗时',
    `output` text COMMENT '执行文件内容输出',
    `error` text COMMENT '节点执行异常日志',
    `next_exec_time` datetime DEFAULT NULL COMMENT '任务下一次执行时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 标签表
CREATE TABLE `tag` (
    `id` int NOT NULL AUTO_INCREMENT,
    `name` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
    `type` smallint NOT NULL COMMENT '1: 任务模块的标签',
    `description` varchar(200) DEFAULT NULL,
    `created_time` datetime DEFAULT NULL,
    `updated_time` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_tag_name_type` (`name`,`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 标签和模块关系表
CREATE TABLE `tag_ref` (
    `id` int NOT NULL AUTO_INCREMENT,
    `tag_id` int NOT NULL,
    `ref_id` int NOT NULL,
    `ref_type` smallint NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_tag_id_ref_id_type` (`tag_id`,`ref_id`,`ref_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 节点表
CREATE TABLE `node` (
    `id` int NOT NULL AUTO_INCREMENT,
    `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '节点名称',
    `description` varchar(200) DEFAULT NULL COMMENT '节点描述',
    `address` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '节点地址，address=ip:port',
    `created_time` datetime DEFAULT NULL,
    `updated_time` datetime DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 用户表
CREATE TABLE `user` (
    `id` int NOT NULL AUTO_INCREMENT,
    `username` varchar(48) NOT NULL,
    `password` varchar(128) NOT NULL,
    `nickname` varchar(64) DEFAULT NULL,
    `email` varchar(255) DEFAULT NULL,
    `about` varchar(200) DEFAULT NULL,
    `created_time` datetime DEFAULT NULL,
    `updated_time` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_uniq_username` (`username`),
    UNIQUE KEY `idx_uniq_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;