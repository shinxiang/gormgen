
CREATE TABLE `user` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '用户id',
    `nickname` varchar(64) NOT NULL COMMENT '昵称',
    `username` varchar(64) NOT NULL COMMENT '用户名',
    `password` varchar(500) NOT NULL COMMENT '登录密码',
    `salt` varchar(36) DEFAULT NULL COMMENT '随机盐',
    `phone` varchar(20) DEFAULT NULL COMMENT '手机号',
    `status` varchar(1) NOT NULL DEFAULT 'Y' COMMENT '状态：Y-启用；N-禁用',
    `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted` tinyint(1) DEFAULT '0' COMMENT '逻辑删除：0-正常，1-删除',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='用户表';


CREATE TABLE `orders` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
    `name` varchar(255) DEFAULT NULL COMMENT '名称',
    `price` decimal(15,3) DEFAULT NULL COMMENT '订单价格',
    `create_time` datetime DEFAULT NULL COMMENT '创建时间',
    `update_time` datetime DEFAULT NULL COMMENT '更新时间',
    `deleted` tinyint(1) DEFAULT NULL COMMENT '删除状态，1表示软删除',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单信息表';
