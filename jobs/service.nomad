job "service" {

  datacenters = ["dc1"]
  type = "batch"

  group "simulation" {

    task "simulated-service" {
      driver = "raw_exec"
      config {
        command = "/bin/echo"
        args    = ["hello", "nomad"]
      }

      resources {
        cpu    = 500 # 500 MHz
        memory = 256 # 256MB

        network {
          mbits = 10
          port "db" {}
        }
      }
    }
  }
}
