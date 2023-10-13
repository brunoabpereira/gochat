package com.example.gochat;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.dao.DataIntegrityViolationException;

@RestController
public class AuthController {
	
	@Autowired
    private UserRepository userRepo;

	private static final Logger log = LoggerFactory.getLogger(AuthController.class);

	@PostMapping("/api/authorize")
	public ResponseEntity<String> authorize(@RequestBody AuthorizeRequest authReq) {
		User user = userRepo.findUserByEmail(authReq.getUseremail());
		if ( user.getUserhash().equals(JwtUtil.hash(authReq.getPassword(),user.getUsersalt())) ){
			String token = JwtUtil.generateToken(authReq.getUseremail());
        	return new ResponseEntity<>(token, HttpStatus.OK);
		}
		return new ResponseEntity<>(null, HttpStatus.UNAUTHORIZED);
	}

	@PostMapping("/api/users")
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
			return new ResponseEntity<>(null,HttpStatus.OK);
		}catch(DataIntegrityViolationException e){
			log.error("test" + e.toString());
			ErrorResponse err = null;
			if (e.toString().contains("constraint_username")){
				err = new ErrorResponse("Username already used.");
			}else if (e.toString().contains("constraint_useremail")){
				err = new ErrorResponse("Email already used.");
			}
			return new ResponseEntity<>(err,HttpStatus.OK);
		}
	}

	@DeleteMapping("/api/users/{id}")
	public void deleteUser(@PathVariable Long id) {
		log.info("deleteUser");
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

class AuthorizeRequest {
	private String useremail;
	private String password;

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