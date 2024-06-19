package main

type Configs struct {
	App      App
	Database Database
}

type App struct {
	Port  int
	Azure Azure
}

type Azure struct {
	SecretName  string
	KeyVaultURL string
	SecretVault string
}

type Database struct {
	Address string
}
