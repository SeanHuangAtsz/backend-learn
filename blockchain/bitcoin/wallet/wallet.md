# Wallet
## P2PK(pay to public key)
![alt text](<pictures/image copy 2.png>)
### ScriptPubKey
![alt text](<pictures/image copy.png>)
### ScriptSig
![alt text](<pictures/image copy 3.png>)
### Address
- A P2PK does not have an official address.
## P2PKH(pay to public key hash)
![alt text](<pictures/image copy 4.png>)
### ScriptPubKey
![alt text](<pictures/image copy 5.png>)
### ScriptSig
![alt text](<pictures/image copy 6.png>)
### Address
![alt text](pictures/image.png)
![alt text](<pictures/image copy 7.png>)
- Different versions
    - p2pkh: 00
    - p2sh: 05
    - p2pkh(testnet): 6f
    - p2sh(testnet): c4
- Uncompressed address to compressed address
![alt text](<pictures/image copy 8.png>)
## P2WPKH(pay to witness public key hash)
### Script
![alt text](pictures/image-1.png)
- It works in the same way as a legacy P2PKH, but it gets unlocked via the Witness field instead of the ScriptSig.
- Script pubkey
![alt text](pictures/image-4.png)
- Witness field
![alt text](pictures/image-2.png)
![alt text](pictures/image-5.png)
- Execution
![!\[alt text\](image-2.png)](pictures/image-3.png)
### Address
- The address for a P2WPKH locking script is the Bech32 encoding of the ScriptPubKey. The ScriptPubKey for a P2WPKH has the following structure: 0014<20-byte hash160(public key)>
![alt text](pictures/image-6.png)