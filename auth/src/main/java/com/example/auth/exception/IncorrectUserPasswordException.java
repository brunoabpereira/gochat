package com.example.auth.exception;

public class IncorrectUserPasswordException extends Exception{

    public IncorrectUserPasswordException(String userEmail){
        super("Incorrect password for user email'" + userEmail);
    }
}
