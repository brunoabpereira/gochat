CREATE SCHEMA gochat;

CREATE TABLE gochat.users
(
    userid SERIAL PRIMARY KEY,
    username character varying(50) NOT NULL,
    userhash character varying(200) NOT NULL,
    usersalt character varying(200) NOT NULL,
    useremail character varying(200) NOT NULL,
    CONSTRAINT constraint_username UNIQUE (username),
    CONSTRAINT constraint_useremail UNIQUE (useremail)
);

CREATE UNIQUE INDEX username_idx ON gochat.users (username);
CREATE UNIQUE INDEX useremail_idx ON gochat.users (useremail);

CREATE TABLE gochat.channels
(
    channelid SERIAL PRIMARY KEY,
    channelname character varying(200) NOT NULL
);

CREATE TABLE gochat.channelmembers
(
    userid integer NOT NULL,
    channelid integer NOT NULL,
    CONSTRAINT channelsmembers_pk PRIMARY KEY (userid, channelid),
    CONSTRAINT fk_users_userid FOREIGN KEY (userid) REFERENCES gochat.users(userid) ON DELETE CASCADE,
    CONSTRAINT fk_channels_channelid FOREIGN KEY (channelid) REFERENCES gochat.channels(channelid) ON DELETE CASCADE
);

CREATE TABLE gochat.messages
(
    messageid SERIAL PRIMARY KEY,
    messagetime timestamp NOT NULL,
    messagetext character varying(200) NOT NULL,
    userid integer NOT NULL,
    channelid integer NOT NULL,
    CONSTRAINT fk_users_userid FOREIGN KEY (userid) REFERENCES gochat.users(userid),
    CONSTRAINT fk_channels_channelid FOREIGN KEY (channelid) REFERENCES gochat.channels(channelid)
);
