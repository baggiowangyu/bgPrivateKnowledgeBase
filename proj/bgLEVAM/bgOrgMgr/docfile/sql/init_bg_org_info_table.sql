DROP TABLE IF EXISTS `bg_org_info`;
CREATE TABLE `bg_org_info` (
  #`auto_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '内建唯一ID，自增',
  `internal_id` varchar(64) NOT NULL COMMENT '部门内部ID，公安系统可以是警综ID',
  `parent_internal_id` varchar(64) NOT NULL COMMENT '父部门内部ID，公安系统可以是警综ID',
  `gbcode_id` varchar(64) COMMENT '部门国标ID',
  `name` varchar(64) NOT NULL COMMENT '部门名称',
  `path` varchar(128) COMMENT '部门路径',
  `type` int NOT NULL COMMENT '部门类型：0-职能部门；1-业务部门；2-派出部门',
  `reg_time` datetime NOT NULL COMMENT '注册时间',
  `update_time` datetime NOT NULL COMMENT '更新时间',
  `from` varchar(64) NOT NULL COMMENT '数据来源',
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
  PRIMARY KEY (`internal_id`)
  #PRIMARY KEY (`auto_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
