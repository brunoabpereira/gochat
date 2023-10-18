package com.example.gochat;

import java.util.Date;
import java.nio.charset.StandardCharsets;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.security.SecureRandom;
import java.util.Random;
import io.jsonwebtoken.*;

public class JwtUtil {
    // SECRET is encoded with base64
    private static final String SECRET = "26SrjQKKdr3Av2S04thIfsXcx4lSInVGjBYk5kUZrlSYFZfmGUZ9t9pcY8Rv8J2026SrjQKKdr3Av2S04thIfsXcx4lSInVGjBYk5kUZrlSYFZfmGUZ9t9pcY8Rv8J20";
    private static final long EXPIRATION_TIME = 600_000; // 10 mins

    public static String generateToken(String subject) {
        return Jwts.builder()
            .setSubject(subject)
            .setExpiration(new Date(System.currentTimeMillis() + EXPIRATION_TIME))
            .signWith(SignatureAlgorithm.HS512, SECRET)
            .compact();
    }

    public static String extractSubject(String token) {
        return Jwts.parser()
            .setSigningKey(SECRET)
			.build()
            .parseSignedClaims(token)
            .getBody()
            .getSubject();
    }

    public static String salt(){
	    Random secRandom = new SecureRandom();
		byte[] salt = new byte[32];
		secRandom.nextBytes(salt);
		return new String(salt, StandardCharsets.UTF_8);
	}

	public static String hash(String str, String salt){
		String hashstr = null;
		try {
			MessageDigest md = MessageDigest.getInstance("SHA-512");

			md.update(salt.getBytes(StandardCharsets.UTF_8));
			byte[] bytes = md.digest(str.getBytes(StandardCharsets.UTF_8));
			
			StringBuilder sb = new StringBuilder();
			for(int i=0; i< bytes.length ;i++){
				sb.append(Integer.toString((bytes[i] & 0xff) + 0x100, 16).substring(1));
			}
			hashstr = sb.toString();

		} catch (NoSuchAlgorithmException e) {
			e.printStackTrace();
		}
		return hashstr;
	}
}
