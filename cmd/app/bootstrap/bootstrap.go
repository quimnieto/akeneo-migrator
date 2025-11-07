package bootstrap

import (
	"context"
	"fmt"
	"log"
	"os"

	"akeneo-migrator/internal/attribute"
	attribute_syncing "akeneo-migrator/internal/attribute/syncing"
	"akeneo-migrator/internal/config"
	"akeneo-migrator/internal/platform/client/akeneo"
	akeneo_storage "akeneo-migrator/internal/platform/storage/akeneo"
	"akeneo-migrator/internal/product"
	product_syncing "akeneo-migrator/internal/product/syncing"
	"akeneo-migrator/internal/reference_entity"
	"akeneo-migrator/internal/reference_entity/syncing"
	"akeneo-migrator/kit/config/static/viper"

	"github.com/spf13/cobra"
)

const CONTEXT = "akeneo-migrator"

// Application contains all application dependencies
type Application struct {
	Config                *config.Config
	SourceClient          *akeneo.Client
	DestClient            *akeneo.Client
	SourceRepository      reference_entity.SourceRepository
	DestRepository        reference_entity.DestRepository
	ReferenceEntitySyncer *syncing.Service
	SourceProductRepo     product.SourceRepository
	DestProductRepo       product.DestRepository
	ProductSyncer         *product_syncing.Service
	SourceAttributeRepo   attribute.SourceRepository
	DestAttributeRepo     attribute.DestRepository
	AttributeSyncer       *attribute_syncing.Service
}

// Run initializes the application and executes CLI commands
func Run() error {
	// 0. Setup default environment variables if not defined
	setupDefaultEnvironmentVariables()

	// 1. Load configuration with Viper
	viperConfig := viper.NewViperConfig()
	err := viperConfig.LoadConfiguration(CONTEXT)
	if err != nil {
		return err
	}

	// 2. Create configuration
	cfg, err := config.LoadConfig(viperConfig)
	if err != nil {
		return fmt.Errorf("error creating configuration: %w", err)
	}

	// 3. Create source client
	sourceClient, err := akeneo.NewClient(akeneo.ClientConfig{
		Host:     cfg.Source.Host,
		ClientID: cfg.Source.ClientID,
		Secret:   cfg.Source.Secret,
		Username: cfg.Source.Username,
		Password: cfg.Source.Password,
	})
	if err != nil {
		return fmt.Errorf("error creating source client: %w", err)
	}

	// 4. Create destination client
	destClient, err := akeneo.NewClient(akeneo.ClientConfig{
		Host:     cfg.Dest.Host,
		ClientID: cfg.Dest.ClientID,
		Secret:   cfg.Dest.Secret,
		Username: cfg.Dest.Username,
		Password: cfg.Dest.Password,
	})
	if err != nil {
		return fmt.Errorf("error creating destination client: %w", err)
	}

	// 5. Create repositories
	sourceRepository := akeneo_storage.NewSourceReferenceEntityRepository(sourceClient)
	destRepository := akeneo_storage.NewDestReferenceEntityRepository(destClient)
	sourceProductRepo := akeneo_storage.NewSourceProductRepository(sourceClient)
	destProductRepo := akeneo_storage.NewDestProductRepository(destClient)
	sourceAttributeRepo := akeneo_storage.NewSourceAttributeRepository(sourceClient)
	destAttributeRepo := akeneo_storage.NewDestAttributeRepository(destClient)

	// 6. Create services
	referenceEntitySyncer := syncing.NewService(sourceRepository, destRepository)
	productSyncer := product_syncing.NewService(sourceProductRepo, destProductRepo)
	attributeSyncer := attribute_syncing.NewService(sourceAttributeRepo, destAttributeRepo)

	// 7. Create application with dependencies
	app := &Application{
		Config:                cfg,
		SourceClient:          sourceClient,
		DestClient:            destClient,
		SourceRepository:      sourceRepository,
		DestRepository:        destRepository,
		ReferenceEntitySyncer: referenceEntitySyncer,
		SourceProductRepo:     sourceProductRepo,
		DestProductRepo:       destProductRepo,
		ProductSyncer:         productSyncer,
		SourceAttributeRepo:   sourceAttributeRepo,
		DestAttributeRepo:     destAttributeRepo,
		AttributeSyncer:       attributeSyncer,
	}

	// 8. Create root command
	rootCmd := &cobra.Command{
		Use:   "akeneo-migrator",
		Short: "CLI tool to migrate data between Akeneo instances",
		Long: `akeneo-migrator is a CLI tool that allows you to synchronize data
between different Akeneo PIM instances, including Reference Entities,
products, categories and other elements.`,
	}

	// 9. Add commands
	syncCmd := createSyncCommand(app)
	rootCmd.AddCommand(syncCmd)

	syncProductCmd := createSyncProductCommand(app)
	rootCmd.AddCommand(syncProductCmd)

	syncAttributeCmd := createSyncAttributeCommand(app)
	rootCmd.AddCommand(syncAttributeCmd)

	// 10. Execute root command
	return rootCmd.Execute()
}

