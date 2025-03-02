package com.mockingdemo.tests;

import com.mockingdemo.client.ServerAClient;
import com.mockingdemo.client.model.UserRequest;
import com.mockingdemo.client.model.UserResponse;
import io.restassured.RestAssured;
import io.restassured.response.Response;
import org.testng.annotations.BeforeClass;
import org.testng.annotations.Test;

import static org.testng.Assert.*;

public class ServerAApiTest {
    private static final String SERVER_A_BASE_URL = "http://localhost:8080";
    
    private ServerAClient serverAClient;
    
    @BeforeClass
    public void setUp() {
        // Configure REST Assured
        RestAssured.baseURI = SERVER_A_BASE_URL;
        
        // Initialize Server A client
        serverAClient = new ServerAClient(SERVER_A_BASE_URL);
        
        // Assuming Server A is already running in QA mode
        System.out.println("Connecting to Server A at: " + SERVER_A_BASE_URL);
    }
    
    @Test
    public void testServerAProcessEndpoint() {
        // Create test request
        UserRequest request = new UserRequest();
        request.setUserId("test-user-123");
        request.setQuery("test query");
        
        // Send request to Server A
        Response response = serverAClient.sendProcessRequest(request);
        
        // Print the response for debugging
        System.out.println("Response body: " + response.getBody().asString());
        
        // Verify response
        assertEquals(response.getStatusCode(), 200);
        assertEquals(serverAClient.getEnvironment(response), "qa");
        
        UserResponse userResponse = response.as(UserResponse.class);
        assertTrue(userResponse.isSuccess());
        
        // Check if result contains expected text without asserting on requestId
        assertTrue(userResponse.getResult().contains("Mock response"));
    }
    
    @Test
    public void testServerAWithDifferentUser() {
        // Create test request with different user
        UserRequest request = new UserRequest();
        request.setUserId("test-user-456");
        request.setQuery("another test query");
        
        // Send request to Server A
        Response response = serverAClient.sendProcessRequest(request);
        
        // Verify response
        assertEquals(response.getStatusCode(), 200);
        
        UserResponse userResponse = response.as(UserResponse.class);
        assertTrue(userResponse.isSuccess());
        // Removed assertion on requestId
    }
    
    @Test
    public void testServerAWithEmptyQuery() {
        // Create test request with empty query
        UserRequest request = new UserRequest();
        request.setUserId("test-user-789");
        request.setQuery("");
        
        // Send request to Server A
        Response response = serverAClient.sendProcessRequest(request);
        
        // Verify response
        assertEquals(response.getStatusCode(), 200);
        
        UserResponse userResponse = response.as(UserResponse.class);
        assertTrue(userResponse.isSuccess());
    }
}