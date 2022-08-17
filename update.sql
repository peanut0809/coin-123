alter table `subscribe_activity` add column `disable` tinyint(4) NOT NULL DEFAULT '0' COMMENT '1.禁用';
alter table `seckill_activity` add column `activity_intro` varchar(255) NOT NULL COMMENT '活动简介';
alter table `seckill_activity` add column `cover_img_url` varchar(128) NOT NULL COMMENT '商品封面图';

DROP TABLE IF EXISTS `banner`;
CREATE TABLE `banner`  (
                           `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
                           `publisher_id` varchar(25) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '发行商id',
                           `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '活动名称',
                           `remarks` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '备注',
                           `image` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '图片',
                           `jump_type` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '页面跳转类型',
                           `jump_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '页面跳转地址',
                           `sort` int(11) UNSIGNED NULL DEFAULT NULL COMMENT '排序',
                           `state` int(1) UNSIGNED NULL DEFAULT 0 COMMENT '状态 0：未上架，1：已上架，2：已下架',
                           `timing_state` int(1) NOT NULL COMMENT '定时状态 1.启用定时器 2.关闭定时器',
                           `goods_on_time` datetime NULL DEFAULT NULL COMMENT '上架时间',
                           `goods_off_time` datetime NULL DEFAULT NULL COMMENT '下架时间',
                           `created_at` datetime NULL DEFAULT NULL COMMENT '创建时间',
                           `updated_at` datetime NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
                           PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 12 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;