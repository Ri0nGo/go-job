-- 任务表
CREATE TABLE `job` (
    `id` int NOT NULL AUTO_INCREMENT,
    `name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '任务名称',
    `exec_type` smallint NOT NULL COMMENT '执行类型 1: shell; 2: http; 3:file',
    `cron_expr` varchar(128) DEFAULT NULL COMMENT 'cron 表达式',
    `node_id` int NOT NULL COMMENT '节点id',
    `internal` json DEFAULT NULL,
    `created_time` datetime DEFAULT NULL,
    `updated_time` datetime DEFAULT NULL,
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
