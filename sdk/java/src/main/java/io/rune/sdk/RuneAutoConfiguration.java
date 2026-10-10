package io.rune.sdk;

import org.springframework.boot.autoconfigure.AutoConfiguration;
import org.springframework.boot.autoconfigure.condition.ConditionalOnProperty;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.Bean;

/**
 * Exposes RuneClient and RuneProperties as Spring beans for application code
 * that wants to do programmatic secret lookups beyond property injection.
 *
 * The actual property injection happens earlier, in RuneEnvironmentPostProcessor.
 */
@AutoConfiguration
@ConditionalOnProperty(prefix = "rune", name = "enabled", havingValue = "true", matchIfMissing = true)
@EnableConfigurationProperties(RuneProperties.class)
public class RuneAutoConfiguration {

    @Bean
    @ConditionalOnProperty(prefix = "rune", name = {"url", "token"})
    public RuneClient runeClient(RuneProperties props) {
        return new RuneClient(props.getUrl(), props.getToken());
    }
}
