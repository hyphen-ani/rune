package io.rune.sdk;

import org.springframework.core.env.MapPropertySource;

import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.logging.Logger;

/**
 * Spring PropertySource that eagerly fetches all secrets for a given namespace
 * from the Rune server and exposes them as Spring properties with the prefix "rune.".
 *
 * A Rune secret with key "db.password" in namespace "myapp" is exposed as
 * the Spring property "rune.db.password", allowing:
 *
 *   spring.datasource.password=${rune.db.password}
 */
public class RunePropertySource extends MapPropertySource {

    private static final Logger log = Logger.getLogger(RunePropertySource.class.getName());

    public RunePropertySource(String name, RuneClient client, String namespace) {
        super(name, loadSecrets(client, namespace));
    }

    private static Map<String, Object> loadSecrets(RuneClient client, String namespace) {
        Map<String, Object> secrets = new HashMap<>();
        List<String> keys = client.listKeys(namespace);
        for (String key : keys) {
            String value = client.getSecret(namespace, key);
            secrets.put("rune." + key, value);
        }
        log.info("[rune] Loaded " + secrets.size() + " secret(s) from namespace '" + namespace + "'");
        return secrets;
    }
}
