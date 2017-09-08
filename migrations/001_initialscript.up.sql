CREATE TABLE "users"
("id" integer primary key autoincrement not null,
"name" text not null unique,
"password" text not null)