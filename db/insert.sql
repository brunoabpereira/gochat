INSERT INTO gochat.channels
(channelname)
VALUES
('general'),
('test');

INSERT INTO gochat.channelmembers
(userid,channelid)
VALUES
-- Add to general channel
(1,3), -- admin
(2,3), -- chris
(3,3), -- paulie
(5,3), -- tony
-- Add to test channel
(2,4), -- chris
(5,4); -- tony

---

DROP TABLE gochat.channelmembers;
DROP TABLE gochat.messages;
DROP TABLE gochat.channels;
DROP TABLE gochat.users;
