package com.mockingdemo.client.model;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class EnhancedRequest {
    private String userId;
    private String query;
    private String timestamp;
    private String requestId;
    private String serverAVersion;
}