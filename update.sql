alter table `subscribe_activity` add column `disable` tinyint(4) NOT NULL DEFAULT '0' COMMENT '1.禁用';
alter table `seckill_activity` add column `activity_intro` varchar(255) NOT NULL COMMENT '活动简介';
alter table `seckill_activity` add column `cover_img_url` varchar(128) NOT NULL COMMENT '商品封面图';
alter table `seckill_activity` add column `disable` tinyint(4) NOT NULL DEFAULT '0' COMMENT '1.禁用';
alter table `seckill_orders` add column `publisher_id` varchar(32) NOT NULL DEFAULT '' COMMENT '发行商ID';