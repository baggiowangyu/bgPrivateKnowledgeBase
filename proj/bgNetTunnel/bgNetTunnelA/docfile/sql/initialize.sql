DROP TABLE IF EXISTS `bg_mapping_table`;
CREATE TABLE `bg_mapping_table`(
  `Mapping_id`    INTEGER       NOT NULL    PRIMARY KEY     AUTO_INCREMENT  COMMENT '映射ID，自增',
  `Mapping_port`  INTEGER       NOT NULL                                    COMMENT '本地映射的PORT',
  `Source_ip`     VARCHAR(255)  NOT NULL                                    COMMENT '源IP',
  `Source_port`   INTEGER       NOT NULL                                    COMMENT '源PORT',
  `Net_type`      VARCHAR(32)   NOT NULL                                    COMMENT '映射的网络类型，只能是TCP或其他，其他一律是UDP',
  `Is_running`    INTEGER       NOT NULL                                    COMMENT '是否处于运行状态'
);

INSERT
INTO `bg_mapping_table`(`Mapping_port`, `Source_ip`, `Source_port`, `Net_type`, `Is_running`)
VALUES(9991, '192.168.231.1', 80, 'TCP', 1);

# INSERT
# INTO `bg_mapping_table`(`Mapping_port`, `Source_ip`, `Source_port`, `Net_type`)
# VALUES(9992, '192.168.231.1', 81, 'TCP');
#
# INSERT
# INTO `bg_mapping_table`(`Mapping_port`, `Source_ip`, `Source_port`, `Net_type`)
# VALUES(9993, '192.168.231.1', 82, 'TCP');
#
# INSERT
# INTO `bg_mapping_table`(`Mapping_port`, `Source_ip`, `Source_port`, `Net_type`)
# VALUES(9994, '192.168.231.1', 83, 'TCP');
#
# INSERT
# INTO `bg_mapping_table`(`Mapping_port`, `Source_ip`, `Source_port`, `Net_type`)
# VALUES(9995, '192.168.231.1', 84, 'TCP');