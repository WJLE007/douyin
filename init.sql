CREATE TABLE `user` (
                        `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id主键',
                        `name` varchar(128) NOT NULL DEFAULT '' COMMENT '用户名',
                        `password` varchar(128) NOT NULL DEFAULT '' COMMENT '密码',
                        `avatar` varchar(255) NOT NULL DEFAULT '' COMMENT '头像地址',
                        `follow_count` bigint NOT NULL DEFAULT '0' COMMENT '关注总数',
                        `follower_count` bigint NOT NULL DEFAULT '0' COMMENT '粉丝总数',
                        `total_favorited` bigint NOT NULL DEFAULT '0' COMMENT '总点赞数',
                        `signature` varchar(255) NOT NULL DEFAULT '' COMMENT '个性签名',
                        `background_image` varchar(255) NOT NULL DEFAULT '' COMMENT '背景图片',
                        `work_count` bigint NOT NULL DEFAULT '0' COMMENT '作品总数',
                        `favorite_count` bigint NOT NULL DEFAULT '0' COMMENT '点赞总数',
                        `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                        `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                        `deleted_at` int DEFAULT '0' COMMENT '软删除',
                        PRIMARY KEY (`id`),
                        UNIQUE KEY `user_name_uindex` (`name`),
                        KEY `user_deleted_at_index` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户表';



CREATE TABLE `video` (
                         `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id主键',
                         `author_id` bigint NOT NULL DEFAULT '0' COMMENT '作者id',
                         `title` varchar(128) NOT NULL DEFAULT '' COMMENT '视频标题',
                         `play_url` varchar(255) NOT NULL DEFAULT '' COMMENT '视频播放地址',
                         `cover_url` varchar(255) NOT NULL DEFAULT '' COMMENT '视频封面地址',
                         `favorite_count` bigint NOT NULL DEFAULT '0' COMMENT '点赞数',
                         `comment_count` bigint NOT NULL DEFAULT '0' COMMENT '评论数',
                         `create_time` bigint NOT NULL COMMENT '创建时间',
                         `deleted_at` int DEFAULT '0' COMMENT '软删除',
                         PRIMARY KEY (`id`),
                         KEY `follow_deleted_at_index` (`deleted_at`),
                         KEY `video_create_time_index` (`create_time`)
) ENGINE=InnoDB AUTO_INCREMENT=35 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='视频表';

CREATE TABLE `message` (
                           `id` bigint NOT NULL AUTO_INCREMENT,
                           `to_user_id` bigint NOT NULL COMMENT '该消息接收者的id',
                           `from_user_id` bigint NOT NULL COMMENT '该消息发送者的id',
                           `content` varchar(255) NOT NULL COMMENT '内容',
                           `create_time` bigint NOT NULL COMMENT '创建时间',
                           PRIMARY KEY (`id`),
                           KEY `message_to_user_id_from_user_id_index` (`to_user_id`,`from_user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=49 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='消息表';


CREATE TABLE `follow` (
                          `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id主键',
                          `user_id` bigint NOT NULL DEFAULT '0' COMMENT '用户id',
                          `follow_id` bigint NOT NULL DEFAULT '0' COMMENT '关注用户id',
                          `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                          PRIMARY KEY (`id`),
                          UNIQUE KEY `user_follow_uindex` (`user_id`,`follow_id`)
) ENGINE=InnoDB AUTO_INCREMENT=24 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='关注表';

CREATE TABLE `favorite` (
                            `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id主键',
                            `user_id` bigint NOT NULL DEFAULT '0' COMMENT '用户id',
                            `video_id` bigint NOT NULL DEFAULT '0' COMMENT '喜欢作品id',
                            `create_time` bigint DEFAULT NULL,
                            PRIMARY KEY (`id`),
                            UNIQUE KEY `user_follow_uindex` (`user_id`,`video_id`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='点赞表';


CREATE TABLE `comment` (
                           `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id主键',
                           `user_id` bigint NOT NULL DEFAULT '0' COMMENT '用户id',
                           `video_id` bigint NOT NULL DEFAULT '0' COMMENT '视频id',
                           `favorite_count` bigint NOT NULL DEFAULT '0' COMMENT '点赞数',
                           `content` varchar(255) NOT NULL DEFAULT '' COMMENT '评论内容',
                           `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                           PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=23 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='评论表';









