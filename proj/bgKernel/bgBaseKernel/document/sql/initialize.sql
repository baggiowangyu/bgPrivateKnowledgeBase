# 系统配置表
# 这里主要用于配置“配置中心/注册中心”的参数
# 配置中心应支持Zookeeper、ETCD、Consul、或其他在线配置中心
DROP TABLE IF EXISTS `sys_base_config`;
CREATE TABLE `sys_base_config`(
    `Unique_id`         varchar(64)     NOT NULL PRIMARY KEY    COMMENT '唯一ID'
);

# 系统角色表
DROP TABLE IF EXISTS `sys_base_role`;
CREATE TABLE `sys_base_role`(
    `Unique_id`         varchar(64)     NOT NULL PRIMARY KEY    COMMENT '唯一ID',
    `Role_id`           varchar(64)     NOT NULL                COMMENT '角色ID',
    `Role_name`         varchar(64)     NOT NULL                COMMENT '角色名称',
#     `Role_range`        varchar(255)    NOT NULL                COMMENT '角色范围，也就是菜单ID'
);

# 系统管理员
INSERT INTO `sys_base_role`(`Unique_id`, `Role_id`, `Role_name`)
VALUES ('rid_r_0000001', 'Role_System_Admin', '系统管理员');
# 安全员
INSERT INTO `sys_base_role`(`Unique_id`, `Role_id`, `Role_name`)
VALUES ('rid_r_0000001', 'Role_Security_Admin', '安全管理员');
# 审计员
INSERT INTO `sys_base_role`(`Unique_id`, `Role_id`, `Role_name`)
VALUES ('rid_r_0000001', 'Role_Audit_Admin', '审计管理员');

# 系统菜单表
DROP TABLE IF EXISTS `sys_base_menu`;
CREATE TABLE `sys_base_menu`(
    `Unique_id`         varchar(64)     NOT NULL PRIMARY KEY    COMMENT '唯一ID',
    `Menu_id`           varchar(64)     NOT NULL                COMMENT '菜单ID',
    `Menu_father_id`    varchar(64)     NOT NULL                COMMENT '父菜单ID',
    `Menu_name`         varchar(64)     NOT NULL                COMMENT '菜单名称',
    `Menu_url`          varchar(512)    NOT NULL                COMMENT '菜单URL',
    `Menu_icon`         varchar(64)                             COMMENT '菜单图标',
    `Menu_show`         int                                     COMMENT '菜单显示状态：2-隐藏，1-显示',
    `Menu_type`         int                                     COMMENT '菜单类型：1-目录，2-菜单，3-按钮',
    `Create_time`       bigint          NOT NULL                COMMENT '创建时间，采用int64，排除时区问题',
    `Update_time`       bigint          NOT NULL                COMMENT '更新时间，采用int64，排除时区问题'
);

# 运维中心
INSERT INTO `sys_base_menu`(`Unique_id`, `Menu_id`, `Menu_father_id`, `Menu_name`, `Menu_url`, `Menu_icon`, `Menu_show`, `Menu_type`, `Create_time`, `Update_time`)
VALUES ('rid_m_000000001', 'MENU_maintaining_center', 'MENU_maintaining_center', '运维中心', '', '', 1, 1, 1579707083, 1579707083);
# 数据中心
INSERT INTO `sys_base_menu`(`Unique_id`, `Menu_id`, `Menu_father_id`, `Menu_name`, `Menu_url`, `Menu_icon`, `Menu_show`, `Menu_type`, `Create_time`, `Update_time`)
VALUES ('rid_m_000000001', 'MENU_maintaining_center', 'MENU_maintaining_center', '运维中心', '', '', 1, 1, 1579707083, 1579707083);
# 指挥中心
INSERT INTO `sys_base_menu`(`Unique_id`, `Menu_id`, `Menu_father_id`, `Menu_name`, `Menu_url`, `Menu_icon`, `Menu_show`, `Menu_type`, `Create_time`, `Update_time`)
VALUES ('rid_m_000000001', 'MENU_maintaining_center', 'MENU_maintaining_center', '运维中心', '', '', 1, 1, 1579707083, 1579707083);
# 监督中心
INSERT INTO `sys_base_menu`(`Unique_id`, `Menu_id`, `Menu_father_id`, `Menu_name`, `Menu_url`, `Menu_icon`, `Menu_show`, `Menu_type`, `Create_time`, `Update_time`)
VALUES ('rid_m_000000001', 'MENU_maintaining_center', 'MENU_maintaining_center', '运维中心', '', '', 1, 1, 1579707083, 1579707083);
#
INSERT INTO `sys_base_menu`(`Unique_id`, `Menu_id`, `Menu_father_id`, `Menu_name`, `Menu_url`, `Menu_icon`, `Menu_show`, `Menu_type`, `Create_time`, `Update_time`)
VALUES ('rid_m_000000001', 'MENU_maintaining_center', 'MENU_maintaining_center', '运维中心', '', '', 1, 1, 1579707083, 1579707083);

