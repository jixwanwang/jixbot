drop table if exists channels;
create table channels (username text primary key);

drop table if exists viewers;
create table viewers (id serial primary key, username text, channel text);

drop table if exists brawlwins;
create table brawlwins (id serial primary key, season int, viewer_id int, wins int);

drop table if exists counts;
create table counts (id serial primary key, type text, viewer_id int, count int);

drop table if exists better_counts;
create table better_counts (viewer_id int primary key, money int default 0, lines_typed int default 0, time_spent int default 0);
create index better_counts_money_idx on better_counts (money);
create index better_counts_lines_typed_idx on better_counts (lines_typed);
create index better_counts_time_spent_idx on better_counts (time_spent);

drop table if exists emotes;
create table emotes (id serial primary key, channel text, emote text);

drop table if exists textcommands;
create table textcommands (id serial primary key, channel text, command text, message text, clearance int, cooldown int);

drop table if exists commands;
create table commands (id serial primary key, channel text, command text);

drop table if exists channel_properties;
create table channel_properties (id serial primary key, channel text, k text, v text);

drop table if exists questions;
create table questions (id serial primary key, channel text, question text, answer text);

drop table if exists quotes;
create table quotes (id serial primary key, channel text, quote_type text, rank int, quote text);
create index quotes_channel_quote_type on quotes(channel, quote_type);
create index quotes_rank on quotes(rank);

insert into channels (username) values ('jixwanwang');

insert into viewers (username, channel) values ('jixwanwang', 'jixwanwang');

insert into brawlwins (season, viewer_id, wins) values (1, 1, 2);

insert into textcommands (channel, command, message, clearance) values ('jixwanwang', '!jix', 'Staff Dansgame', 0);
insert into textcommands (channel, command, message, clearance) values ('_global', '!jixbot', 'I am created!', 1);

insert into commands (channel, command) values ('jixwanwang', 'brawl');

insert into channel_properties (channel, k, v) values ('jixwanwang', 'currency', 'JixCoin');