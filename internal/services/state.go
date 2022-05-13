package services

import (
	"context"
	"github.com/NethermindEth/juno/internal/config"
	"github.com/NethermindEth/juno/internal/log"
	"github.com/NethermindEth/juno/pkg/db"
	"github.com/NethermindEth/juno/pkg/db/state"
	"go.uber.org/zap"
	"sync"
)

// StateService is the service to store and search all the information related with the state,
// that is contract code and contract storage
var StateService stateService

type stateService struct {
	running bool
	manager *state.Manager
	logger  *zap.SugaredLogger
	wg      sync.WaitGroup
}

// Run starts the service. After this the service is ready to store and search the information.
func (service *stateService) Run() error {
	// Init service logger
	service.logger = log.Default.Named(log.NameSlot("State service"))

	// Check if the service is already running
	if service.running {
		service.logger.Warn("Service is already running")
		return nil
	}
	// Init the databases
	codeDatabase := db.NewKeyValueDb(config.Runtime.DbPath+"/code", 0)
	storageDatabase := db.NewBlockSpecificDatabase(db.NewKeyValueDb(config.Runtime.DbPath+"/storage", 0))
	service.manager = state.NewStateManager(codeDatabase, *storageDatabase)

	// Set service as running
	service.running = true

	service.logger.Info("Service running")
	return nil
}

// Close closes the service.
func (service *stateService) Close(_ context.Context) {
	// Check if the service is running
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

// StoreCode stores the contract code.
func (service *stateService) StoreCode(contractAddress string, code state.ContractCode) {
	service.wg.Add(1)
	defer service.wg.Done()

	service.logger.With("Contract address", contractAddress).Debug("StoreCode")

	service.manager.PutCode(contractAddress, &code)
}

// GetCode search the contract code for the given contract address.
func (service *stateService) GetCode(contractAddress string) *state.ContractCode {
	service.wg.Add(1)
	defer service.wg.Done()

	service.logger.With("Contract address", contractAddress).Debug("GetCode")

	code := service.manager.GetCode(contractAddress)

	return code
}

// StoreStorage stores the contract storage at the given block. Notice: this does not update the storage, if you want
// to update the previous storage, then use the Update method.
func (service *stateService) StoreStorage(contractAddress string, blockNumber uint64, storage state.ContractStorage) {
	service.wg.Add(1)
	defer service.wg.Done()

	service.logger.With("Contract address", contractAddress, "Block number", blockNumber).Debug("StoreStorage")

	service.manager.PutStorage(contractAddress, blockNumber, &storage)
}

// GetStorage search the updated contract storage for the given contract address and block number.
func (service *stateService) GetStorage(contractAddress string, blockNumber uint64) *state.ContractStorage {
	service.wg.Add(1)
	defer service.wg.Done()

	service.logger.With("Contract address", contractAddress, "Block number", blockNumber).Debug("GetStorage")

	storage := service.manager.GetStorage(contractAddress, blockNumber)

	return storage
}

// UpdateStorage updates the storage for the given contract address at the given block number.
func (service *stateService) UpdateStorage(contractAddress string, blockNumber uint64, storage state.ContractStorage) {
	service.wg.Add(1)
	defer service.wg.Done()

	service.logger.With("Contract address", contractAddress, "Block number", blockNumber).Debug("UpdateStorage")

	oldStorage := service.manager.GetStorage(contractAddress, blockNumber)
	if oldStorage == nil {
		service.manager.PutStorage(contractAddress, blockNumber, &storage)
	} else {
		oldStorage.Update(storage)
		service.manager.PutStorage(contractAddress, blockNumber, oldStorage)
	}
}