# 系统角色/菜单关联表
DROP TABLE IF EXISTS `sys_base_role_menu`;
CREATE TABLE `sys_base_role_menu`(
     `Unique_id`        int AUTO_INCREMENT NOT NULL PRIMARY KEY    COMMENT '唯一ID',
     `Role_id`          varchar(64)     NOT NULL                COMMENT '角色ID',
     `Menu_id`          varchar(64)     NOT NULL                COMMENT '菜单ID',
     `Reverse01`        varchar(255)                            COMMENT '预留字段'
);

# 组织架构表
DROP TABLE IF EXISTS `sys_base_organization`;
CREATE TABLE `sys_base_organization`(
    `Unique_id`         varchar(64)     NOT NULL PRIMARY KEY    COMMENT '唯一ID',
    `Org_code`          varchar(64)     NOT NULL                COMMENT '组织编码',
    `Org_name`          varchar(64)     NOT NULL                COMMENT '组织名称',
    `Org_path`          varchar(64)     NOT NULL                COMMENT '组织路径',
    `Parent_id`         varchar(64)     NOT NULL                COMMENT '父组织编码',
    `Create_time`       bigint          NOT NULL                COMMENT '创建时间，采用int64，排除时区问题',
    `Update_time`       bigint          NOT NULL                COMMENT '更新时间，采用int64，排除时区问题',
    `Source`            varchar(64)     NOT NULL                COMMENT '数据来源',
    `Order_index`       int             NOT NULL                COMMENT '排序索引',
    `Is_hide`           boolean         NOT NULL                COMMENT '组织是否隐藏',
    `Is_disable`        boolean         NOT NULL                COMMENT '组织是否禁用',
    `Extend`            varchar(255)                            COMMENT '扩展信息'
);

# 初始化组织架构
## 省级根节点
INSERT INTO `sys_base_organization`(`Unique_id`, `Org_code`, `Org_name`, `Org_path`, `Parent_id`, `Create_time`, `Update_time`, `Source`, `Order_index`, `Is_hide`, `Is_disable`, `Extend`)
VALUES ('rid_o_000000000', '440000', '广东省公安厅', '/440000', '440000', 1579707083, 1579707083, 'Internal', 1, 0, 0, '');

