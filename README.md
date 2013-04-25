Create MySQL DB with the following schema.

create table `pages` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `title` varchar(100),
    `body` text,
    `created_at` datetime,
    `updated_at` datetime,
    primary key (`id`)
);
