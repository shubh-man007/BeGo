CREATE TABLE `users` (
  `id` int PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `is_admin` bool NOT NULL DEFAULT false
);

ALTER TABLE `orders`
    ADD COLUMN `user_id` int NOT NULL,
    ADD CONSTRAINT `user_id_fk` FOREIGN KEY (`user_id`)
        REFERENCES `users` (`id`);