// createSyncCommand creates the sync command
func createSyncCommand(app *Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync [entity-name]",
		Short: "Synchronizes all records from a specific Reference Entity",
		Long: `Synchronizes all records from a Reference Entity from the source Akeneo 
to the destination Akeneo. Requires the entity name as an argument.

Example:
  akeneo-migrator sync brands
  akeneo-migrator sync brands --debug`,
		Args: cobra.ExactArgs(1),
		Run:  runSyncCommand(app),
	}

	// Add debug mode flag
	cmd.Flags().Bool("debug", false, "Enable debug mode to see record contents")

	return cmd
}

// runSyncCommand executes the synchronization logic
func runSyncCommand(app *Application) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		entityName := args[0]
		ctx := context.Background()

		// Get debug flag
		debug, _ := cmd.Flags().GetBool("debug") //nolint:errcheck // flag is optional

		fmt.Printf("🚀 Starting synchronization for entity: %s\n", entityName)
		if debug {
			fmt.Println("🔍 Debug mode enabled")
		}

		// Execute synchronization using the service
		fmt.Printf("📋 Synchronizing Reference Entity '%s'...\n", entityName)
		fmt.Println("   1️⃣  Syncing entity definition...")
		fmt.Println("   2️⃣  Syncing attributes...")
		fmt.Println("   3️⃣  Syncing records...")

		result, err := app.ReferenceEntitySyncer.Sync(ctx, entityName)
		if err != nil {
			log.Printf("❌ Synchronization error: %v\n", err)
			return
		}

		fmt.Printf("📊 Found %d records to synchronize\n", result.TotalRecords)

		// Show progress for each record
		if debug {
			for _, syncErr := range result.Errors {
				fmt.Printf("❌ Error in record '%s': %s\n", syncErr.Code, syncErr.Message)
			}
		}

		// Final summary
		fmt.Println("\n📋 Synchronization summary:")
		fmt.Printf("   ✅ Successfully synchronized records: %d\n", result.SuccessCount)
		fmt.Printf("   ❌ Records with errors: %d\n", result.ErrorCount)
		fmt.Printf("   📊 Total processed: %d\n", result.TotalRecords)

		if result.ErrorCount > 0 {
			fmt.Println("\n⚠️  Synchronization completed with some errors.")
			if !debug {
				fmt.Println("💡 Run with --debug to see error details")
			}
		} else {
			fmt.Println("\n🎉 Synchronization completed successfully!")
		}
	}
}

// setupDefaultEnvironmentVariables sets up default environment variables
func setupDefaultEnvironmentVariables() {
	if os.Getenv("ENVIRONMENT") == "" {
		_ = os.Setenv("ENVIRONMENT", "local")
	}
	if os.Getenv("CONFIG_PATH") == "" {
		_ = os.Setenv("CONFIG_PATH", "akeneo-migrator")
	}
}

