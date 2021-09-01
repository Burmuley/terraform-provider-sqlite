terraform {
  required_providers {
    sqlite = {
      version = "0.1"
      source = "burmuley.com/edu/sqlite"
    }
  }
}

provider "sqlite" {
  path = "local.db"
}
