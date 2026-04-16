env "local" {
  src = "file://migrations"
  dev = getenv("ATLAS_DEV_URL")
  url = getenv("DATABASE_URL")
}