## 2025-04-18 user表新增邮箱字段

```mysql
ALTER TABLE `go-job`.`user` 
ADD COLUMN `email` varchar(255) NULL,
ADD UNIQUE INDEX `idx_uniq_email`(`email`);
```


## 2025-04-22 job表新增通知

```mysql
ALTER TABLE `go-job`.`job`
    ADD COLUMN `notify_status` tinyint(1) NULL DEFAULT 2 COMMENT '通知状态 1启用；2停用',
    ADD COLUMN `notify_type` smallint NULL COMMENT '通知方式 1邮件',
    ADD COLUMN `notify_strategy` smallint NULL COMMENT '通知策略 1成功后通知；2失败后通知；3总是通知';
```

## 2025-04-23 job表新增用户id

```mysql
ALTER TABLE `go-job`.`job` 
ADD COLUMN `user_id` int NOT NULL COMMENT '用户id',
ADD COLUMN `notify_mark` varchar(255) NULL COMMENT '通知方式的具体内存，如邮箱地址';
```

## 2025-05-02 删除job表中的通知字段

```mysql
#     `notify_status` tinyint(1) DEFAULT '2' COMMENT '通知状态 1启用；2停用',
#     `notify_type` smallint DEFAULT NULL COMMENT '通知方式 1邮件',
#     `notify_strategy` smallint DEFAULT NULL COMMENT '通知策略 1成功后通知；2失败后通知；3总是通知',
#     `notify_mark` varchar(255) DEFAULT NULL COMMENT '通知方式的具体内存，如邮箱地址',
ALTER TABLE job_record
    DROP COLUMN notify_status,
    DROP COLUMN notify_type,
    DROP COLUMN notify_strategy,
    DROP COLUMN notify_mark;
```

## 2025-05-12 job_record 表新增索引

```mysql
CREATE INDEX idx_job_id ON job_record(job_id);
CREATE INDEX idx_status ON job_record(status);
```


## 2025-06-08 新增身份认证表

```sql
CREATE TABLE `auth_identity` (
     `id` int NOT NULL AUTO_INCREMENT,
     `user_id` int NOT NULL,
     `type` smallint NOT NULL COMMENT '授权类型, 1:github',
     `identity` varchar(128) NOT NULL COMMENT '身份标识',
     `name` varchar(128) DEFAULT NULL COMMENT '授权平台的用户名',
     `created_time` datetime DEFAULT NULL,
     `updated_time` datetime DEFAULT NULL,
     PRIMARY KEY (`id`),
     KEY `idx_type_identity` (`type`,`identity`) COMMENT '类型和身份标识唯一'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
```

## 2025-06-18 用户表新增登录时间字段

```mysql
alter table `user`
    add login_time datetime null;
```