## 地市叶子节点
INSERT INTO `sys_base_organization`(`Unique_id`, `Org_code`, `Org_name`, `Org_path`, `Parent_id`, `Create_time`, `Update_time`, `Source`, `Order_index`, `Is_hide`, `Is_disable`, `Extend`)
VALUES ('rid_o_000000001', '440100', '广州市公安局', '/440000/440100', '440000', 1579707083, 1579707083, 'Internal', 1, 0, 0, '');
INSERT INTO `sys_base_organization`(`Unique_id`, `Org_code`, `Org_name`, `Org_path`, `Parent_id`, `Create_time`, `Update_time`, `Source`, `Order_index`, `Is_hide`, `Is_disable`, `Extend`)
VALUES ('rid_o_000000002', '440200', '韶关市公安局', '/440000/440200', '440000', 1579707083, 1579707083, 'Internal', 1, 0, 0, '');
INSERT INTO `sys_base_organization`(`Unique_id`, `Org_code`, `Org_name`, `Org_path`, `Parent_id`, `Create_time`, `Update_time`, `Source`, `Order_index`, `Is_hide`, `Is_disable`, `Extend`)
VALUES ('rid_o_000000003', '440300', '深圳市公安局', '/440000/440300', '440000', 1579707083, 1579707083, 'Internal', 1, 0, 0, '');
INSERT INTO `sys_base_organization`(`Unique_id`, `Org_code`, `Org_name`, `Org_path`, `Parent_id`, `Create_time`, `Update_time`, `Source`, `Order_index`, `Is_hide`, `Is_disable`, `Extend`)
VALUES ('rid_o_000000004', '440400', '珠海市公安局', '/440000/440400', '440000', 1579707083, 1579707083, 'Internal', 1, 0, 0, '');
INSERT INTO `sys_base_organization`(`Unique_id`, `Org_code`, `Org_name`, `Org_path`, `Parent_id`, `Create_time`, `Update_time`, `Source`, `Order_index`, `Is_hide`, `Is_disable`, `Extend`)
VALUES ('rid_o_000000005', '440500', '汕头市公安局', '/440000/440500', '440000', 1579707083, 1579707083, 'Internal', 1, 0, 0, '');
INSERT INTO `sys_base_organization`(`Unique_id`, `Org_code`, `Org_name`, `Org_path`, `Parent_id`, `Create_time`, `Update_time`, `Source`, `Order_index`, `Is_hide`, `Is_disable`, `Extend`)
VALUES ('rid_o_000000006', '440600', '佛山市公安局', '/440000/440600', '440000', 1579707083, 1579707083, 'Internal', 1, 0, 0, '');
INSERT INTO `sys_base_organization`(`Unique_id`, `Org_code`, `Org_name`, `Org_path`, `Parent_id`, `Create_time`, `Update_time`, `Source`, `Order_index`, `Is_hide`, `Is_disable`, `Extend`)
VALUES ('rid_o_000000007', '440700', '江门市公安局', '/440000/440700', '440000', 1579707083, 1579707083, 'Internal', 1, 0, 0, '');
INSERT INTO `sys_base_organization`(`Unique_id`, `Org_code`, `Org_name`, `Org_path`, `Parent_id`, `Create_time`, `Update_time`, `Source`, `Order_index`, `Is_hide`, `Is_disable`, `Extend`)
VALUES ('rid_o_000000008', '440800', '湛江市公安局', '/440000/440800', '440000', 1579707083, 1579707083, 'Internal', 1, 0, 0, '');
INSERT INTO `sys_base_organization`(`Unique_id`, `Org_code`, `Org_name`, `Org_path`, `Parent_id`, `Create_time`, `Update_time`, `Source`, `Order_index`, `Is_hide`, `Is_disable`, `Extend`)
VALUES ('rid_o_000000009', '440900', '茂名市公安局', '/440000/440900', '440000', 1579707083, 1579707083, 'Internal', 1, 0, 0, '');
INSERT INTO `sys_base_organization`(`Unique_id`, `Org_code`, `Org_name`, `Org_path`, `Parent_id`, `Create_time`, `Update_time`, `Source`, `Order_index`, `Is_hide`, `Is_disable`, `Extend`)
VALUES ('rid_o_000000010', '441200', '肇庆市公安局', '/440000/441200', '440000', 1579707083, 1579707083, 'Internal', 1, 0, 0, '');
INSERT INTO `sys_base_organization`(`Unique_id`, `Org_code`, `Org_name`, `Org_path`, `Parent_id`, `Create_time`, `Update_time`, `Source`, `Order_index`, `Is_hide`, `Is_disable`, `Extend`)
VALUES ('rid_o_000000011', '441300', '惠州市公安局', '/440000/441300', '440000', 1579707083, 1579707083, 'Internal', 1, 0, 0, '');
INSERT INTO `sys_base_organization`(`Unique_id`, `Org_code`, `Org_name`, `Org_path`, `Parent_id`, `Create_time`, `Update_time`, `Source`, `Order_index`, `Is_hide`, `Is_disable`, `Extend`)
VALUES ('rid_o_000000012', '441400', '梅州市公安局', '/440000/441400', '440000', 1579707083, 1579707083, 'Internal', 1, 0, 0, '');
INSERT INTO `sys_base_organization`(`Unique_id`, `Org_code`, `Org_name`, `Org_path`, `Parent_id`, `Create_time`, `Update_time`, `Source`, `Order_index`, `Is_hide`, `Is_disable`, `Extend`)
VALUES ('rid_o_000000013', '441500', '汕尾市公安局', '/440000/441500', '440000', 1579707083, 1579707083, 'Internal', 1, 0, 0, '');
INSERT INTO `sys_base_organization`(`Unique_id`, `Org_code`, `Org_name`, `Org_path`, `Parent_id`, `Create_time`, `Update_time`, `Source`, `Order_index`, `Is_hide`, `Is_disable`, `Extend`)
VALUES ('rid_o_000000014', '441600', '河源市公安局', '/440000/441600', '440000', 1579707083, 1579707083, 'Internal', 1, 0, 0, '');
INSERT INTO `sys_base_organization`(`Unique_id`, `Org_code`, `Org_name`, `Org_path`, `Parent_id`, `Create_time`, `Update_time`, `Source`, `Order_index`, `Is_hide`, `Is_disable`, `Extend`)
VALUES ('rid_o_000000015', '441700', '阳江市公安局', '/440000/441700', '440000', 1579707083, 1579707083, 'Internal', 1, 0, 0, '');
INSERT INTO `sys_base_organization`(`Unique_id`, `Org_code`, `Org_name`, `Org_path`, `Parent_id`, `Create_time`, `Update_time`, `Source`, `Order_index`, `Is_hide`, `Is_disable`, `Extend`)
VALUES ('rid_o_000000016', '441800', '清远市公安局', '/440000/441800', '440000', 1579707083, 1579707083, 'Internal', 1, 0, 0, '');
INSERT INTO `sys_base_organization`(`Unique_id`, `Org_code`, `Org_name`, `Org_path`, `Parent_id`, `Create_time`, `Update_time`, `Source`, `Order_index`, `Is_hide`, `Is_disable`, `Extend`)
VALUES ('rid_o_000000017', '441900', '东莞市公安局', '/440000/441900', '440000', 1579707083, 1579707083, 'Internal', 1, 0, 0, '');
INSERT INTO `sys_base_organization`(`Unique_id`, `Org_code`, `Org_name`, `Org_path`, `Parent_id`, `Create_time`, `Update_time`, `Source`, `Order_index`, `Is_hide`, `Is_disable`, `Extend`)
VALUES ('rid_o_000000018', '442000', '中山市公安局', '/440000/442000', '440000', 1579707083, 1579707083, 'Internal', 1, 0, 0, '');
INSERT INTO `sys_base_organization`(`Unique_id`, `Org_code`, `Org_name`, `Org_path`, `Parent_id`, `Create_time`, `Update_time`, `Source`, `Order_index`, `Is_hide`, `Is_disable`, `Extend`)
VALUES ('rid_o_000000019', '445100', '潮州市公安局', '/440000/445100', '440000', 1579707083, 1579707083, 'Internal', 1, 0, 0, '');
INSERT INTO `sys_base_organization`(`Unique_id`, `Org_code`, `Org_name`, `Org_path`, `Parent_id`, `Create_time`, `Update_time`, `Source`, `Order_index`, `Is_hide`, `Is_disable`, `Extend`)
VALUES ('rid_o_000000020', '445200', '揭阳市公安局', '/440000/445200', '440000', 1579707083, 1579707083, 'Internal', 1, 0, 0, '');
INSERT INTO `sys_base_organization`(`Unique_id`, `Org_code`, `Org_name`, `Org_path`, `Parent_id`, `Create_time`, `Update_time`, `Source`, `Order_index`, `Is_hide`, `Is_disable`, `Extend`)
VALUES ('rid_o_000000021', '445300', '云浮市公安局', '/440000/445300', '440000', 1579707083, 1579707083, 'Internal', 1, 0, 0, '');

