drop table if exists 3manst;
create table 3manst (
	id bigint auto_increment primary key, 
	board binary(144) not null, 
	moats bit(3) not null,
	movesnext tinyint not null, 
	castling bit(6) not null, 
	enpassant binary(4) not null,
	halfmoveclock tinyint not null, 
	fullmovenumber smallint not null, 
	alive bit(3) not null,
	unique key everything(
		board, moats, movesnext, castling,
		enpassant,
		halfmoveclock, fullmovenumber,
		alive
	)
) ENGINE = InnoDB;

drop table if exists 3manplayer;
create table 3manplayer (
	id bigint auto_increment primary key,
	auth varbinary(100) not null
	-- name varchar(100) not null,
) engine = InnoDB default charset=utf8;

drop table if exists chessuser;
create table chessuser (
	id bigint auto_increment primary key,
	login varchar(20) unique key,
	passwd varchar(100) not null,
	name varchar(100),
	player bigint not null unique key,
	constraint
		foreign key (player) references 3manplayer (id)
		on update restrict
) engine = InnoDB default charset=utf8;

drop table if exists chessbot;
create table chessbot (
	id bigint auto_increment primary key,
	whoami varbinary(20) not null, -- ai type identifier
	owner bigint not null,
	ownname varchar(50),
	player bigint not null unique key,
	settings varbinary(500),
	unique key everything ( whoami, owner, settings ),
	constraint
		foreign key (owner) references chessuser (id)
		on update restrict,
	constraint
		foreign key (player) references 3manplayer (id)
		on update restrict
) engine = InnoDB default charset=utf8;

drop table if exists 3mangp;
create table 3mangp (
	id bigint auto_increment primary key, 
	lastmove bigint not null, 
	white bigint not null, 
	gray bigint not null, 
	black bigint not null, 
	created datetime not null,
	constraint
		foreign key (white) references 3manplayer (id)
		on update restrict,
	constraint
		foreign key (gray) references 3manplayer (id)
		on update restrict,
	constraint
		foreign key (black) references 3manplayer (id)
		on update restrict
) ENGINE = InnoDB;

drop table if exists 3manmv;
create table 3manmv (
	id bigint auto_increment primary key,
	fromto binary(4) not null,
	beforegame bigint not null,
	afterstate bigint not null,
	promotion tinyint,
	who bigint not null,
	constraint
		foreign key (beforegame) references 3mangp (id)
		on update restrict,
	constraint
		foreign key (who) references 3manplayer (id)
		on update restrict,
	constraint
		foreign key (afterstate) references 3manst (id)
		on update restrict,
	unique onemove(fromto, beforegame, promotion, who)
) engine = InnoDB;

alter table 3mangp
	add constraint
		foreign key (lastmove) references 3manmv (id)
		on update restrict;


-- vi:ft=mysql
