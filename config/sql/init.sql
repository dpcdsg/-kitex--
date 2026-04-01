create table tiktok.`user`
(
    `id`               bigint auto_increment not null,
    `username`         varchar(255)                                                                 not null unique,
    `password`         varchar(255)                                                                 not null,
    `avatar`           varchar(255) default 'https://files.ozline.icu/images/avatar.jpg'            not null comment 'url',
    `signature`        varchar(255) default ''                                                      not null comment '255charmax',
    `created_at`       timestamp    default current_timestamp                                       not null,
    `updated_at`       timestamp    default current_timestamp                                       not null on update current_timestamp,
    `deleted_at`       timestamp    default null null,
    constraint `id`
        primary key (`id`)
) engine=InnoDB auto_increment=10000 default charset=utf8mb4;

create table tiktok.`product`
(
    `id`          bigint auto_increment not null,
    `seller_id`   bigint                              not null,
    `name`        varchar(255)                        not null,
    `description` text                                not null,
    `price`       bigint                              not null comment 'unit: cent',
    `stock`       bigint                              not null def·ault 0,
    `image_url`   varchar(512)                        not null default '',
    `category`    varchar(64)                         not null default '',
    `status`      tinyint   default 1                 not null comment '1:on 0:off',
    `created_at`  timestamp default current_timestamp not null,
    `updated_at`  timestamp default current_timestamp not null on update current_timestamp,
    `deleted_at`  timestamp default null null,
    constraint `id`
        primary key (`id`),
        index `idx_seller` (`seller_id`),
        index `idx_category` (`category`),
        index `idx_status` (`status`)
) engine=InnoDB default charset=utf8mb4;

create table tiktok.`seckill_activity`
(
    `id`              bigint auto_increment not null,
    `product_id`      bigint                              not null,
    `seckill_price`   bigint                              not null comment 'unit: cent',
    `total_stock`     bigint                              not null,
    `available_stock` bigint                              not null,
    `start_time`      timestamp                           not null,
    `end_time`        timestamp                           not null,
    `status`          tinyint   default 0                 not null comment '0:pending 1:active 2:ended',
    `created_at`      timestamp default current_timestamp not null,
    `updated_at`      timestamp default current_timestamp not null on update current_timestamp,
    `deleted_at`      timestamp default null null,
    constraint `id`
        primary key (`id`),
        index `idx_product` (`product_id`),
        index `idx_status` (`status`),
        index `idx_time` (`start_time`, `end_time`),
    constraint `seckill_product`
        foreign key (`product_id`)
            references tiktok.`product` (`id`)
            on delete cascade
) engine=InnoDB default charset=utf8mb4;

create table tiktok.`order`
(
    `id`            bigint auto_increment not null,
    `order_no`      varchar(64)                         not null unique,
    `user_id`       bigint                              not null,
    `product_id`    bigint                              not null,
    `activity_id`   bigint                              not null,
    `amount`        bigint                              not null comment 'unit: cent',
    `status`        tinyint   default 0                 not null comment '0:pending 1:paid 2:cancelled 3:expired',
    `created_at`    timestamp default current_timestamp not null,
    `updated_at`    timestamp default current_timestamp not null on update current_timestamp,
    `deleted_at`    timestamp default null null,
    constraint `id`
        primary key (`id`),
        unique index `idx_order_no` (`order_no`),
        index `idx_user` (`user_id`),
        index `idx_user_activity` (`user_id`, `activity_id`),
        index `idx_status` (`status`),
    constraint `order_user`
        foreign key (`user_id`)
            references tiktok.`user` (`id`)
            on delete cascade,
    constraint `order_product`
        foreign key (`product_id`)
            references tiktok.`product` (`id`)
            on delete cascade,
    constraint `order_activity`
        foreign key (`activity_id`)
            references tiktok.`seckill_activity` (`id`)
            on delete cascade
) engine=InnoDB default charset=utf8mb4;
