package com.hostmk.sdk.services;

import com.hostmk.sdk.http.HostMKHttpClient;
import com.hostmk.sdk.models.ListResult;
import com.hostmk.sdk.models.fraud.FraudAnalysis;
import com.hostmk.sdk.models.fraud.FraudRule;
import com.hostmk.sdk.models.fraud.FraudViolation;

import java.util.Map;

/** Fraud rules, analysis, and violation listing. */
public final class FraudService {
    private final HostMKHttpClient client;

    public FraudService(HostMKHttpClient client) {
        this.client = client;
    }

    public ListResult<FraudRule> listRules() {
        return client.list("/api/v1/fraud/rules", Map.of(), "rules", FraudRule.class);
    }

    public FraudRule createRule(FraudRule rule) {
        return client.post("/api/v1/fraud/rules", Map.of(), rule, FraudRule.class);
    }

    public FraudRule updateRule(FraudRule rule) {
        return client.put("/api/v1/fraud/rules", Map.of(), rule, FraudRule.class);
    }

    public void deleteRule(String id) {
        client.delete("/api/v1/fraud/rules", Map.of("id", id));
    }

    public FraudAnalysis analyze(Map<String, Object> attributes) {
        return client.post("/api/v1/fraud/analyze", Map.of(), attributes, FraudAnalysis.class);
    }

    public ListResult<FraudViolation> listViolations() {
        return client.list("/api/v1/fraud/violations", Map.of(), "events", FraudViolation.class);
    }
}
