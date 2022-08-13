package main

var databaseInfoColumnsV1 = [...]columnDesc{
	{"schema_ordinal": "integer"},
}

var registeredUserColumnsV1 = [...]columnDesc{
	{"creation_time": "timestamp without time zone"},
	{"email": "character varying"},
	{"email_validated": "boolean"},
	{"password_cost": "integer"},
	{"password_hash": "bytea"},
	{"rights": "text"},
	{"username": "character varying"},
	{"version": "integer"},
}

var pendingValidationColumnsV1 = [...]columnDesc{
	{"username": "character varying"},
	{"uuid": "uuid"},
}

var schemaForV1 = schema{
	"jutzo_database_info":      databaseInfoColumnsV1[0:],
	"jutzo_registered_user":    registeredUserColumnsV1[0:],
	"jutzo_pending_validation": pendingValidationColumnsV1[0:],
}

// DDL for creating the current schema if it doesn't exist
var ddlV1 = [...]string{
	`drop table if exists jutzo_database_info cascade `,
	`drop table if exists jutzo_pending_validation cascade `,
	`drop table if exists jutzo_registered_user cascade `,
	`create table if not exists jutzo_registered_user
		(
		username        varchar(256)            not null
			constraint username_key
			primary key,
		email           varchar(256)            not null,
		email_validated boolean   default false not null,
		version         integer   default 1     not null,
		creation_time   timestamp default now() not null,
		password_cost   integer                 not null,
		password_hash   bytea                   not null,
		rights          text      default 'blog,login'::text
	)`,
	`comment on column jutzo_registered_user.username is 'The user provided username.'`,
	`comment on column jutzo_registered_user.rights is 'The comma delimited rights for the user'`,
	`alter table jutzo_registered_user owner to jutzo`,
	`create table if not exists jutzo_pending_validation
		(
		uuid     uuid default gen_random_uuid() not null
			constraint uuid_key
			primary key,
		username varchar(256)                   not null
			constraint foreign_key_name
			references jutzo_registered_user
			on update cascade on delete cascade
		)`,
	`alter table jutzo_pending_validation owner to jutzo`,
	`create unique index if not exists email_idx on jutzo_registered_user (email)`,
	`create table jutzo_database_info
		(
			schema_ordinal integer default 1
		)`,
	`comment on table jutzo_database_info is 'Information about the database'`,
	`alter table jutzo_database_info owner to jutzo`,
	`insert into jutzo_database_info (schema_ordinal) values (0)`,
}

// Define the V1 schema ordinal
const schemaV1 = 0
