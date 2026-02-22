# andasy.hcl app configuration file generated for nfs-api on Saturday, 21-Feb-26 11:56:35 CAT
#
# See https://github.com/quarksgroup/andasy-cli for information about how to use this file.

app_name = "nfs-api"

app {

  env = {
    APP_ENV = "production"
  }

  port = 5002

  compute {
    cpu      = 1
    memory   = 1024
    cpu_kind = "shared"
  }

  process {
    name = "nfs-api"
  }

}
