drop table if exists orders;
drop table if exists order_type;
drop table if exists order_status;
drop table if exists transactions;
drop table if exists symbols;
drop table if exists accounts;
drop table if exists account_statuses;
drop table if exists audit_log;
drop table if exists users;

create table `users`
(
	user_id varchar(1024) primary key,
	username varchar(255) not null unique,
	registered datetime not null default current_timestamp,
	last_action datetime default current_timestamp
);

create table `audit_log`
(
	audit_id bigint primary key auto_increment,
	user_id varchar(1024) not null,
	`date` datetime not null default current_timestamp,
	action varchar(255) not null,
	notes text,
	foreign key (user_id) references `users`(user_id)
);


create table `account_statuses`
(
	account_status_id bigint primary key auto_increment,
	name varchar(255) not null unique,
	`type` varchar(255) not null unique
);


create table `accounts`
(
	account_id bigint primary key auto_increment,
	account_name varchar(255) not null,
	user_id varchar(1024) not null,
	created datetime not null default current_timestamp,
	account_status_id bigint not null,	-- Is the account open or closed (from previous attempts?)
	foreign key (user_id) references `users`(user_id),
	foreign key (account_status_id) references `account_statuses`(account_status_id)
);

-- Symbols (like MCD, IVV) and any attributes associated with them of 

create table `symbols`
(
	symbol varchar(255) primary key
);

insert into symbols (symbol) values ('SPAXX');	-- Money market fund for new accounts

--  Positions on an account
-- drop table if exists positions;
-- create table `positions`
-- (
-- 	position_id bigint primary key auto_increment,
-- 	account_id bigint not null,
-- 	created datetime not null default current_timestamp,
-- 	foreign key (account_id) references `accounts`(account_id)
-- );

-- drop table if exists transaction_type;
-- create table `transaction_type`
-- (
-- 	transaction_type_id bigint primary key auto_increment,
-- 	name varchar(255) not null unique,
-- 	`type` varchar(255) not null unique
-- );
-- 
-- insert into transaction_type (name,`type`) values ('Trade','trade');
-- insert into transaction_type (name,`type`) values ('Transfer','transfer');

-- Transactions that occurred (buy/sell)

create table `transactions`
(
	transaction_id bigint primary key auto_increment,
	transaction_date datetime not null default current_timestamp,
	from_symbol varchar(255),
	to_symbol varchar(255),
	quantity float not null,
	sender bigint,	
	receiver bigint,
	foreign key (from_symbol) references `symbols`(symbol),
	foreign key (to_symbol) references `symbols`(symbol),
	foreign key (sender) references `accounts`(account_id),
	foreign key (receiver) references `accounts`(account_id)
);

-- Status types across all possible statuses (Completed, In Progress, Cancelled)

create table `order_status`
(
	status_id bigint primary key auto_increment,
	name varchar(255) not null unique,
	`type` varchar(255) not null unique
);

insert into order_status (name,`type`) values ('Completed','completed');
insert into order_status (name,`type`) values ('In Progress','in_progress');
insert into order_status (name,`type`) values ('Cancelled','cancelled');

-- Limit order, market order, etc

create table `order_type`
(
	order_type_id bigint primary key auto_increment,
	name varchar(255) not null unique,
	`type` varchar(255) not null unique
);

insert into order_type (name,`type`) values ('Market Order','market');
insert into order_type (name,`type`) values ('Limit Order','limit');

-- Orders an account has entered, regardless of state

create table `orders`
(
	order_id bigint primary key auto_increment,
	status_id bigint not null,
	order_type_id bigint not null,
	account_id bigint not null,
	limit_price float,				-- if the order is a limit order, there should be a limit order price (not for market orders)
	foreign key (status_id) references `order_status`(status_id),
	foreign key (order_type_id) references `order_type`(order_type_id),
	foreign key (account_id) references `accounts`(account_id)
);

commit;

