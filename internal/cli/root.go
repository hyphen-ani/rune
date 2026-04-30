package cli

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{

	Use:   "rune",
	Short: "A lightweight, secure secrets manager",
	Long: `
      _ __ _   _ _ __   ___
     | '__| | | | '_ \ / _ \
     | |  | |_| | | | |  __/
     |_|   \__,_|_| |_|\___|

rune is a minimal, local-first secrets manager designed for simplicity and security.
It provides encrypted storage for secrets using AES-GCM, with keys derived from a user passphrase via Argon2.

The vault starts in a sealed state and must be explicitly unsealed before accessing any data.

Key features:

- Encrypted storage at rest
- Passphrase-based key derivation
- Explicit seal/unseal workflow
- Lightweight and zero external dependencies

Use rune to securely manage secrets in local or self-hosted environments without heavy infrastructure.
`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
