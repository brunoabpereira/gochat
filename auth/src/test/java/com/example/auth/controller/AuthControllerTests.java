package com.example.auth.controller;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.BeforeEach;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.autoconfigure.web.servlet.WebMvcTest;
import org.springframework.boot.test.mock.mockito.MockBean;
import org.springframework.test.web.servlet.MockMvc;

import com.example.auth.exception.IncorrectUserPasswordException;
import com.example.auth.exception.UserWithEmailNotFoundException;
import com.example.auth.model.AuthorizeRequest;
import com.example.auth.service.AuthService;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.http.MediaType;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.AdditionalMatchers.not;
import static org.mockito.Mockito.when;

@WebMvcTest( controllers = AuthController.class )
@AutoConfigureMockMvc
class AuthControllerTests {

	@Autowired
	private MockMvc mvc;
	@Autowired
	private ObjectMapper objectMapper;
	@MockBean
    private AuthService authService;

	@BeforeEach
	public void setUp() throws Exception {
		String userEmail = "test@example.com";
		String userPassword = "test";
		when(authService.authorize(userEmail, userPassword)).thenReturn("token");
		when(authService.authorize(eq(userEmail), not(eq(userPassword)))).thenThrow(new IncorrectUserPasswordException(""));
		when(authService.authorize(not(eq(userEmail)), anyString())).thenThrow(new UserWithEmailNotFoundException(""));
	}

	@Test
	void contextLoads() {
	}

	@Test
	void authorizeWithCorrectPassword() throws Exception {
		String authReq = objectMapper.writeValueAsString(new AuthorizeRequest("test@example.com", "test"));
		this.mvc.perform(post("/api/authorize")
		.content(authReq)
		.contentType(MediaType.APPLICATION_JSON))
		.andExpect(status().is(200));
	}

	@Test
	void authorizeWithIncorrectPassword() throws Exception {
		String authReq = objectMapper.writeValueAsString(new AuthorizeRequest("test@example.com", "blabla"));
		this.mvc.perform(post("/api/authorize")
		.content(authReq)
		.contentType(MediaType.APPLICATION_JSON))
		.andExpect(status().is(401));
	}

	@Test
	void authorizeWithIncorrectEmail() throws Exception {
		String authReq = objectMapper.writeValueAsString(new AuthorizeRequest("incorrect@example.com", "123"));
		this.mvc.perform(post("/api/authorize")
		.content(authReq)
		.contentType(MediaType.APPLICATION_JSON))
		.andExpect(status().is(400));
	}

}
