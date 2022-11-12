create table item
(
    id           int(10) unsigned not null auto_increment,
    item_name    varchar(256)     not null,
    icon_url     varchar(256)     null,
    price        decimal(10, 2)   not null default 0.0,
    seller_id    int(10) unsigned not null,
    booked_cnt   int(10) unsigned not null default 0,
    is_del       tinyint(1)       not null default 0 comment '0.未删除 1.已删除',
    effect_begin timestamp        not null default 0 comment '上架时间',
    effect_end   timestamp        not null default 0 comment '下架时间',
    update_time  timestamp        not null default current_timestamp on update current_timestamp,
    create_time  timestamp        not null default current_timestamp,
    images       text             null,
    PRIMARY KEY (id),
    KEY `seller_id_index` (seller_id)
) engine = innodb
  default charset = utf8mb4;
