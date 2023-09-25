# TokenMigration contract build commands

1. build - solc --bin --abi contract/migration_contract.sol -o build-muticall
2. abigen - abigen --bin=build_multicall/Migrator.bin --abi=build_multicall/Migrator.abi --pkg=multicall --out=gen-multicall/migrator.go

# Run make file

-> make -f MakeFile run
