create database if not exists `payment_bridge` charset = 'utf8';
use `payment_bridge`;

create table if not exists `assets`
(
    user_address char(42) not null,
    title varchar(255),
    image_url text,
    image_preview_url text,
    image_thumbnail_url text,
    description text,
    contract_address char(42) not null,
    token_id varchar(128) not null,
    num_sales int(10),
    owner char(42) not null,
    owner_img_url text,
    creator char(42) not null,
    creator_img_url text,
    token_metadata text,
    slug varchar(255) not null,
    is_delete int(8)
) ENGINE = InnoDB DEFAULT CHARSET=utf8mb4;

create table if not exists `collections`
(
    slug varchar(255) not null,
    owner char(42) not null,
    name varchar(255),
    banner_image_url text,
    description text,
    image_url text,
    large_image_url text,
    is_delete int(8),
    create_date varchar(32)
) ENGINE = InnoDB DEFAULT CHARSET=utf8mb4;

create table if not exists `contracts`
(
    address char(42) not null,
    contract_name varchar(255),
    contract_type varchar(64),
    symbol varchar(64),
    schema_name varchar(32),
    total_supply varchar(128),
    Description text
) ENGINE = InnoDB DEFAULT CHARSET=utf8mb4;

create table if not exists `assets_top_ownerships`
(
    contract_address char(42) not null,
    token_id varchar(128) not null,
    owner char(42) not null,
    profile_img_url text,
    quantity varchar(128)
) ENGINE = InnoDB DEFAULT CHARSET=utf8mb4;

create table if not exists `traits`
(
    contract_address char(42) not null,
    token_id varchar(128) not null,
    trait_type varchar(255),
    value varchar(255),
    display_type varchar(255),
    max_value int(10),
    trait_count int(10),
    order varchar(255),
    is_delete int(8)
) ENGINE = InnoDB DEFAULT CHARSET=utf8mb4;