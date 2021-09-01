resource "sqlite_table" "test_table" {
  name = "users"

  column {
    name = "id"
    type = "INTEGER"
    constraints {
      not_null = true
      primary_key = true
    }
  }

  column {
    name = "name"
    type = "TEXT"
    constraints {
      not_null = true
    }
  }

  column {
    name = "last_name"
    type = "TEXT"
    constraints {
      not_null = true
    }
  }

  column {
    name = "password"
    type = "TEXT"
    constraints {
      not_null = true
      default = "123"
    }
  }
}

resource "sqlite_table" "test_table2" {
  name = "projects"
  column {
    name = "id"
    type = "INTEGER"
    constraints {
      not_null = true
      primary_key = true
    }
  }

  column {
    name = "user"
    type = "INTEGER"
    constraints {
      not_null = true
    }
  }

  column {
    name = "name"
    type = "TEXT"
    constraints {
      not_null = true
    }
  }
}

resource "sqlite_index" "test_index1" {
  depends_on = [sqlite_table.test_table]
  name = "users_index"
  table = sqlite_table.test_table.name
  columns = ["id", "name"]
}

//resource "sqlite_index" "test_index2" {
//  name = "test_index2"
//  table = sqlite_table.test_table2.name
//  unique = true
//}
