-- 创建数据库
CREATE DATABASE IF NOT EXISTS aifriends_db
CHARACTER SET utf8mb4
COLLATE utf8mb4_unicode_ci;

-- 创建用户
CREATE USER 'aifriends'@'%' IDENTIFIED BY '123456';

-- 分配权限
GRANT ALL PRIVILEGES ON aifriends_db.* TO 'aifriends'@'%';

FLUSH PRIVILEGES;