package com.example.auth.model;

public class DeleteRequest {
	private String useremail;

	public DeleteRequest(){}

	public DeleteRequest(String useremail){
		this.useremail = useremail;
	}

	public void setUseremail(String useremail) {
        this.useremail = useremail;
    }

	public String getUseremail() {
        return this.useremail;
    }
}