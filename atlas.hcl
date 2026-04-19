data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "ariga.io/atlas-provider-gorm/gormschema",
    "--path", "./internal/model",
    "--dialect", "postgres",
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url
  dev = "docker://postgres/15/dev"
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

env "local" {
  src = "file://migrations"
  dev = getenv("ATLAS_DEV_URL")
  url = getenv("DATABASE_URL")
  migration {
    dir              = "file://migrations"
    baseline         = "202604180001"
    revisions_schema = "public"
  }
}
