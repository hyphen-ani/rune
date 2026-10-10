package io.rune.sdk;

public class RuneException extends RuntimeException {
    public RuneException(String message) {
        super(message);
    }

    public RuneException(String message, Throwable cause) {
        super(message, cause);
    }
}
