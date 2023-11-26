package com.example.auth.exception;

public class UseremailAlreadyUsedException extends Exception{

    public UseremailAlreadyUsedException(String userEmail){
        super("User email '" + userEmail + "' already used.");
    }
}

