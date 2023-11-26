package com.example.auth.service;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.dao.DataIntegrityViolationException;
import org.springframework.stereotype.Service;

import com.example.auth.model.User;
import com.example.auth.repository.UserRepository;
import com.example.auth.util.AuthUtil;
import com.example.auth.exception.IncorrectUserPasswordException;
import com.example.auth.exception.UserWithEmailNotFoundException;
import com.example.auth.exception.UsernameAlreadyUsedException;
import com.example.auth.exception.UseremailAlreadyUsedException;

@Service
public class AuthService {

    @Autowired
    private UserRepository userRepo;
	@Autowired
	private AuthUtil auth;
    @Value("${jwtAdminUsername}")
    private final String jwtAdminUsername = "admin";

    public AuthService(){

    }

    public String authorize(String userEmail, String password) throws UserWithEmailNotFoundException, IncorrectUserPasswordException {
        User user = userRepo.findUserByEmail(userEmail);

        if (user == null) {
            throw new UserWithEmailNotFoundException(userEmail);
        }
        
        if (!user.getUserhash().equals(auth.hash(password, user.getUsersalt()))) {
            throw new IncorrectUserPasswordException(userEmail);
        }

        return auth.generateToken(user.getUsername());
    }

    public void registerUser(String userName, String userEmail, String password) throws UsernameAlreadyUsedException, UseremailAlreadyUsedException{
        String passwordSalt = auth.salt();
		String passwordHash = auth.hash(password, passwordSalt);
		User user = new User(userName, passwordHash, passwordSalt, userEmail);
		try {
			userRepo.saveAndFlush(user);
		}catch(DataIntegrityViolationException exception){
			if (exception.toString().contains("constraint_username")){
                throw new UsernameAlreadyUsedException(userName);
			}else if (exception.toString().contains("constraint_useremail")){
                throw new UseremailAlreadyUsedException(userEmail);
			}else {
                throw exception;
            }
		}
    }

    public void deleteUser(String userEmail) throws UserWithEmailNotFoundException{
        User user = userRepo.findUserByEmail(userEmail);
        if ( user == null ){
            throw new UserWithEmailNotFoundException(userEmail);
        }else {
            userRepo.deleteById(user.getUserid().intValue());
        }
    }

    public boolean tokenSubjectIsAdmin(String jwtToken){
        return auth.extractSubject(jwtToken).equals(jwtAdminUsername);
    }

}