## 区县叶子节点
### 广州市属区县
INSERT INTO `sys_base_organization`(`Unique_id`, `Org_code`, `Org_name`, `Org_path`, `Parent_id`, `Create_time`, `Update_time`, `Source`, `Order_index`, `Is_hide`, `Is_disable`, `Extend`)
VALUES ('rid_o_000000101', '440103', '荔湾区公安分局', '/440000/440100/440103', '440100', 1579707083, 1579707083, 'Internal', 1, 0, 0, '');
INSERT INTO `sys_base_organization`(`Unique_id`, `Org_code`, `Org_name`, `Org_path`, `Parent_id`, `Create_time`, `Update_time`, `Source`, `Order_index`, `Is_hide`, `Is_disable`, `Extend`)
VALUES ('rid_o_000000102', '440104', '越秀区公安分局', '/440000/440100/440104', '440100', 1579707083, 1579707083, 'Internal', 1, 0, 0, '');
INSERT INTO `sys_base_organization`(`Unique_id`, `Org_code`, `Org_name`, `Org_path`, `Parent_id`, `Create_time`, `Update_time`, `Source`, `Order_index`, `Is_hide`, `Is_disable`, `Extend`)
VALUES ('rid_o_000000103', '440105', '海珠区公安分局', '/440000/440100/440105', '440100', 1579707083, 1579707083, 'Internal', 1, 0, 0, '');
INSERT INTO `sys_base_organization`(`Unique_id`, `Org_code`, `Org_name`, `Org_path`, `Parent_id`, `Create_time`, `Update_time`, `Source`, `Order_index`, `Is_hide`, `Is_disable`, `Extend`)
VALUES ('rid_o_000000104', '440106', '天河区公安分局', '/440000/440100/440106', '440100', 1579707083, 1579707083, 'Internal', 1, 0, 0, '');
INSERT INTO `sys_base_organization`(`Unique_id`, `Org_code`, `Org_name`, `Org_path`, `Parent_id`, `Create_time`, `Update_time`, `Source`, `Order_index`, `Is_hide`, `Is_disable`, `Extend`)
VALUES ('rid_o_000000105', '440111', '白云区公安分局', '/440000/440100/440111', '440100', 1579707083, 1579707083, 'Internal', 1, 0, 0, '');
INSERT INTO `sys_base_organization`(`Unique_id`, `Org_code`, `Org_name`, `Org_path`, `Parent_id`, `Create_time`, `Update_time`, `Source`, `Order_index`, `Is_hide`, `Is_disable`, `Extend`)
VALUES ('rid_o_000000106', '440112', '黄埔区公安分局', '/440000/440100/440112', '440100', 1579707083, 1579707083, 'Internal', 1, 0, 0, '');
INSERT INTO `sys_base_organization`(`Unique_id`, `Org_code`, `Org_name`, `Org_path`, `Parent_id`, `Create_time`, `Update_time`, `Source`, `Order_index`, `Is_hide`, `Is_disable`, `Extend`)
VALUES ('rid_o_000000107', '440113', '番禺区公安分局', '/440000/440100/440113', '440100', 1579707083, 1579707083, 'Internal', 1, 0, 0, '');
INSERT INTO `sys_base_organization`(`Unique_id`, `Org_code`, `Org_name`, `Org_path`, `Parent_id`, `Create_time`, `Update_time`, `Source`, `Order_index`, `Is_hide`, `Is_disable`, `Extend`)
VALUES ('rid_o_000000108', '440114', '花都区公安分局', '/440000/440100/440114', '440100', 1579707083, 1579707083, 'Internal', 1, 0, 0, '');
INSERT INTO `sys_base_organization`(`Unique_id`, `Org_code`, `Org_name`, `Org_path`, `Parent_id`, `Create_time`, `Update_time`, `Source`, `Order_index`, `Is_hide`, `Is_disable`, `Extend`)
VALUES ('rid_o_000000109', '440115', '南沙区公安分局', '/440000/440100/440115', '440100', 1579707083, 1579707083, 'Internal', 1, 0, 0, '');
INSERT INTO `sys_base_organization`(`Unique_id`, `Org_code`, `Org_name`, `Org_path`, `Parent_id`, `Create_time`, `Update_time`, `Source`, `Order_index`, `Is_hide`, `Is_disable`, `Extend`)
VALUES ('rid_o_000000110', '440117', '从化区公安分局', '/440000/440100/440117', '440100', 1579707083, 1579707083, 'Internal', 1, 0, 0, '');
INSERT INTO `sys_base_organization`(`Unique_id`, `Org_code`, `Org_name`, `Org_path`, `Parent_id`, `Create_time`, `Update_time`, `Source`, `Order_index`, `Is_hide`, `Is_disable`, `Extend`)
VALUES ('rid_o_000000111', '440118', '增城区公安分局', '/440000/440100/440118', '440100', 1579707083, 1579707083, 'Internal', 1, 0, 0, '');

