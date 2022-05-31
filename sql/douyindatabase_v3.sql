-- user
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`
(
    `id`           bigint unsigned NOT NULL AUTO_INCREMENT COMMENT ='用户id',
    `user_name`    varchar(100)    NOT NULL UNIQUE COMMENT ='用户名,必须唯一',
    `password`     varchar(100)    NOT NULL COMMENT = '加密过的密码',
    `follow_count` bigint unsigned DEFAULT 0 COMMENT ='关注数',
    `fans_count`   bigint unsigned DEFAULT 0 COMMENT ='粉丝数',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

INSERT INTO user(user_name, password) VALUES ('gsz', '123456');


-- relation
DROP TABLE IF EXISTS `relation`;
CREATE TABLE `relation`
(
    `id`        bigint unsigned NOT NULL AUTO_INCREMENT,
    `user_id`   bigint unsigned NOT NULL COMMENT ='用户的id',
    `follow_id` bigint unsigned NOT NULL COMMENT ='用户关注的用户的id',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_relation_follow` (`user_id`, `follow_id`) COMMENT ='方便查询用户关注的用户，且不能重复',
    INDEX `idx_relation_fans` (`follow_id`, `user_id`) COMMENT ='方便查询用户的粉丝'
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

-- video
DROP TABLE IF EXISTS `video`;
CREATE TABLE `video`
(
    `id`             bigint unsigned NOT NULL AUTO_INCREMENT COMMENT ='视频id',
    `author_id`      bigint unsigned NOT NULL COMMENT ='视频的作者id',
    `play_url`       varchar(255)    NOT NULL COMMENT ='视频文件的存储地址',
    `cover_url`      varchar(255)    NOT NULL COMMENT ='视频封面的地址',
    `favorite_count` bigint unsigned DEFAULT 0 COMMENT ='视频的点赞数',
    `comment_count`  bigint unsigned DEFAULT 0 COMMENT ='视频的评论数',
    `title`          varchar(255)    DEFAULT NULL COMMENT ='视频标题',
    `created_at`     bigint unsigned NOT NULL COMMENT ='视频的创建日期',
    PRIMARY KEY (`id`),
    INDEX `idx_author_id` (`author_id`),
    INDEX `idx_video_created_at` (`created_at`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;


-- favorite
DROP TABLE IF EXISTS `favorite`;
CREATE TABLE `favorite`
(
    `user_id`  bigint unsigned NOT NULL COMMENT ='用户id',
    `video_id` bigint unsigned NOT NULL COMMENT ='用户点赞的视频id',
    PRIMARY KEY (`user_id`, `video_id`) COMMENT ='联合主键保证用户点赞某个视频不会重复，因为不需要根据视频id查询点赞的用户，所以不需要建video_id索引'
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

-- comment
DROP TABLE IF EXISTS `comment`;
CREATE TABLE `comment`
(
    `id`          bigint unsigned NOT NULL AUTO_INCREMENT,
    `user_id`     bigint unsigned NOT NULL COMMENT ='发布评论的用户的id',
    `video_id`    bigint unsigned NOT NULL COMMENT ='被评论的视频的id',
    `content`     text            NOT NULL COMMENT ='评论内容',
    `create_date` bigint unsigned NOT NULL COMMENT ='评论发布的时间',
    PRIMARY KEY (`id`),
    INDEX `idx_comment_video_id` (`video_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;