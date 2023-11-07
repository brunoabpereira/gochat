CREATE OR REPLACE FUNCTION getUsersChannels(username_arg character varying(200))
RETURNS TABLE (channelid int, channelname character varying(200)) 
LANGUAGE plpgsql
AS 
$$
BEGIN
	RETURN QUERY
	SELECT c.channelid, c.channelname FROM users u 
	JOIN channelmembers cm ON u.userid = cm.userid
	JOIN channels c ON cm.channelid = c.channelid
	WHERE u.username = username_arg;
END
$$;

---

CREATE OR REPLACE PROCEDURE addUserToChannels(username_arg character varying(200), channelname_arg character varying(200))
LANGUAGE plpgsql
AS 
$$
BEGIN
	INSERT INTO channelmembers
	(userid,channelid)
	VALUES
	(
		(SELECT userid FROM users WHERE username = username_arg),
		(SELECT channelid FROM channels WHERE channelname = channelname_arg)
	);
END
$$;
