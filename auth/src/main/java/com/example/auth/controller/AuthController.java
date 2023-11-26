package com.example.auth.controller;

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
import org.springframework.beans.factory.annotation.Value;
import jakarta.servlet.http.HttpServletResponse;
import jakarta.servlet.http.Cookie;

import com.example.auth.service.AuthService;
import com.example.auth.exception.IncorrectUserPasswordException;
import com.example.auth.exception.UserWithEmailNotFoundException;
import com.example.auth.exception.UseremailAlreadyUsedException;
import com.example.auth.exception.UsernameAlreadyUsedException;
import com.example.auth.model.AuthorizeRequest;
import com.example.auth.model.RegisterUserRequest;
import com.example.auth.model.ErrorResponse;
import com.example.auth.model.DeleteRequest;

@RestController
@RequestMapping("/api")
public class AuthController {

	private static final Logger log = LoggerFactory.getLogger(AuthController.class);

	@Value("${jwtCookieName}")
	private final String jwtCookieName = "JWTID";
	@Autowired
	private AuthService authService;

	@PostMapping("/authorize")
	public void authorize(@RequestBody AuthorizeRequest authReq, HttpServletResponse response) {
		String token = null;
		try{
			token = authService.authorize(authReq.getUseremail(), authReq.getPassword());
		} catch (UserWithEmailNotFoundException exception) {
			log.error(exception.toString());
			response.setStatus(HttpStatus.BAD_REQUEST.value());
		} catch (IncorrectUserPasswordException exception) {
			log.error(exception.toString());
			response.setStatus(HttpStatus.UNAUTHORIZED.value());
		}
		if (token != null){
			Cookie cookie = new Cookie(jwtCookieName,token);
			cookie.setPath("/");
			cookie.setHttpOnly(true); // cookie will be hidden from scripts on the client side
			cookie.setSecure(true);
			response.addCookie(cookie);
			response.setStatus(HttpStatus.OK.value());
		}
	}

	@PostMapping("/users")
	public ResponseEntity<ErrorResponse> registerUser(@RequestBody RegisterUserRequest regReq) {
		try {
			authService.registerUser(regReq.getUsername(), regReq.getUseremail(), regReq.getPassword());
		}catch(UsernameAlreadyUsedException exception){
			log.error(exception.toString());
			ErrorResponse err = new ErrorResponse("Username already used.");
			return new ResponseEntity<>(err, HttpStatus.BAD_REQUEST);
		} catch (UseremailAlreadyUsedException exception){
			log.error(exception.toString());
			ErrorResponse err = new ErrorResponse("Email already used.");
			return new ResponseEntity<>(err, HttpStatus.BAD_REQUEST);
		}
		return new ResponseEntity<>(null, HttpStatus.OK);
	}

	@DeleteMapping("/users")
	public void deleteUser(@CookieValue(jwtCookieName) String jwtCookie, @RequestBody DeleteRequest delReq, HttpServletResponse response) {
		if ( authService.tokenSubjectIsAdmin(jwtCookie) ){
			try {
				authService.deleteUser(delReq.getUseremail());
				response.setStatus(HttpStatus.OK.value());
			} catch (UserWithEmailNotFoundException exception) {
				response.setStatus(HttpStatus.BAD_REQUEST.value());
			}
		}else {
			response.setStatus(HttpStatus.UNAUTHORIZED.value());
		}
	}

}
