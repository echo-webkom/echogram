resource "azurerm_storage_container" "storage_container" {
  name                  = "${var.prefix}-storage-container"
  storage_account_name  = var.sa_name
  container_access_type = "private"
}
