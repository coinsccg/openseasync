create database if not exists `payment_bridge` charset = 'utf8';
use `payment_bridge`;

create table if not exists `assets`
(
    title varchar(255),
    image_url text,
    image_preview_url text,
    image_thumbnail_url text,
    description text,
    address char(42) not null,
    token_id varchar(128) not null,
    contract_name varchar(64),
    symbol varchar(32) not null,
    schema_name varchar(10) not null,
    total_supply varchar(64),
    owner char(42) not null,
    owner_img_url text,
    creator char(42) not null,
    creator_img_url text,
    create_date varchar(64),
    slug varchar(255) not null,
    collection_img_url text,
    collection_banner_img_url text,
    collection_description text,
    collection_large_img_url text
) ENGINE = InnoDB DEFAULT CHARSET=utf8mb4;

create table if not exists `collections`
(
    slug varchar(255) not null,
    owner char(42) not null,
    name varchar(255),
    banner_image_url text,
    description text,
    image_url text,
    large_image_url text
) ENGINE = InnoDB DEFAULT CHARSET=utf8mb4;