// createSyncProductCommand creates the sync-product command
func createSyncProductCommand(app *Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync-product [identifier]",
		Short: "Synchronizes a product hierarchy by its common identifier",
		Long: `Synchronizes a complete product hierarchy from the source Akeneo to the destination Akeneo.

For SIMPLE products: Common → Child Products
For CONFIGURABLE products: Common → Models → Variant Products

Requires the common product/model identifier as an argument.

Example:
  akeneo-migrator sync-product COMMON-001
  akeneo-migrator sync-product COMMON-001 --debug
  akeneo-migrator sync-product COMMON-001 --single  # Sync only the product, not hierarchy`,
		Args: cobra.ExactArgs(1),
		Run:  runSyncProductCommand(app),
	}

	// Add flags
	cmd.Flags().Bool("debug", false, "Enable debug mode to see product contents")
	cmd.Flags().Bool("single", false, "Sync only the single product, not the entire hierarchy")

	return cmd
}

// runSyncProductCommand executes the product synchronization logic
func runSyncProductCommand(app *Application) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		identifier := args[0]
		ctx := context.Background()

		// Get flags
		debug, _ := cmd.Flags().GetBool("debug")   //nolint:errcheck // flag is optional
		single, _ := cmd.Flags().GetBool("single") //nolint:errcheck // flag is optional

		fmt.Printf("🚀 Starting synchronization for product: %s\n", identifier)
		if debug {
			fmt.Println("🔍 Debug mode enabled")
		}

		var result *product_syncing.SyncResult
		var err error

		if single {
			// Sync only the single product
			fmt.Printf("📥 Fetching product '%s' from source...\n", identifier)
			result, err = app.ProductSyncer.Sync(ctx, identifier)
		} else {
			// Sync entire hierarchy
			fmt.Printf("📥 Fetching product hierarchy for '%s' from source...\n", identifier)
			result, err = app.ProductSyncer.SyncHierarchy(ctx, identifier)
		}

		if err != nil {
			log.Printf("❌ Synchronization error: %v\n", err)
			return
		}

		// Show result
		if result.Success {
			fmt.Println("\n📋 Synchronization Summary:")
			fmt.Printf("   📦 Models synced: %d\n", result.ModelsSynced)
			fmt.Printf("   📦 Products synced: %d\n", result.ProductsSynced)
			fmt.Printf("   📊 Total synced: %d\n", result.TotalSynced)
			fmt.Printf("\n✅ Hierarchy '%s' synchronized successfully!\n", result.Identifier)
		} else {
			fmt.Printf("❌ Failed to synchronize '%s': %s\n", result.Identifier, result.Error)
		}
	}
}

// createSyncAttributeCommand creates the sync-attribute command
func createSyncAttributeCommand(app *Application) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync-attribute [code]",
		Short: "Synchronizes an attribute by its code",
		Long: `Synchronizes a single attribute from the source Akeneo to the destination Akeneo.

Requires the attribute code as an argument.

Example:
  akeneo-migrator sync-attribute sku
  akeneo-migrator sync-attribute description --debug`,
		Args: cobra.ExactArgs(1),
		Run:  runSyncAttributeCommand(app),
	}

	// Add debug flag
	cmd.Flags().Bool("debug", false, "Enable debug mode to see attribute contents")

	return cmd
}

// runSyncAttributeCommand executes the attribute synchronization logic
func runSyncAttributeCommand(app *Application) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		code := args[0]
		ctx := context.Background()

		// Get debug flag
		debug, _ := cmd.Flags().GetBool("debug") //nolint:errcheck // flag is optional

		fmt.Printf("🚀 Starting synchronization for attribute: %s\n", code)
		if debug {
			fmt.Println("🔍 Debug mode enabled")
		}

		// Execute synchronization using the service
		result, err := app.AttributeSyncer.Sync(ctx, code)
		if err != nil {
			log.Printf("❌ Synchronization error: %v\n", err)
			return
		}

		// Show result
		if result.Success {
			fmt.Printf("\n✅ Attribute '%s' synchronized successfully!\n", result.Code)
		} else {
			fmt.Printf("❌ Failed to synchronize '%s': %s\n", result.Code, result.Error)
		}
	}
}
