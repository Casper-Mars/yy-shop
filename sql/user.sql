create table user
(
    id       int(64)      not null auto_increment,
    username varchar(64)  not null,
    password varchar(128) not null,
    phone    varchar(18),
    nickname varchar(20),
    avatar   varchar(50),
    PRIMARY KEY (id)
) engine = innodb
  default charset = utf8mb4;