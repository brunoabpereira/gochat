INSERT INTO gochat.users
(userid,username,userhash)
VALUES
(0,'admin','ABC'),
(1,'chris','ABC'),
(2,'paulie','ABC'),
(3,'tony','ABC');


INSERT INTO gochat.channels
(channelid,channelname)
VALUES
(0,'general'),
(1,'test');


INSERT INTO gochat.channelsmembers
(userid,channelid)
VALUES
-- Add to general channel
(0,0), -- admin
(1,0), -- chris
(2,0), -- paulie
(3,0), -- tony
-- Add to test channel
(1,1), -- chris
(3,1); -- tony


