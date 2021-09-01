# Terraform provider for SQLite database engine

**!!! WARNING !!!**<br>
This is an educational project. Not intended for any production use!
<br>**!!! WARNING !!!**

Here is my first implemetation of SQLite plugin (provider) for Terraform.

There's no many docs and tests yet. The code is intended for local development and use for now.

Plugin was tested on MacOS X, but should work on any other platform as well.

Prequisites:
* Go >= 1.16
* Terraform >= 1.0

Provider supports:
* Creating database tables with basic constraints like `NOT NULL`, `PRIMARY KEY` and `DEFAULT`
* Creating database basic indexes

Provider **does not** support:
* Table/index schema upgrades (in case of schema changes resources will be recreated)
* Resources import

Known issues:
* If you delete/recreate table having any indexes on it, provider will not automatically recreate indexes; it will be recreated on the next run of `terraform apply`

## How to build
1. Clone repository and make Go download dependencies
   ```shell
   git clone https://github.com/Burmuley/terraform-provider-sqlite
   cd terraform-provider-sqlite
   go get
   ```
2. Run script `build.sh` to build and install provider to the Terraform local cache, or run these commands to achieve the same
   ```shell
   go build -o terraform-provider-sqlite
   mkdir -p ~/.terraform.d/plugins/burmuley.com/edu/sqlite/0.1/darwin_amd64
   mv terraform-provider-sqlite ~/.terraform.d/plugins/burmuley.com/edu/sqlite/0.1/darwin_amd64
   ```

## How to use
And example code is located in the `example` directory.
```shell
cd example
terraform init
```
Then Terraform should find local provider and successfully initialize the local module.
```shell
Initializing the backend...

Initializing provider plugins...
- Reusing previous version of burmuley.com/edu/sqlite from the dependency lock file
- Installing burmuley.com/edu/sqlite v0.1.0...
- Installed burmuley.com/edu/sqlite v0.1.0 (unauthenticated)

Terraform has been successfully initialized!

```

Simply run the `apply` command and confirm you want to provision "the infrastructure". :)
```shell
terraform apply
```
```shell
Terraform will perform the following actions:

  # sqlite_index.test_index1 will be created
  + resource "sqlite_index" "test_index1" {
      + columns = [
          + "id",
          + "name",
        ]
      + created = (known after apply)
      + id      = (known after apply)
      + name    = "users_index"
      + table   = "users"
    }

  # sqlite_table.test_table will be created
  + resource "sqlite_table" "test_table" {
      + created = (known after apply)
      + id      = (known after apply)
      + name    = "users"

      + column {
          + name = "id"
          + type = "INTEGER"

          + constraints {
              + not_null    = true
              + primary_key = true
            }
        }
      + column {
          + name = "name"
          + type = "TEXT"

          + constraints {
              + not_null    = true
              + primary_key = false
            }
        }
      + column {
          + name = "last_name"
          + type = "TEXT"

          + constraints {
              + not_null    = true
              + primary_key = false
            }
        }
      + column {
          + name = "password"
          + type = "TEXT"

          + constraints {
              + default     = "123"
              + not_null    = true
              + primary_key = false
            }
        }
    }

  # sqlite_table.test_table2 will be created
  + resource "sqlite_table" "test_table2" {
      + created = (known after apply)
      + id      = (known after apply)
      + name    = "projects"

      + column {
          + name = "id"
          + type = "INTEGER"

          + constraints {
              + not_null    = true
              + primary_key = true
            }
        }
      + column {
          + name = "user"
          + type = "INTEGER"

          + constraints {
              + not_null    = true
              + primary_key = false
            }
        }
      + column {
          + name = "name"
          + type = "TEXT"

          + constraints {
              + not_null    = true
              + primary_key = false
            }
        }
    }

Plan: 3 to add, 0 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

sqlite_table.test_table2: Creating...
sqlite_table.test_table: Creating...
sqlite_table.test_table2: Creation complete after 0s [id=projects]
sqlite_table.test_table: Creation complete after 0s [id=users]
sqlite_index.test_index1: Creating...
sqlite_index.test_index1: Creation complete after 0s [id=users_index]

Apply complete! Resources: 3 added, 0 changed, 0 destroyed.
```
