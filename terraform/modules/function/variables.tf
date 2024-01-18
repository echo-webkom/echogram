variable "rg_name" {
  type        = string
  description = "The name of the resource group."
}

variable "sa_name" {
  type        = string
  description = "The name of the storage account."
}

variable "sa_primary_connection_string" {
  type        = string
  description = "The primary access key for the storage account."
}

variable "sa_primary_access_key" {
  type        = string
  description = "The primary access key for the storage account."
}

variable "location" {
  type        = string
  description = "The location/region where the resource group will be created."
}

variable "prefix" {
  type        = string
  description = "The prefix which should be used for all resources in this example"
}
