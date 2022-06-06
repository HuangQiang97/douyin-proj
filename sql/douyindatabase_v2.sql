-- user
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_name` varchar(100) NOT NULL UNIQUE,
  `password` varchar(100) NOT NULL,
  `follow_count` bigint unsigned DEFAULT 0 ,
  `fans_count` bigint unsigned DEFAULT 0 ,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


INSERT INTO user(user_name,password) VALUES('gsz','123456');


-- relation
DROP TABLE IF EXISTS `relation`;
CREATE TABLE `relation` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `follow_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_relation_follow` (`user_id`,`follow_id`),
  INDEX `idx_relation_fans` (`follow_id`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- video
DROP TABLE IF EXISTS `video`;
CREATE TABLE `video` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `author_id` bigint unsigned NOT NULL,
  `play_url` varchar(255) NOT NULL,
  `cover_url` varchar(255) NOT NULL,
  `favorite_count` bigint unsigned DEFAULT 0,
  `comment_count` bigint unsigned DEFAULT 0,
  `title` varchar(255) DEFAULT NULL,
  `created_at` bigint unsigned NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_author_id` (`author_id`),
  INDEX `idx_video_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;



-- favorite
DROP TABLE IF EXISTS `favorite`;
CREATE TABLE `favorite` (
  `user_id` bigint unsigned NOT NULL,
  `video_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`user_id`,`video_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- comment
DROP TABLE IF EXISTS `comment`;
CREATE TABLE `comment` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `video_id` bigint unsigned NOT NULL,
  `content` text NOT NULL,
  `create_date` bigint unsigned NOT NULL,
   PRIMARY KEY (`id`),
   INDEX `idx_comment_video_id` (`video_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;










