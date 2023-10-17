provider "feldera" {
  endpoint = "http://localhost:8080"
}

/** For local testing
You can use the `make build-local` for building the plugin
 and then use the following code for testing the plugin locally.
terraform {
  required_providers {
    feldera = {
      source  = "terraform.local/local/feldera"
      version = "1.0.0"
    }
  }
}

provider "feldera" {
  endpoint = "http://localhost:8080"
}
*/
