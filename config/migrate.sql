create table services
(
	id integer
		constraint services_pk
			primary key autoincrement,
	uuid varchar not null,
	name varchar not null,
	display varchar not null,
	instances text,
	grayscale text,
	load_mode varchar default 'random'
);

create table routes
(
	id integer
		constraint routes_pk
			primary key autoincrement,
	uuid varchar not null,
	pattern varchar,
	match_type varchar default 'regex',
	method varchar default 'ALL',
	service varchar,
	authorize int2 default 0,
	authorize_prefix varchar default 'AUTH',
	custom_token int2 default 0,
	rewrite varchar,
	description int
);

create table instances
(
	id integer
		constraint instances_pk
			primary key autoincrement,
	uuid varchar not null,
	url varchar not null
);

create table rewrites
(
	id integer
		constraint rewrites_pk
			primary key autoincrement,
	uuid varchar not null,
	origin varchar,
	rewrite varchar
);

