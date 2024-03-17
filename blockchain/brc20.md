# BRC20
## Technologies that make inscriptions possible
- Data field in transaction fields
- Taproot bitcoin update that made it possible to profitably add much more data than was possible before: photos, gifs, audio, video, etc.
- The Ordinals standard, which inherently provides the foundation for the bitcoin network to be able to assign ownership of some data(so that the code guarantees your ownership of that data) .
## Ordinals
- A special standard was invented — Ordinals Protocol. Its purpose is to assign each satoshi an ordinal number corresponding to the order in which it was created by miners (1, 2, 3, 4, …). Initially, satoshis have no order in the bitcoin network, so you have to invent one.
- Ordinals protocol sets a certain rule that you always send the lowest numbered sats when transferring.For example, if I have sats #11–20 and send you 1 satoshi, then according to the Ordinals protocol, I send you #11. This way it becomes possible to number every satoshi in the bitcoin network. In addition, we can choose which satoshis we want to transfer. If I want to send you #12, I can transfer 1 satoshi to another wallet and then 1 satoshi to you. The first satoshi will be #11 and the second will be #12.
## How do Inscriptions work?
- Inscriptions solved the owner of data problem by embedding data into satoshis. Once this is done, it can be said that the owner of that satoshi (ownership of satoshi was dealt with above in the first part about Ordinals) is the owner of the data.
- Since this satoshi is identifiable (it is numbered), we can trace its history back to the transaction in which the data was embedded in it, thus allowing us to determine that this satoshi contains such and such data.
## What is brc-20
- What BRC-20 inherently does is it allows users to simulate a balance update using multiple inscriptions.
- It does this by recording each transaction affecting the balance and calculating the final result offchain(not on the bitcoin network itself) . Since the history of changes affecting the balance will never change(they are already in the blockchain) , they can be represented using inscriptions.
### Token creation and transfer
- To create a token in a simplified version of BRC-20, I would make an inscribe with roughly the following content:
```json
{ 
 "protocol": "brc-20",
 "operation": "deploy",
 "symbol": "ABM_Token",
 "maxSupply": "21000000",
 "perOrdinalMintLimit": "1000"
}
```
- A mint can look like this:
```json
{
"protocol": "brc-20",
"operation": "mint",
"symbol": "ABM_Token",
"amount": "1000"
}
```
- To send these tokens to someone else, including selling them on some marketplace, we inscribe the transfer and the desired amount in this way:
```json
{ 
 "protocol": "brc-20",
 "operation": "transfer",
 "symbol": "ABM_Token",
 "amount": "500"
}
```
- Just like with Inscriptions, we need a special indexer to find all these BRC-20 inscribed inscriptions and calculate the balances itself. Balances are stored in bitcoin in the sense that they can be calculated from transactions, directly using the bitcoin network. But, as mentioned above, to do this, everyone must interpret these transactions in the same way (order of inscriptions, value of the tokens themselves, etc.)
