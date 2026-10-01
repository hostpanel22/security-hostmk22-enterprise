package com.hostmk.sdk.models.fraud;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.time.Instant;
import java.util.List;

/** Result of the fraud analysis endpoint. */
@JsonIgnoreProperties(ignoreUnknown = true)
public record FraudAnalysis(
        @JsonProperty("risk_score") int riskScore,
        @JsonProperty("risk_level") String riskLevel,
        boolean blocked,
        List<String> matched,
        @JsonProperty("analyzed_at") Instant analyzedAt
) {
}
