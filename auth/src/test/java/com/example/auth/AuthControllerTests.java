package com.example.auth;

import org.junit.jupiter.api.Test;
import org.mockito.Mockito;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.mock.mockito.MockBean;
import org.springframework.test.web.servlet.MockMvc;
import org.springframework.http.MediaType;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;
import org.aspectj.lang.annotation.Before;


@SpringBootTest
@AutoConfigureMockMvc
class AuthControllerTests {

	@Autowired
	private MockMvc mvc;

	@MockBean
    private UserRepository userRepo;

	// @Before
	// public void setUp() {
	// 	User user1 = new User("test", "123", "123", "test@example.com");
	// 	Mockito.when(userRepo.findUserByEmail(user1.getUseremail())).thenReturn(user1);
	// }

	@Test
	void contextLoads() {
	}

	@Test
	void testAuthorize() throws Exception {
		this.mvc.perform(post("/api/authorize")
		.contentType(MediaType.APPLICATION_JSON))
		.andExpect(status().isOk());
	}

}
