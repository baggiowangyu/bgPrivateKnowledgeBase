DROP TABLE IF EXISTS `bg_device_info`;
CREATE TABLE `bg_device_info` (
  `auto_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '内建唯一ID，自增',
  `internal_id` varchar(64) NOT NULL COMMENT '设备内部ID',
  `gbcode_id` varchar(64) NOT NULL COMMENT '设备国标ID',
  `name` varchar(64) NOT NULL COMMENT '设备名称',
  `type` int NOT NULL COMMENT '设备类型：0-无网络设备；1-有线网络设备；2-无线网络设备',
  `reg_time` datetime NOT NULL COMMENT '注册时间',
  `update_time` datetime NOT NULL COMMENT '更新时间',
  `from` varchar(64) NOT NULL COMMENT '数据来源',
  `user_id` varchar(64) COMMENT '使用人ID',
  `user_name` varchar(64) COMMENT '使用人姓名',
  `reserve01` varchar(64) COMMENT '预留字段01',
  `reserve02` varchar(64) COMMENT '预留字段02',
  `reserve03` varchar(64) COMMENT '预留字段03',
  `reserve04` varchar(64) COMMENT '预留字段04',
  `reserve05` varchar(64) COMMENT '预留字段05',
  `reserve06` varchar(64) COMMENT '预留字段06',
  `reserve07` varchar(64) COMMENT '预留字段07',
  `reserve08` varchar(64) COMMENT '预留字段08',
  `reserve09` varchar(64) COMMENT '预留字段09',
  `reserve10` varchar(64) COMMENT '预留字段10',
  PRIMARY KEY (`auto_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;