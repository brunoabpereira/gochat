package com.example.gochat;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.beans.factory.annotation.Autowired;
import java.util.List;

@RestController
public class AuthController {
	
	private static final Logger log = LoggerFactory.getLogger(AuthController.class);
	
	@Autowired
    private UserRepository userRepo;

	@PostMapping("/api/authorize")
	public String authorize(@RequestParam String username, @RequestParam String password) {
		List<User> listUsers = userRepo.findAll();
		for (User user : listUsers){
			log.info(user.getUsername());
		}
		log.info("authorize: " + username + " " + password);
		return "";
	}

	@PostMapping("/api/users")
	public void registerUser(@PathVariable Long id) {
		log.info("registerUser");
	}

	@DeleteMapping("/api/users/{id}")
	public void deleteUser(@PathVariable Long id) {
		log.info("deleteUser");
	}

}