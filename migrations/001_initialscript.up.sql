CREATE TABLE "users"
("id" integer primary key autoincrement not null,
"name" text not null unique,
"password" text not null);

CREATE TABLE "sessions"
("id" integer primary key autoincrement not null,
"sessionid" text not null unique,
"userid" integer not null,
"timeCreated" datetime not null,
foreign key(userid) references users(id)
);