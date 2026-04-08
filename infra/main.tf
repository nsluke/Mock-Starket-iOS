terraform {
  required_version = ">= 1.5"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

# ---------- Networking ----------

data "aws_vpc" "default" {
  default = true
}

data "aws_subnets" "default" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default.id]
  }
}

# ---------- Security Groups ----------

resource "aws_security_group" "db" {
  name        = "mockstarket-db"
  description = "RDS PostgreSQL access from App Runner"
  vpc_id      = data.aws_vpc.default.id

  ingress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"] # App Runner VPC connector will restrict this
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = { Name = "mockstarket-db" }
}

# ---------- RDS PostgreSQL ----------

resource "aws_db_instance" "postgres" {
  identifier     = "mockstarket"
  engine         = "postgres"
  engine_version = "16.4"
  instance_class = "db.t4g.micro"

  allocated_storage     = 20
  max_allocated_storage = 50
  storage_type          = "gp3"

  db_name  = "mockstarket"
  username = "mockstarket"
  password = var.db_password

  vpc_security_group_ids = [aws_security_group.db.id]
  publicly_accessible    = true # Needed for App Runner without VPC connector
  skip_final_snapshot    = true

  backup_retention_period = 7
  deletion_protection     = false

  tags = { Name = "mockstarket" }
}

# ---------- ECR Repository ----------

resource "aws_ecr_repository" "backend" {
  name                 = "mockstarket-backend"
  image_tag_mutability = "MUTABLE"
  force_delete         = true

  image_scanning_configuration {
    scan_on_push = true
  }
}

# ---------- Secrets Manager ----------

resource "aws_secretsmanager_secret" "polygon_api_key" {
  name                    = "mockstarket/polygon-api-key"
  recovery_window_in_days = 0
}

resource "aws_secretsmanager_secret_version" "polygon_api_key" {
  secret_id     = aws_secretsmanager_secret.polygon_api_key.id
  secret_string = var.polygon_api_key
}

# ---------- App Runner ----------

resource "aws_iam_role" "apprunner_ecr" {
  name = "mockstarket-apprunner-ecr"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = { Service = "build.apprunner.amazonaws.com" }
    }]
  })
}

resource "aws_iam_role_policy_attachment" "apprunner_ecr" {
  role       = aws_iam_role.apprunner_ecr.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSAppRunnerServicePolicyForECRAccess"
}

resource "aws_iam_role" "apprunner_instance" {
  name = "mockstarket-apprunner-instance"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = { Service = "tasks.apprunner.amazonaws.com" }
    }]
  })
}

resource "aws_iam_role_policy" "apprunner_secrets" {
  name = "secrets-access"
  role = aws_iam_role.apprunner_instance.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect   = "Allow"
      Action   = ["secretsmanager:GetSecretValue"]
      Resource = [aws_secretsmanager_secret.polygon_api_key.arn]
    }]
  })
}

resource "aws_apprunner_service" "backend" {
  service_name = "mockstarket-api"

  source_configuration {
    authentication_configuration {
      access_role_arn = aws_iam_role.apprunner_ecr.arn
    }

    image_repository {
      image_identifier      = "${aws_ecr_repository.backend.repository_url}:latest"
      image_repository_type = "ECR"

      image_configuration {
        port = "8080"

        runtime_environment_variables = {
          PORT                    = "8080"
          DATABASE_URL            = "postgres://${aws_db_instance.postgres.username}:${var.db_password}@${aws_db_instance.postgres.endpoint}/${aws_db_instance.postgres.db_name}?sslmode=require"
          MARKET_DATA_SOURCE      = "polygon"
          POLYGON_API_KEY         = var.polygon_api_key
          POLYGON_BASE_URL        = "https://api.polygon.io"
          POLYGON_WS_ENABLED      = "false"
          POLYGON_POLL_INTERVAL_MS = "30000"
          DEV_MODE                = "true"
          CORS_ORIGINS            = "*"
          LOG_LEVEL               = "info"
          STARTING_CASH           = "100000"
          MAX_WS_CLIENTS          = "1000"
        }
      }
    }

    auto_deployments_enabled = false
  }

  instance_configuration {
    cpu               = "0.25 vCPU"
    memory            = "0.5 GB"
    instance_role_arn = aws_iam_role.apprunner_instance.arn
  }

  health_check_configuration {
    protocol            = "TCP"
    interval            = 10
    timeout             = 5
    healthy_threshold   = 1
    unhealthy_threshold = 10
  }

  tags = { Name = "mockstarket-api" }

  depends_on = [aws_iam_role_policy_attachment.apprunner_ecr]
}

# ---------- ECR + App Runner for Web Frontend ----------

resource "aws_ecr_repository" "web" {
  name                 = "mockstarket-web"
  image_tag_mutability = "MUTABLE"
  force_delete         = true

  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "aws_apprunner_service" "web" {
  service_name = "mockstarket-web"

  source_configuration {
    authentication_configuration {
      access_role_arn = aws_iam_role.apprunner_ecr.arn
    }

    image_repository {
      image_identifier      = "${aws_ecr_repository.web.repository_url}:latest"
      image_repository_type = "ECR"

      image_configuration {
        port = "3000"

        runtime_environment_variables = {
          NEXT_PUBLIC_API_URL = "https://${aws_apprunner_service.backend.service_url}"
          NEXT_PUBLIC_WS_URL  = "wss://${aws_apprunner_service.backend.service_url}/ws"
        }
      }
    }

    auto_deployments_enabled = false
  }

  instance_configuration {
    cpu    = "0.25 vCPU"
    memory = "1 GB"
  }

  health_check_configuration {
    protocol            = "TCP"
    interval            = 10
    timeout             = 5
    healthy_threshold   = 1
    unhealthy_threshold = 10
  }

  tags = { Name = "mockstarket-web" }

  depends_on = [aws_iam_role_policy_attachment.apprunner_ecr]
}

# ---------- Outputs ----------

output "api_url" {
  value       = "https://${aws_apprunner_service.backend.service_url}"
  description = "Backend API URL"
}

output "web_url" {
  value       = "https://${aws_apprunner_service.web.service_url}"
  description = "Web frontend URL"
}

output "ecr_backend" {
  value       = aws_ecr_repository.backend.repository_url
  description = "ECR repository URL for backend"
}

output "ecr_web" {
  value       = aws_ecr_repository.web.repository_url
  description = "ECR repository URL for web frontend"
}

output "db_endpoint" {
  value       = aws_db_instance.postgres.endpoint
  description = "RDS PostgreSQL endpoint"
}
