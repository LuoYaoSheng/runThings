/*
 Navicat Premium Data Transfer

 Source Server         : docker
 Source Server Type    : MySQL
 Source Server Version : 80028
 Source Host           : 127.0.0.1:3306
 Source Schema         : eq

 Target Server Type    : MySQL
 Target Server Version : 80028
 File Encoding         : 65001

 Date: 17/05/2022 16:16:40
*/

SET NAMES utf8mb4;
SET
FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for eq_log
-- ----------------------------
DROP TABLE IF EXISTS `eq_log`;
CREATE TABLE `eq_log`
(
    `id`          int                                                           NOT NULL AUTO_INCREMENT,
    `sn`          varchar(255)                                                  NOT NULL,
    `product_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
    `status`      int                                                           DEFAULT '0',
    `title`       varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '',
    `content`     varchar(255)                                                  DEFAULT '',
    `create_time` datetime(6) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP (6),
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='设备日志表';

-- ----------------------------
-- Table structure for eq_alarm_rule
-- ----------------------------
DROP TABLE IF EXISTS `eq_alarm_rule`;
CREATE TABLE `eq_alarm_rule`
(
    `id`         int unsigned NOT NULL AUTO_INCREMENT,
    `name`       varchar(255) NOT NULL                                         DEFAULT '' COMMENT '规则名',
    `level`      tinyint      NOT NULL                                         DEFAULT '0' COMMENT '告警等级',
    `code`       varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '产品code',
    `sn`         varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '设备sn',
    `content`    text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '告警规则',
    `created_by` bigint unsigned DEFAULT '0' COMMENT '创建者',
    `updated_by` bigint unsigned DEFAULT '0' COMMENT '更新者',
    `created_at` datetime                                                      DEFAULT NULL COMMENT '创建时间',
    `updated_at` datetime                                                      DEFAULT NULL COMMENT '更新时间',
    `deleted_at` datetime                                                      DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='告警规则表';

-- ----------------------------
-- Records of eq_alarm_rule
-- ----------------------------
BEGIN;
INSERT INTO `eq_alarm_rule`
VALUES (1, '温度过高', 1, '1100800013', '', '[{\"property\":\"temperature\",\"condition\":0,\"value\":70}]', 1, 1,
        '2022-05-21 10:26:13', '2022-05-21 10:58:26', NULL);
INSERT INTO `eq_alarm_rule`
VALUES (2, '温度过低', 1, '1100800013', '', '[{\"property\":\"temperature\",\"condition\":2,\"value\":10}]', 1, 0,
        '2022-05-21 10:58:49', '2022-05-21 10:58:49', NULL);
INSERT INTO `eq_alarm_rule`
VALUES (3, '温湿度异常', 0, '1100800013', '',
        '[{\"property\":\"temperature\",\"condition\":0,\"value\":50},{\"property\":\"humidity\",\"condition\":0,\"value\":60}]',
        1, 0, '2022-05-21 10:59:59', '2022-05-21 10:59:59', NULL);
COMMIT;

SET
FOREIGN_KEY_CHECKS = 1;
