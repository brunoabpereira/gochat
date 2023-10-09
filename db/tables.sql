CREATE SCHEMA gochat;

CREATE TABLE gochat.users
(
    userid integer NOT NULL,
    username character varying(50) NOT NULL,
    userhash character varying(200) NOT NULL,
    CONSTRAINT users_pk PRIMARY KEY (userid)
);


CREATE TABLE gochat.channels
(
    channelid integer NOT NULL,
    channelname character varying(200) NOT NULL,
    CONSTRAINT channel_pk PRIMARY KEY (channelid)
);

CREATE TABLE gochat.channelsmembers
(
    userid integer NOT NULL,
    channelid integer NOT NULL,
    CONSTRAINT channelsmembers_pk PRIMARY KEY (userid, channelid),
    CONSTRAINT fk_users_userid FOREIGN KEY (userid) REFERENCES gochat.users(userid),
    CONSTRAINT fk_channels_channelid FOREIGN KEY (channelid) REFERENCES gochat.channels(channelid)
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
