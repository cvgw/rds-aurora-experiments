variable "engine" {
  default = "aurora-mysql"
}
variable "engine_version" {
  default = "5.7.12"
}
variable "instance_class" {
  default = "db.t2.small"
}
variable "description" {
  default = "test subnet group created by tf"
}
variable "instance_id" {}
variable "cluster_id" {}
variable "subnet_group_name" {}
variable "master_username" {}
variable "master_user_password" {}
variable "sg_ids" {
  type = "list"
}
variable "subnet_ids" {
  type = "list"
}


resource "aws_db_subnet_group" "default" {
  name = "${var.subnet_group_name}"
  description = "${var.description}"
  subnet_ids = "${var.subnet_ids}"
}

resource "aws_rds_cluster" "default" {
  cluster_identifier      = "${var.cluster_id}"
  engine                  = "${var.engine}"
  db_subnet_group_name    = "${aws_db_subnet_group.default.name}"
  master_username         = "${var.master_username}"
  master_password         = "${var.master_user_password}"
  vpc_security_group_ids  = "${var.sg_ids}"
}

resource "aws_rds_cluster_instance" "default" {
  count              = 1
  identifier         = "${var.instance_id}"
  cluster_identifier = "${aws_rds_cluster.default.id}"
  engine             = "${var.engine}"
  instance_class     = "${var.instance_class}"
}
