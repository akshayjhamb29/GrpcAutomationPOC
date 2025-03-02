package com.mockingdemo.client.model;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class ProcessedResponse {
    private String result;
    private boolean success;
    private String processingTime;
    private String serverBId;
}