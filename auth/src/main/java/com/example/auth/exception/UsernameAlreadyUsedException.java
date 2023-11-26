package com.example.auth.exception;

public class UsernameAlreadyUsedException extends Exception{

    public UsernameAlreadyUsedException(String userName){
        super("User name '" + userName + "' already used.");
    }
}

