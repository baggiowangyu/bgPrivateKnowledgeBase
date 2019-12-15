SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

# 首先是采集站信息
DROP TABLE IF EXISTS `business_base_workstation_info`;
CREATE TABLE `business_base_workstation_info` (
    `auto_id` int NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT '自增主键',
    `device_id` varchar(40) NOT NULL COMMENT '一般来讲应该为国标ID，我们这里预留40个字符，以便进行其他扩展',
    `deivce_name` varchar(128) COMMENT '设备名称',
    `department_id` varchar(40) COMMENT '所属部门ID，若没有值则为未使用',
    `department_name` varchar(128) COMMENT '所属部门名称',
    `manager_id` varchar(40) COMMENT '管理者ID',
    `manager_name` varchar(128) COMMENT '管理者姓名',
    `manager_phone_number` varchar(40) COMMENT '管理者联系电话',
    `ip_address` varchar(16) COMMENT '设备IP地址',
    `device_latitude` double COMMENT '设备所在位置纬度坐标',
    `device_longtitude` double COMMENT '设备所在位置经度坐标',
    `device_height` double COMMENT '设备所在位置海拔高度',
    `device_state` int NOT NULL COMMENT '设备状态',
    `detail` varchar(255) COMMENT '备注信息'
);

DROP TABLE IF EXISTS `business_base_collector_device_type`;
CREATE TABLE `business_base_collector_device_type`(
    `auto_id` bigint(20) NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT '自增主键',
    `device_type` int NOT NULL COMMENT '设备类型ID',
    `device_type_name` varchar(128) NOT NULL COMMENT '设备类型名称'
);

INSERT INTO `business_base_collector_device_type`(`device_type`, `device_type_name`)
VALUES (128, '采集工作站');
INSERT INTO `business_base_collector_device_type`(`device_type`, `device_type_name`)
VALUES (131, '安防监控探头');
INSERT INTO `business_base_collector_device_type`(`device_type`, `device_type_name`)
VALUES (150, '普通执法记录仪');
INSERT INTO `business_base_collector_device_type`(`device_type`, `device_type_name`)
VALUES (151, '4G执法记录仪');

# 执法采集设备信息(执法仪、或其他设备)
DROP TABLE IF EXISTS `business_base_collector_device_info`;
CREATE TABLE `business_base_collector_device_info`(
    `auto_id` bigint(20) NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT '自增主键',
    `device_id` varchar(40) NOT NULL COMMENT '一般来讲应该为国标ID，我们这里预留40个字符，以便进行其他扩展',
    `device_type` int NOT NULL COMMENT '设备类型，具体参见设备类型字典表',
    `deivce_name` varchar(128) COMMENT '设备名称',
    `department_id` varchar(40) COMMENT '所属部门ID，若没有值则为未使用',
    `department_name` varchar(128) COMMENT '所属部门名称',
    `wearer_id` varchar(40) COMMENT '设备佩戴者，或者说使用人',
    `wearer_name` varchar(128) COMMENT '设备佩戴者姓名',
    `device_state` int NOT NULL COMMENT '设备状态',
    `detail` varchar(255) COMMENT '备注信息'
);

# 存储服务器信息
DROP TABLE IF EXISTS `business_base_storage_server_info`;
CREATE TABLE `business_base_storage_server_info` (
    `auto_id` bigint(20) NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT '自增主键',
    `device_id` varchar(40) NOT NULL COMMENT '一般来讲应该为国标ID，我们这里预留40个字符，以便进行其他扩展',
    `deivce_name` varchar(128) COMMENT '设备名称',
    `department_id` varchar(40) COMMENT '所属部门ID，若没有值则为未使用',
    `department_name` varchar(128) COMMENT '所属部门名称',
    `manager_id` varchar(40) COMMENT '管理者ID',
    `manager_name` varchar(128) COMMENT '管理者姓名',
    `manager_phone_number` varchar(40) COMMENT '管理者联系电话',
    `ip_address` varchar(16) COMMENT '设备IP地址',
    `device_state` int NOT NULL COMMENT '设备状态',
    `detail` varchar(255) COMMENT '备注信息'
);