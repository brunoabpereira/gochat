package com.example.auth.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;

import com.example.auth.model.User;

public interface UserRepository extends JpaRepository<User, Integer> {
    @Query(
        value = "SELECT * FROM users u WHERE u.useremail LIKE :email", 
        nativeQuery = true
    )
    User findUserByEmail(@Param("email") String email);
}