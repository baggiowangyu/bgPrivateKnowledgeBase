-- 路由信息
DROP TABLE IF EXISTS `api_router_info`;
CREATE TABLE `api_router_info` (
    `Id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `Protocol` VARCHAR(255) NOT NULL,       -- 路由协议，目前只支持HTTP
    `Gateway_api` VARCHAR(255) NOT NULL,    -- 网关对外开放的API
    `Service_api` VARCHAR(255) NOT NULL     -- 实际服务提供的API
);

-- 接口提供者信息
DROP TABLE IF EXISTS `api_service_info`;
CREATE TABLE `api_service_info`
(
    `Id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `Gateway_api` VARCHAR(255) NOT NULL,    -- 网关对外开放的API，
    `Service_address` VARCHAR(255) NOT NULL -- 实际服务的地址
);

-- 接口访问统计信息
DROP TABLE IF EXISTS `api_access_statistics`;
CREATE TABLE `api_access_statistics` (
    `Id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `Gateway_api` VARCHAR(255) NOT NULL,    -- 网关对外开放的API，
    `access_count` BIGINT NOT NULL
);

-- 接口访问日志
DROP TABLE IF EXISTS `api_access_log`;
CREATE TABLE `api_access_log`
(
    `Id`          INTEGER PRIMARY KEY AUTOINCREMENT,
    `Gateway_api` VARCHAR(255) NOT NULL, -- 网关对外开放的API，
    `Accesser_address` VARCHAR(255) NOT NULL,
    `Access_time` DATETIME NOT NULL
);