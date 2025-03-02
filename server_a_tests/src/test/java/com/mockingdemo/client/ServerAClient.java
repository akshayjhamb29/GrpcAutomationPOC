package com.mockingdemo.client;

import com.mockingdemo.client.model.UserRequest;
import com.mockingdemo.client.model.UserResponse;
import io.restassured.RestAssured;
import io.restassured.http.ContentType;
import io.restassured.response.Response;

public class ServerAClient {
    private final String baseUrl;
    
    public ServerAClient(String baseUrl) {
        this.baseUrl = baseUrl;
    }
    
    public Response sendProcessRequest(UserRequest request) {
        return RestAssured.given()
                .contentType(ContentType.JSON)
                .body(request)
                .when()
                .post(baseUrl + "/api/process");
    }
    
    public UserResponse sendProcessRequestAndGetResponse(UserRequest request) {
        return sendProcessRequest(request)
                .then()
                .statusCode(200)
                .extract()
                .as(UserResponse.class);
    }
    
    public String getEnvironment(Response response) {
        return response.getHeader("X-Environment");
    }
}