const express = require('express');
const app = express();
const { ethers } = require('ethers');
const crypto = require('crypto');

app.use(express.json());

// --- Configuration ---

const imageRegistryABI = [
  {
    anonymous: false,
    inputs: [
      {
        indexed: true,
        internalType: 'bytes32',
        name: 'imageHash',
        type: 'bytes32',
      },
    ],
    name: 'ImageRegistered',
    type: 'event',
  },
  {
    inputs: [
      {
        internalType: 'bytes32',
        name: 'imageHash',
        type: 'bytes32',
      },
    ],
    name: 'registerImage',
    outputs: [],
    stateMutability: 'nonpayable',
    type: 'function',
  },
  {
    inputs: [
      {
        internalType: 'bytes32',
        name: 'imageHash',
        type: 'bytes32',
      },
    ],
    name: 'isImageTrusted',
    outputs: [
      {
        internalType: 'bool',
        name: '',
        type: 'bool',
      },
    ],
    stateMutability: 'view',
    type: 'function',
  },
  {
    inputs: [
      {
        internalType: 'bytes32',
        name: '',
        type: 'bytes32',
      },
    ],
    name: 'trustedImages',
    outputs: [
      {
        internalType: 'bool',
        name: '',
        type: 'bool',
      },
    ],
    stateMutability: 'view',
    type: 'function',
  },
];

// Replace with your deployed contract address (from Ganache).
const imageRegistryAddress = '0x5CeC6C7F106d58d29Dd89B9cd0302f39f3DbcDD0';

// Connect to Ganache (assumed to be running on http://localhost:8545).
const provider = new ethers.providers.JsonRpcProvider('http://localhost:8545');

// Replace <PRIVATE_KEY> with one of the private keys provided by Ganache.
const privateKey =
  '0x79e65362c62ceba88406297944a4063fd821621509bef8e3c8edba558d546ec3';
const signer = new ethers.Wallet(privateKey, provider);

// Create a contract instance with a signer (to send transactions).
const imageRegistryContract = new ethers.Contract(
  imageRegistryAddress,
  imageRegistryABI,
  signer
);

// --- Endpoint to Register Image on Blockchain ---

app.post('/register', async (req, res) => {
  const { repository, digest } = req.body;
  console.log(
    `Received image registration for repository: ${repository}, digest: ${digest}`
  );

  // Assume digest is provided in a format like "sha256:<hash>".
  // Remove the "sha256:" prefix if present.
  let imageHash = digest.startsWith('sha256:') ? digest.slice(7) : digest;

  // If the hash is not 64 hex characters (32 bytes), compute a SHA256 hash of repository+digest.
  if (imageHash.length !== 64) {
    imageHash = crypto
      .createHash('sha256')
      .update(repository + digest)
      .digest('hex');
  }

  // Ensure the hash has a 0x prefix.
  imageHash = '0x' + imageHash;
  console.log('Computed image hash:', imageHash);

  try {
    // Call the smart contract's registerImage function.
    let tx = await imageRegistryContract.registerImage(imageHash);
    console.log('Transaction sent. Waiting for confirmation...');
    let receipt = await tx.wait();
    console.log('Transaction confirmed:', receipt.transactionHash);
    res.status(200).json({
      status: 'Image registered successfully',
      txHash: receipt.transactionHash,
    });
  } catch (error) {
    console.error('Error registering image on blockchain:', error);
    res
      .status(500)
      .json({ status: 'Error registering image', error: error.toString() });
  }
});

const PORT = 5001;
app.listen(PORT, () => {
  console.log(`Blockchain registration service listening on port ${PORT}`);
});
