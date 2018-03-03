job "evelina" {
  type        = "batch"
  datacenters = ["dc1"]

  meta {
    pr = ""
  }

  parameterized {
    meta_required = ["pr"]
  }

  task "tc" {
    driver = "docker"

    config {
      image = "remmelt/go-builder"
      command = "echo"
      args    = ["${NOMAD_META_PR}"]
//      args    = ["version"]
    }

    resources {
      cpu    = 1000
      memory = 256
    }
  }
}
