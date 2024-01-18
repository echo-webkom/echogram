resource "azurerm_resource_group" "rg" {
  name     = "${var.prefix}-resources"
  location = local.location

  tags = {
    environment = local.environment
  }
}

resource "azurerm_storage_account" "sa" {
  name                     = "${var.prefix}storage"
  resource_group_name      = azurerm_resource_group.rg.name
  location                 = local.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

module "function" {
  source = "./modules/function"

  rg_name                      = azurerm_resource_group.rg.name
  sa_name                      = azurerm_storage_account.sa.name
  sa_primary_connection_string = azurerm_storage_account.sa.primary_connection_string
  sa_primary_access_key        = azurerm_storage_account.sa.primary_access_key
  location                     = local.location
  prefix                       = var.prefix
}

module "blob_storage" {
  source = "./modules/blob_storage"

  sa_name = azurerm_storage_account.sa.name
  prefix  = var.prefix
}