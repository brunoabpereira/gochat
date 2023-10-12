INSERT INTO gochat.users
(username,userhash,usersalt,useremail)
VALUES
('admin','','','admin@example.com'),
('chris','','','chris@example.com'),
('paulie','','','paulie@example.com'),
('tony','','','tony@example.com');

INSERT INTO gochat.channels
(channelname)
VALUES
('general'),
('test');

INSERT INTO gochat.channelmembers
(userid,channelid)
VALUES
-- Add to general channel
(1,1), -- admin
(2,1), -- chris
(3,1), -- paulie
(4,1), -- tony
-- Add to test channel
(2,2), -- chris
(4,2); -- tony

---

DROP TABLE gochat.channelsmembers;
DROP TABLE gochat.messages;
DROP TABLE gochat.channels;
DROP TABLE gochat.users;
