// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract ImageRegistry {
    // Mapping to record trusted images (using a 32-byte hash as key)
    mapping(bytes32 => bool) public trustedImages;

    event ImageRegistered(bytes32 indexed imageHash);

    // Register an image hash as trusted
    function registerImage(bytes32 imageHash) public {
        trustedImages[imageHash] = true;
        emit ImageRegistered(imageHash);
    }

    // Query if an image hash is trusted
    function isImageTrusted(bytes32 imageHash) public view returns (bool) {
        return trustedImages[imageHash];
    }
}
