data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "./loader",
  ]
}

env "development" {
  src = data.external_schema.gorm.url
  dev = "postgresql://postgres:password@localhost:5432/vn_dev?sslmode=disable"
  url = "postgresql://postgres:password@localhost:5432/vn?sslmode=disable"
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \" \" }}"
    }
  }
  exclude = ["atlas_schema_revisions"]
}

diff {
  skip {
    drop_schema = true
    drop_table = true
  }
}