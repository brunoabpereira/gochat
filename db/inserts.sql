\COPY users (username,userhash,usersalt,useremail) FROM 'users.csv' WITH (FORMAT CSV, HEADER, DELIMITER ',');

INSERT INTO gochat.channels
(channelname)
VALUES
('General'),
('Bada Bing!');

call addUserToChannels('tony', 'General');
call addUserToChannels('paulie', 'General');
call addUserToChannels('chrissy', 'General');
call addUserToChannels('bobby', 'General');
call addUserToChannels('silvio', 'General');

call addUserToChannels('tony', 'Bada Bing!');
call addUserToChannels('paulie', 'Bada Bing!');
call addUserToChannels('bobby', 'Bada Bing!');