DROP TABLE IF EXISTS `sys_base_user`;
CREATE TABLE `sys_base_user`(
    `Unique_id`         varchar(64)     NOT NULL PRIMARY KEY    COMMENT '唯一ID',
    `User_code`         varchar(64)     NOT NULL                COMMENT '用户编码（相当于工号，警察为警号，辅警为辅警号，其他人为相对应工号），也作为系统登录账号',
    `User_pin`          varchar(256)    NOT NULL                COMMENT '用户口令',
    `User_name`         varchar(64)     NOT NULL                COMMENT '用户名称',
    `Org_id`            varchar(64)     NOT NULL                COMMENT '所在组织ID',
    `Role_id`           varchar(64)     NOT NULL                COMMENT '系统角色（系统角色默认有三个：系统管理员、安全管理员、审计管理员）',
    `Is_disable`        boolean         NOT NULL                COMMENT '是否被标记为禁用',
    `Is_deleted`        boolean         NOT NULL                COMMENT '是否被标记为删除',
    `Can_multi_login`   boolean         NOT NULL                COMMENT '是否允许多重登录',
    `Gender`            int             NOT NULL                COMMENT '性别：0-女；1-男；2-其他',
    `Create_time`       bigint          NOT NULL                COMMENT '创建时间',
    `Update_time`       bigint          NOT NULL                COMMENT '更新时间',
    `Source`            varchar(64)     NOT NULL                COMMENT '数据来源'
);

# 初始化用户数据
## 系统内置三个管理员，默认密码："123456"，默认盐"gosuncn"
INSERT INTO `sys_base_user`(`Unique_id`, `User_code`, `User_pin`, `User_name`, `Org_id`, `Role_id`, `Is_disable`, `Is_deleted`, `Can_multi_login`, `Gender`, `Create_time`, `Update_time`, `Source`)
VALUES ('rid_u_000000000', 'sysadmin', '74026ff78f299e9f0be4a81ddcdac0e14cee8fff3e8af34a70bba536e0373348', '系统管理员', '000000', '', 0, 0, 1, 1, 1579707083, 1579707083, 'Internal');