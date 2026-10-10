package io.rune.sdk;

import java.io.IOException;
import java.net.URI;
import java.net.URLEncoder;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.List;

/**
 * Minimal HTTP client for the Rune secrets API.
 * Uses java.net.http (Java 11+) with zero extra dependencies.
 */
public class RuneClient {

    private final String baseUrl;
    private final String token;
    private final HttpClient http;

    public RuneClient(String baseUrl, String token) {
        this.baseUrl = baseUrl.endsWith("/") ? baseUrl.substring(0, baseUrl.length() - 1) : baseUrl;
        this.token = token;
        this.http = HttpClient.newBuilder().build();
    }

    /**
     * Lists all secret keys available in the given namespace.
     * Calls GET /secret/list?namespace=<ns>
     */
    public List<String> listKeys(String namespace) {
        String url = baseUrl + "/secret/list?namespace=" + encode(namespace);
        String body = get(url);
        return parseJsonStringArray(body);
    }

    /**
     * Retrieves the value of a secret by key in the given namespace.
     * Calls GET /secret/get?namespace=<ns>&key=<key>
     */
    public String getSecret(String namespace, String key) {
        String url = baseUrl + "/secret/get?namespace=" + encode(namespace) + "&key=" + encode(key);
        String body = get(url);
        return extractJsonStringValue(body, "value");
    }

    private String get(String url) {
        HttpRequest req = HttpRequest.newBuilder()
                .uri(URI.create(url))
                .header("Authorization", "Bearer " + token)
                .GET()
                .build();
        try {
            HttpResponse<String> resp = http.send(req, HttpResponse.BodyHandlers.ofString());
            if (resp.statusCode() != 200) {
                throw new RuneException("Rune server returned " + resp.statusCode() + ": " + resp.body().trim());
            }
            return resp.body();
        } catch (IOException | InterruptedException e) {
            throw new RuneException("Failed to reach Rune server at " + url, e);
        }
    }

    // Minimal JSON array parser for ["key1","key2"] — avoids Jackson dependency.
    private List<String> parseJsonStringArray(String json) {
        List<String> result = new ArrayList<>();
        if (json == null || json.isBlank()) return result;
        String trimmed = json.trim();
        if (!trimmed.startsWith("[")) return result;
        // Strip [ ] and split on ","
        trimmed = trimmed.substring(1, trimmed.length() - 1).trim();
        if (trimmed.isEmpty()) return result;
        for (String part : trimmed.split(",")) {
            String val = part.trim();
            if (val.startsWith("\"") && val.endsWith("\"")) {
                result.add(val.substring(1, val.length() - 1));
            }
        }
        return result;
    }

    // Minimal {"value":"..."} extractor — avoids Jackson dependency.
    private String extractJsonStringValue(String json, String fieldName) {
        if (json == null) return null;
        String search = "\"" + fieldName + "\":\"";
        int start = json.indexOf(search);
        if (start == -1) return null;
        start += search.length();
        int end = json.indexOf("\"", start);
        if (end == -1) return null;
        // Unescape basic JSON escape sequences
        return json.substring(start, end).replace("\\n", "\n").replace("\\\"", "\"").replace("\\\\", "\\");
    }

    private String encode(String value) {
        return URLEncoder.encode(value, StandardCharsets.UTF_8);
    }
}
