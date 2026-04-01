package cmd

import (
	internalConfig "github.com/albert-upert/template-backend-users/config"
	"github.com/albert-upert/template-backend-users/module/redirect"
	"github.com/albert-upert/template-backend-users/module/roles"
	"github.com/albert-upert/template-backend-users/module/users"
	"github.com/albert-upert/template-backend-utils-framework/common/logger"
	"github.com/albert-upert/template-backend-utils-framework/config"
	"github.com/albert-upert/template-backend-utils-framework/fiber"
	"github.com/albert-upert/template-backend-utils-framework/otel"
	"github.com/albert-upert/template-backend-utils-framework/postgres"
	"github.com/albert-upert/template-backend-utils-framework/redis"
	"github.com/albert-upert/template-backend-utils-libraries/idp"
	"github.com/albert-upert/template-backend-utils-libraries/middleware"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var serve = &cobra.Command{
	Use:   "serve",
	Short: "Start the application",
	RunE:  serveE,
}

func serveE(cmd *cobra.Command, args []string) error {
	fx.New(
		// provide libraries
		logger.Module,
		config.Module,
		fiber.Module,
		otel.Module,
		postgres.Module,
		redis.Module,

		// supply config source & resolvers
		fx.Supply(
			fx.Annotated{
				Group: "config_options",
				Target: config.WithSources(
					// FileSource (Base Config / Local Dev)
					config.FileSource("config/config.json"),

					// Consul KV (Dynamic Config)
					// config.KVSource(
					// 	// Prefix: config/users/<env>/
					// 	fmt.Sprintf("config/users/%s/", func() string {
					// 		env := os.Getenv("APP_ENV")
					// 		if env == "" {
					// 			return "dev"
					// 		}
					// 		return env
					// 	}()),

					// 	// Client
					// 	config.ConsulClient(
					// 		os.Getenv("CONSUL_HTTP_ADDR"),
					// 		os.Getenv("CONSUL_HTTP_TOKEN"),
					// 	),

					// 	config.KVDefaultMapper(3),
					// ),

					// EnvSource (Highest Priority Overrides)
					config.EnvSource("APP", config.DefaultEnvMapper()),
				),
			},

			fx.Annotated{
				Group: "config_options",
				Target: config.WithResolvers(
					// FileResolver creates a resolver that reads values from files.
					// It resolves values with the "file://" scheme.
					//
					// Examples:
					//   - "file:///etc/secrets/api.key" -> contents of /etc/secrets/api.key
					//   - "file://./config/private.pem" -> contents of ./config/private.pem
					//
					// Returns:
					//   - Resolver: A resolver that can be used with ReadInConfig
					//
					// Security Note:
					// The path is validated to prevent directory traversal attacks.
					// Only files within local paths are generally allowed, depending on validation logic.
					config.FileResolver(),

					// EnvResolver creates a resolver that fetches values from environment variables.
					// It resolves values with the "env://" scheme.
					//
					// Examples:
					//   - "env://DATABASE_URL" -> value of DATABASE_URL environment variable
					//   - "env://API_KEY" -> value of API_KEY environment variable
					//
					// Returns:
					//   - Resolver: A resolver that can be used with ReadInConfig
					config.EnvResolver(),

					// Base64Resolver creates a resolver that decodes base64-encoded values.
					// It resolves values with the "base64://" scheme.
					//
					// Examples:
					//   - "base64://SGVsbG8gV29ybGQ=" -> "Hello World"
					//   - "base64://c2VjcmV0LXRva2Vu" -> "secret-token"
					//
					// Returns:
					//   - Resolver: A resolver that can be used with ReadInConfig
					config.Base64Resolver(),

					// VaultResolver creates a resolver that fetches secrets from HashiCorp Vault.
					// It resolves values with the "vault://" scheme and returns JSON-encoded secrets.
					//
					// Parameters:
					//   - url: Vault server URL (e.g., "https://vault:8200")
					//   - token: Vault authentication token
					//   - mountPath: Vault secrets engine mount path (e.g., "secret")
					//   - vaultPath: Base path for secrets (can be empty if secrets are at root)
					//
					// Examples:
					//   - "vault://database" -> JSON secret from secret/data/database
					//   - "vault://api/keys" -> JSON secret from secret/data/api/keys
					//
					// Returns:
					//   - Resolver: A resolver that can be used with ReadInConfig
					//
					// Note:
					// The resolver returns the JSON representation of the secret from Vault.
					// The key is extracted from the portion of the string after "vault://".
					// config.VaultResolver("", "", "", ""),
				),
			},
		),

		// provide used configurations
		fx.Provide(
			config.ProvideConfig[internalConfig.ApplicationConfig](),
			internalConfig.Postgres,
			internalConfig.Redis,
			internalConfig.Fiber,
			internalConfig.Otel,
			internalConfig.Logger,
			internalConfig.InternalApp,
			middleware.NewAuthorizationMiddleware,
		),

		middleware.HealthModule,
		roles.Module,
		users.Module,
		redirect.Module,
		fx.Provide(
			fx.Annotate(
				idp.NewIDP,
				fx.As(new(idp.IDPProvider)),
			),
		),
	).Run()

	return nil
}
