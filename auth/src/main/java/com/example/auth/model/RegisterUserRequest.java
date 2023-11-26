package com.example.auth.model;

public class RegisterUserRequest {
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
