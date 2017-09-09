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

CREATE TABLE "markers"
(
"id" integer primary key autoincrement not null,
"lat" real not null,
"lon" real not null,
"userid" integer not null,
"text" string not null,
"timeCreated" datetime not null,
foreign key(userid) references users(id)
)