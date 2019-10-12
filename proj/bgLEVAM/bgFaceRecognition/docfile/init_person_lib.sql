DROP TABLE IF EXISTS `uap_tb_gmvcs_key_person_type_db`;
CREATE TABLE `uap_tb_gmvcs_key_person_type_db` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `db_name` varchar(50) NOT NULL COMMENT '重点人员库名称',
  `db_desc` varchar(255) DEFAULT NULL COMMENT '重点人员库描述信息',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime NOT NULL COMMENT '更新时间',
  `enable` bit(1) NOT NULL COMMENT '是否启用',
  `is_deleted` bit(1) NOT NULL COMMENT '删除标记',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8 COMMENT='重点人员库（用于人证核验）';

INSERT
INTO `uap_tb_gmvcs_key_person_type_db`(`db_name`, `db_desc`, `create_time`, `update_time`, `enable`, `is_deleted`)
VALUES ('重点人口库-0', '重点人口库-0', '2019-10-12 00:00:00', '2019-10-12 00:00:00', TRUE, FALSE);

INSERT
INTO `uap_tb_gmvcs_key_person_type_db`(`db_name`, `db_desc`, `create_time`, `update_time`, `enable`, `is_deleted`)
VALUES ('重点人口库-1', '重点人口库-1', '2019-10-12 00:00:00', '2019-10-12 00:00:00', TRUE, FALSE);

INSERT
INTO `uap_tb_gmvcs_key_person_type_db`(`db_name`, `db_desc`, `create_time`, `update_time`, `enable`, `is_deleted`)
VALUES ('重点人口库-2', '重点人口库-2', '2019-10-12 00:00:00', '2019-10-12 00:00:00', TRUE, FALSE);


