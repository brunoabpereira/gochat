package com.example.auth.exception;

public class UserWithEmailNotFoundException extends Exception{

    public UserWithEmailNotFoundException(String userEmail){
        super("User with email '" + userEmail + "' not found.");
    }
}

