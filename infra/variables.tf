variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "db_password" {
  description = "RDS PostgreSQL password"
  type        = string
  sensitive   = true
}

variable "polygon_api_key" {
  description = "Polygon.io API key for real market data"
  type        = string
  sensitive   = true
}
