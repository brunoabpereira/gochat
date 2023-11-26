package com.example.auth.model;

public class AuthorizeRequest {
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
