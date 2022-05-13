package services

import (
	"context"
	"github.com/NethermindEth/juno/internal/config"
	"github.com/NethermindEth/juno/internal/log"
	"github.com/NethermindEth/juno/pkg/db"
	"github.com/NethermindEth/juno/pkg/db/abi"
	"go.uber.org/zap"
	"sync"
)

// ABIService is the ABI lookup and storage service for contracts.
// This service is unique for the entire application. To start the service, call the Run method.
var ABIService abiService

type abiService struct {
	running bool
	manager *abi.Manager
	logger  *zap.SugaredLogger
	wg      sync.WaitGroup
}

// Run starts the service. After this, the service is ready to save and search ABI.
func (service *abiService) Run() error {
	// Check if the service is already running
	if service.running {
		service.logger.Warn("Service is already running")
		return nil
	}
	// Init the database
	database := db.NewKeyValueDb(config.Runtime.DbPath+"/abi", 0)
	service.manager = abi.NewABIManager(database)

	service.logger = log.Default.Named(log.NameSlot("ABI service"))
	// Set service as running
	service.running = true
	service.logger.Info("service started")
	return nil
}

// Close closes the service.
func (service *abiService) Close(_ context.Context) {
	if !service.running {
		service.logger.Warn("Service is not running")
		return
	}

	service.logger.Info("Closing service...")
	service.running = false
	service.logger.Info("Waiting to finish current requests...")
	service.wg.Wait()
	// TODO: Close the database manager
	service.logger.Info("Closed")
}

// StoreABI stores the ABI in the database. If any error occurs then panic.
func (service *abiService) StoreABI(contractAddress string, abi abi.Abi) {
	service.wg.Add(1)
	defer service.wg.Done()

	service.logger.With("Contract address", contractAddress).Debug("StoreABI")

	service.manager.PutABI(contractAddress, &abi)
}

// GetABI searches for the ABI associated with the given contract address.
// If the ABI is not found, then it returns nil. If any error happens then panic.
func (service *abiService) GetABI(contractAddress string) *abi.Abi {
	service.wg.Add(1)
	defer service.wg.Done()

	service.logger.With("Contract address", contractAddress).Debug("GetABI")

	return service.manager.GetABI(contractAddress)
}
