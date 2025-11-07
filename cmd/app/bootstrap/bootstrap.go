package bootstrap

import (
	"context"
	"fmt"
	"log"
	"os"

	"akeneo-migrator/internal/config"
	"akeneo-migrator/internal/platform/client/akeneo"
	akeneo_storage "akeneo-migrator/internal/platform/storage/akeneo"
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

	// 6. Create services
	referenceEntitySyncer := syncing.NewService(sourceRepository, destRepository)

	// 7. Create application with dependencies
	app := &Application{
		Config:                cfg,
		SourceClient:          sourceClient,
		DestClient:            destClient,
		SourceRepository:      sourceRepository,
		DestRepository:        destRepository,
		ReferenceEntitySyncer: referenceEntitySyncer,
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
		debug, _ := cmd.Flags().GetBool("debug")

		fmt.Printf("🚀 Starting synchronization for entity: %s\n", entityName)
		if debug {
			fmt.Println("🔍 Debug mode enabled")
		}

		// Execute synchronization using the service
		fmt.Printf("📥 Fetching records from entity '%s'...\n", entityName)
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
		os.Setenv("ENVIRONMENT", "local")
	}
	if os.Getenv("CONFIG_PATH") == "" {
		os.Setenv("CONFIG_PATH", "akeneo-migrator")
	}
}
