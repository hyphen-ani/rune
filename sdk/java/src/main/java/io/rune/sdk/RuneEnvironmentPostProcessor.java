package io.rune.sdk;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.env.EnvironmentPostProcessor;
import org.springframework.core.env.ConfigurableEnvironment;

import java.util.logging.Logger;

/**
 * Runs before any @Bean methods, after application.yml has been loaded.
 * Reads rune.url / rune.token / rune.namespace from the environment,
 * fetches all secrets, and inserts them as the highest-precedence PropertySource
 * so ${rune.<key>} expressions resolve correctly.
 *
 * Registered via META-INF/spring/org.springframework.boot.env.EnvironmentPostProcessor.imports
 */
public class RuneEnvironmentPostProcessor implements EnvironmentPostProcessor {

    private static final Logger log = Logger.getLogger(RuneEnvironmentPostProcessor.class.getName());

    @Override
    public void postProcessEnvironment(ConfigurableEnvironment environment, SpringApplication application) {
        String enabled   = environment.getProperty("rune.enabled", "true");
        String url       = environment.getProperty("rune.url");
        String token     = environment.getProperty("rune.token");
        String namespace = environment.getProperty("rune.namespace", "default");
        String failFast  = environment.getProperty("rune.fail-fast", "true");

        if (!"true".equalsIgnoreCase(enabled)) {
            log.info("[rune] Starter disabled via rune.enabled=false");
            return;
        }

        if (url == null || url.isBlank() || token == null || token.isBlank()) {
            if ("true".equalsIgnoreCase(failFast)) {
                throw new RuneException("[rune] rune.url and rune.token must be configured");
            }
            log.warning("[rune] rune.url or rune.token not configured; skipping secret injection");
            return;
        }

        try {
            RuneClient client = new RuneClient(url, token);
            RunePropertySource source = new RunePropertySource("rune", client, namespace);
            environment.getPropertySources().addFirst(source);
        } catch (RuneException e) {
            if ("true".equalsIgnoreCase(failFast)) {
                throw new RuneException("[rune] Failed to load secrets from Rune server: " + e.getMessage(), e);
            }
            log.warning("[rune] Could not load secrets (fail-fast=false): " + e.getMessage());
        }
    }
}
