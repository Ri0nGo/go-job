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