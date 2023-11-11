package com.example.auth;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.BeforeEach;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.autoconfigure.web.servlet.WebMvcTest;
import org.springframework.boot.test.mock.mockito.MockBean;
import org.springframework.test.web.servlet.MockMvc;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.http.MediaType;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;
import static org.mockito.Mockito.when;

@WebMvcTest( controllers = AuthController.class )
@AutoConfigureMockMvc
class AuthControllerTests {

	@Autowired
	private MockMvc mvc;
	@MockBean
    private UserRepository userRepo;
	@MockBean
	private AuthUtil auth;
	@Autowired
	private ObjectMapper objectMapper;

	@BeforeEach
	public void setUp() {
		String userPassword = "test";
		String userHash = "abc";
		User user = new User("test", userHash, "123", "test@example.com");
		when(userRepo.findUserByEmail(user.getUseremail())).thenReturn(user);
		when(auth.hash(userPassword, user.getUsersalt())).thenReturn(userHash);
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
