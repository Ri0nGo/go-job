## 20250-04-18 user表新增邮箱字段

```mysql
ALTER TABLE `go-job`.`user` 
ADD COLUMN `email` varchar(255) NULL AFTER `nickname`,
ADD UNIQUE INDEX `idx_uniq_email`(`email`);
```