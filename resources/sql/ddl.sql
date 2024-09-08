create database my_ecommerce_system;
use my_ecommerce_system;

CREATE TABLE sys_user (
    id BIGINT NOT NULL,
    name VARCHAR(50) NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    PRIMARY KEY (id)
);
