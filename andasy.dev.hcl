# andasy.dev.hcl app configuration file for nfs-api-dev (development)

app_name = "nfs-api-dev"

app {

  env = {}

  port = 5002

  compute {
    cpu      = 1
    memory   = 1024
    cpu_kind = "shared"
  }

  process {
    name = "nfs-api-dev"
  }

}
