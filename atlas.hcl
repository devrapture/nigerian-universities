env "local" {
  src = attr.migration.dir
  dev = getenv("ATLAS_DEV_URL")
  url = getenv("DATABASE_URL")
}