resource "feldera_program" "example" {
  name = "my-example-program"
  description = "This is an example program"
  code = <<EOT
      create table VENDOR (
        id bigint not null primary key,
        name varchar,
        address varchar
      );
  EOT
}