DROP TABLE IF EXISTS `uap_tb_gmvcs_key_person_recognition`;
CREATE TABLE `uap_tb_gmvcs_key_person_recognition` (
  `person_id` varchar(255) NOT NULL COMMENT '人员唯一ID',
  `key_person_type_db_id` bigint(20) NOT NULL COMMENT '人员所属',
  `id_card` varchar(20) NOT NULL COMMENT '身份证号',
  `person_name` varchar(20) NOT NULL COMMENT '姓名',
  `birth_date` date DEFAULT NULL COMMENT '生日',
  `key_person_detail` varchar(255) DEFAULT NULL COMMENT '重点人员描述',
  `person_district` varchar(50) DEFAULT NULL COMMENT '区域',
  `person_address` varchar(100) DEFAULT NULL COMMENT '地址',
  `person_reg_img` varchar(255) DEFAULT NULL COMMENT '重点人员注册头像存储路径',
  `remarks` text COMMENT '备注',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime NOT NULL COMMENT '更新时间',
  `is_deleted` bit(1) NOT NULL COMMENT '删除标记',
  PRIMARY KEY (`person_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='重点人员信息（用于人证核验）';

INSERT
INTO `uap_tb_gmvcs_key_person_recognition`(
     `person_id`, `key_person_type_db_id`, `id_card`, `person_name`, `birth_date`, `key_person_detail`,
     `person_district`, `person_address`, `person_reg_img`, `create_time`, `update_time`, `is_deleted`)
VALUES ('440000198001010001', 0, '440000198001010001', '欧礼斌', '19800101', '犯了点事', '广东省广州市', '黄埔区开创大道2819号',
        '/1.jpg', '2019-10-12 00:00:00', '2019-10-12 00:00:00', FALSE);

INSERT
INTO `uap_tb_gmvcs_key_person_recognition`(
  `person_id`, `key_person_type_db_id`, `id_card`, `person_name`, `birth_date`, `key_person_detail`,
  `person_district`, `person_address`, `person_reg_img`, `create_time`, `update_time`, `is_deleted`)
VALUES ('440000198001020002', 0, '440000198001020002', '位旭彤', '19800102', '加班太狠', '广东省广州市', '黄埔区开创大道2819号',
        '/2.jpg', '2019-10-12 00:00:00', '2019-10-12 00:00:00', FALSE);

INSERT
INTO `uap_tb_gmvcs_key_person_recognition`(
  `person_id`, `key_person_type_db_id`, `id_card`, `person_name`, `birth_date`, `key_person_detail`,
  `person_district`, `person_address`, `person_reg_img`, `create_time`, `update_time`, `is_deleted`)
VALUES ('440000198001030003', 0, '440000198001030003', '方彦波', '19800103', '不请吃饭', '广东省广州市', '黄埔区开创大道2819号',
        '/3.jpg', '2019-10-12 00:00:00', '2019-10-12 00:00:00', FALSE);

INSERT
INTO `uap_tb_gmvcs_key_person_recognition`(
  `person_id`, `key_person_type_db_id`, `id_card`, `person_name`, `birth_date`, `key_person_detail`,
  `person_district`, `person_address`, `person_reg_img`, `create_time`, `update_time`, `is_deleted`)
VALUES ('440000198001040004', 1, '440000198001040004', '邱业兴', '19800104', '这家伙啥都不懂', '广东省广州市', '黄埔区开创大道2819号',
        '/4.jpg', '2019-10-12 00:00:00', '2019-10-12 00:00:00', FALSE);

INSERT
INTO `uap_tb_gmvcs_key_person_recognition`(
  `person_id`, `key_person_type_db_id`, `id_card`, `person_name`, `birth_date`, `key_person_detail`,
  `person_district`, `person_address`, `person_reg_img`, `create_time`, `update_time`, `is_deleted`)
VALUES ('440000198001050005', 1, '440000198001050005', '谢辉', '19800105', '长得太白', '广东省广州市', '黄埔区开创大道2819号',
        '/5.jpg', '2019-10-12 00:00:00', '2019-10-12 00:00:00', FALSE);

INSERT
INTO `uap_tb_gmvcs_key_person_recognition`(
  `person_id`, `key_person_type_db_id`, `id_card`, `person_name`, `birth_date`, `key_person_detail`,
  `person_district`, `person_address`, `person_reg_img`, `create_time`, `update_time`, `is_deleted`)
VALUES ('440000198001060006', 1, '440000198001060006', '李达祥', '19800106', '眼神太犀利', '广东省广州市', '黄埔区开创大道2819号',
        '/6.jpg', '2019-10-12 00:00:00', '2019-10-12 00:00:00', FALSE);

INSERT
INTO `uap_tb_gmvcs_key_person_recognition`(
  `person_id`, `key_person_type_db_id`, `id_card`, `person_name`, `birth_date`, `key_person_detail`,
  `person_district`, `person_address`, `person_reg_img`, `create_time`, `update_time`, `is_deleted`)
VALUES ('440000198001070007', 2, '440000198001070007', '魏永高', '19800107', '舔包不拉人', '广东省广州市', '黄埔区开创大道2819号',
        '/7.jpg', '2019-10-12 00:00:00', '2019-10-12 00:00:00', FALSE);

INSERT
INTO `uap_tb_gmvcs_key_person_recognition`(
  `person_id`, `key_person_type_db_id`, `id_card`, `person_name`, `birth_date`, `key_person_detail`,
  `person_district`, `person_address`, `person_reg_img`, `create_time`, `update_time`, `is_deleted`)
VALUES ('440000198001080008', 2, '440000198001080008', '周子海', '19800108', '枪打得太好', '广东省广州市', '黄埔区开创大道2819号',
        '/8.jpg', '2019-10-12 00:00:00', '2019-10-12 00:00:00', FALSE);

INSERT
INTO `uap_tb_gmvcs_key_person_recognition`(
  `person_id`, `key_person_type_db_id`, `id_card`, `person_name`, `birth_date`, `key_person_detail`,
  `person_district`, `person_address`, `person_reg_img`, `create_time`, `update_time`, `is_deleted`)
VALUES ('440000198001090009', 2, '440000198001090009', '陈国艺', '19800109', '机场高架王', '广东省广州市', '黄埔区开创大道2819号',
        '/9.jpg', '2019-10-12 00:00:00', '2019-10-12 00:00:00', FALSE);
