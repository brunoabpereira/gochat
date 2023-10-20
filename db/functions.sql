CREATE OR REPLACE FUNCTION getUsersChannels(username_arg character varying(200))
RETURNS TABLE (channelid int, channelname character varying(200)) 
LANGUAGE plpgsql
AS $$
BEGIN
	RETURN QUERY
	SELECT c.channelid, c.channelname FROM users u 
	JOIN channelmembers cm ON u.userid = cm.userid
	JOIN channels c ON cm.channelid = c.channelid
	WHERE u.username = username_arg;
END;$$