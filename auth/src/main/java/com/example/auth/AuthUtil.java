package com.example.auth;

import java.util.Date;
import java.nio.charset.StandardCharsets;
import java.security.KeyFactory;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.security.SecureRandom;
import java.security.PublicKey;
import java.security.PrivateKey;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.security.spec.*;
import java.util.Random;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;
import io.jsonwebtoken.Jwts;

@Component
public class AuthUtil {
    private static final long EXPIRATION_TIME = 1800_000; // 30 mins
	private PublicKey pubKey;
	private PrivateKey privKey;

	@Autowired
	public void setPubKey(@Value("${pubKeyFilename}") String filename) throws Exception{
		byte[] keyBytes = Files.readAllBytes(Paths.get(filename));
		X509EncodedKeySpec spec = new X509EncodedKeySpec(keyBytes);
		KeyFactory kf = KeyFactory.getInstance("RSA");
		pubKey = kf.generatePublic(spec);
	}

	@Autowired
	public void setPrivKey(@Value("${privKeyFilename}") String filename) throws  Exception{
		byte[] keyBytes = Files.readAllBytes(Paths.get(filename));
		PKCS8EncodedKeySpec spec = new PKCS8EncodedKeySpec(keyBytes);
		KeyFactory kf = KeyFactory.getInstance("RSA");
		privKey = kf.generatePrivate(spec);
	}

    public String generateToken(String subject) {
        return Jwts.builder()
            .subject(subject)
            .expiration(new Date(System.currentTimeMillis() + EXPIRATION_TIME))
            .signWith(privKey)
            .compact();
    }

    public String extractSubject(String token) {
        return Jwts.parser()
			.verifyWith(pubKey)
			.build()
            .parseSignedClaims(token)
            .getPayload()
            .getSubject();
    }

    public String salt(){
	    Random secRandom = new SecureRandom();
		byte[] salt = new byte[32];
		secRandom.nextBytes(salt);
		StringBuilder sb = new StringBuilder();
		for(int i=0; i< salt.length ;i++){
			sb.append(Integer.toString((salt[i] & 0xff) + 0x100, 16).substring(1));
		}
		return sb.toString();
	}

	public String hash(String str, String salt){
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
