package com.example.auth.model;

import jakarta.persistence.Entity;
import jakarta.persistence.Table;
import jakarta.persistence.Id;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;

@Entity
@Table(name = "users", schema="gochat")
public class User {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long userid;
    private String username;
    private String userhash;
    private String usersalt;
    private String useremail;

    public User() {}

    public User(String username, String userhash, String usersalt, String useremail) {
        this.username = username;
        this.userhash = userhash;
        this.usersalt = usersalt;
        this.useremail = useremail;
    }

    public void setUsername(String username) {
        this.username = username;
    }

    public void setUserhash(String userhash) {
        this.userhash = userhash;
    }

    public void setUsersalt(String usersalt) {
        this.usersalt = usersalt;
    }

    public void setUseremail(String useremail) {
        this.useremail = useremail;
    }

    public String getUsername() {
        return this.username;
    }

    public String getUserhash() {
        return this.userhash;
    }

    public String getUsersalt() {
        return this.usersalt;
    }

    public String getUseremail() {
        return this.useremail;
    }

    public Long getUserid() {
        return this.userid;
    }
}