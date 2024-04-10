// SPDX-License-Identifier: BUSL-1.1
pragma solidity =0.8.12;

import "src/core/EigenDAServiceManager.sol";
import "eigenlayer-scripts/utils/ExistingDeploymentParser.sol";
import "forge-std/Test.sol";
import "forge-std/Script.sol";
import "forge-std/StdJson.sol";

contract Deployer_Mainnet is ExistingDeploymentParser {

    string public existingCoreDeploymentPath  = string(bytes("./script/deploy/mainnet/mainnet_addresses.json"));
    string public existingDADeploymentPath = string(bytes("script/deploy/mainnet/mainnet_deployment_data.json"));

    address registryCoordinator;
    address stakeRegistry;

    EigenDAServiceManager eigenDAServiceManagerImplementation;

    function run() external {
        _parseDeployedContracts(existingCoreDeploymentPath);
        registryCoordinator = stdJson.readAddress(existingDADeploymentPath, ".addresses.registryCoordinator");
        stakeRegistry = stdJson.readAddress(existingDADeploymentPath, ".addresses.stakeRegistry");

        vm.startBroadcast();

        eigenDAServiceManagerImplementation = new EigenDAServiceManager(
            avsDirectory,
            IRegistryCoordinator(registryCoordinator),
            IStakeRegistry(stakeRegistry)
        );

        vm.stopBroadcast();
    }
}