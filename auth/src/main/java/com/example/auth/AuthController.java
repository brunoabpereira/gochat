package com.example.auth;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.CookieValue;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.dao.DataIntegrityViolationException;
import jakarta.servlet.http.HttpServletResponse;
import jakarta.servlet.http.Cookie;

@RestController
@RequestMapping("/api")
public class AuthController {
	
	@Autowired
    private UserRepository userRepo;
	private static final String jwtCookieName = "JWTID";
	private static final Logger log = LoggerFactory.getLogger(AuthController.class);

	@PostMapping("/authorize")
	public void authorize(@RequestBody AuthorizeRequest authReq, HttpServletResponse response) {
		User user = userRepo.findUserByEmail(authReq.getUseremail());
		if (user == null) {
			response.setStatus(HttpStatus.BAD_REQUEST.value());
		}else if ( user.getUserhash().equals(JwtUtil.hash(authReq.getPassword(),user.getUsersalt())) ){
			String token = JwtUtil.generateToken(user.getUsername());
			Cookie cookie = new Cookie(jwtCookieName,token);
			cookie.setPath("/");
			cookie.setHttpOnly(true); // this cookie will be hidden from scripts on the client side
			cookie.setSecure(true);
			response.addCookie(cookie);
			response.setStatus(HttpStatus.OK.value());
		}else {
			response.setStatus(HttpStatus.UNAUTHORIZED.value());
		}
	}

	@PostMapping("/users")
	public ResponseEntity<ErrorResponse> registerUser(@RequestBody RegisterUserRequest regReq) {
		String passwordSalt = JwtUtil.salt();
		String passwordHash = JwtUtil.hash(regReq.getPassword(),passwordSalt);
		User user = new User(
			regReq.getUsername(),
			passwordHash,
			passwordSalt,
			regReq.getUseremail());
		try {
			userRepo.saveAndFlush(user);
			return new ResponseEntity<>(null, HttpStatus.OK);
		}catch(DataIntegrityViolationException e){
			log.error("test" + e.toString());
			ErrorResponse err = null;
			if (e.toString().contains("constraint_username")){
				err = new ErrorResponse("Username already used.");
			}else if (e.toString().contains("constraint_useremail")){
				err = new ErrorResponse("Email already used.");
			}
			return new ResponseEntity<>(err, HttpStatus.BAD_REQUEST);
		}
	}

	@DeleteMapping("/users")
	public void deleteUser(@CookieValue(jwtCookieName) String jwtCookie, @RequestBody DeleteRequest delReq, HttpServletResponse response) {
		if ( JwtUtil.extractSubject(jwtCookie).equals("admin") ){
			User user = userRepo.findUserByEmail(delReq.getUseremail());
			if ( user == null ){
				response.setStatus(HttpStatus.BAD_REQUEST.value());
			}else {
				userRepo.deleteById(user.getUserid().intValue());
				response.setStatus(HttpStatus.OK.value());
			}
		}else {
			response.setStatus(HttpStatus.UNAUTHORIZED.value());
		}
	}

}

class ErrorResponse{
	private String error;

	ErrorResponse(String error){
		this.error = error;
	}

	public String getError(){
		return this.error;
	}

	public void setError(String error){
		this.error = error;
	}
}

class DeleteRequest {
	private String useremail;

	public DeleteRequest(){}

	public DeleteRequest(String useremail){
		this.useremail = useremail;
	}

	public void setUseremail(String useremail) {
        this.useremail = useremail;
    }

	public String getUseremail() {
        return this.useremail;
    }
}

class AuthorizeRequest {
	private String useremail;
	private String password;

	public AuthorizeRequest(){}

	public AuthorizeRequest(String useremail, String password){
		this.useremail = useremail;
		this.password = password;
	}

	public void setUseremail(String useremail) {
        this.useremail = useremail;
    }

	public void setPassword(String password) {
        this.password = password;
    }

	public String getUseremail() {
        return this.useremail;
    }

	public String getPassword() {
        return this.password;
    }
}

class RegisterUserRequest {
	private String username;
	private String useremail;
	private String password;

	public RegisterUserRequest(){}

	public RegisterUserRequest(String username, String useremail, String password){
		this.username = username;
		this.useremail = useremail;
		this.password = password;
	}

	public void setUsername(String username) {
        this.username = username;
    }

	public void setUseremail(String useremail) {
        this.useremail = useremail;
    }

	public void setPassword(String password) {
        this.password = password;
    }

	public String getUsername() {
        return this.username;
    }

	public String getUseremail() {
        return this.useremail;
    }

	public String getPassword() {
        return this.password;
    }
}