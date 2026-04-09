terraform {
  required_version = ">= 1.5"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.0"
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
  description = "RDS PostgreSQL access"
  vpc_id      = data.aws_vpc.default.id

  ingress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
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
  publicly_accessible    = true
  skip_final_snapshot    = true

  backup_retention_period = 7
  deletion_protection     = false

  tags = { Name = "mockstarket" }
}

# ---------- ECR Repositories ----------

resource "aws_ecr_repository" "backend" {
  name                 = "mockstarket-backend"
  image_tag_mutability = "MUTABLE"
  force_delete         = true

  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "aws_ecr_repository" "web" {
  name                 = "mockstarket-web"
  image_tag_mutability = "MUTABLE"
  force_delete         = true

  image_scanning_configuration {
    scan_on_push = true
  }
}

# ---------- CloudWatch Log Groups ----------

resource "aws_cloudwatch_log_group" "backend" {
  name              = "/aws/ecs/mockstarket-api"
  retention_in_days = 14
}

resource "aws_cloudwatch_log_group" "web" {
  name              = "/aws/ecs/mockstarket-web"
  retention_in_days = 14
}

# ---------- IAM: ECS Task Execution Role ----------

resource "aws_iam_role" "ecs_execution" {
  name = "mockstarket-ecs-execution"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = { Service = "ecs-tasks.amazonaws.com" }
    }]
  })
}

resource "aws_iam_role_policy_attachment" "ecs_execution" {
  role       = aws_iam_role.ecs_execution.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

# ---------- IAM: ECS Express Infrastructure Role ----------

resource "aws_iam_role" "ecs_infrastructure" {
  name = "mockstarket-ecs-infrastructure"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = { Service = "ecs.amazonaws.com" }
    }]
  })
}

resource "aws_iam_role_policy_attachment" "ecs_infrastructure" {
  role       = aws_iam_role.ecs_infrastructure.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonECSInfrastructureRoleforExpressGatewayServices"
}

# ---------- ECS Express Mode: Backend API ----------

resource "aws_ecs_express_gateway_service" "backend" {
  service_name           = "mockstarket-api"
  execution_role_arn     = aws_iam_role.ecs_execution.arn
  infrastructure_role_arn = aws_iam_role.ecs_infrastructure.arn

  cpu    = "256"
  memory = "512"

  health_check_path = "/api/v1/system/health"

  primary_container {
    image          = "${aws_ecr_repository.backend.repository_url}:latest"
    container_port = 8080

    aws_logs_configuration {
      log_group         = aws_cloudwatch_log_group.backend.name
      log_stream_prefix = "ecs"
    }

    environment {
      name  = "PORT"
      value = "8080"
    }
    environment {
      name  = "DATABASE_URL"
      value = "postgres://${aws_db_instance.postgres.username}:${var.db_password}@${aws_db_instance.postgres.endpoint}/${aws_db_instance.postgres.db_name}?sslmode=require"
    }
    environment {
      name  = "MARKET_DATA_SOURCE"
      value = "polygon"
    }
    environment {
      name  = "POLYGON_API_KEY"
      value = var.polygon_api_key
    }
    environment {
      name  = "POLYGON_BASE_URL"
      value = "https://api.polygon.io"
    }
    environment {
      name  = "POLYGON_WS_ENABLED"
      value = "false"
    }
    environment {
      name  = "POLYGON_POLL_INTERVAL_MS"
      value = "30000"
    }
    environment {
      name  = "DEV_MODE"
      value = "true"
    }
    environment {
      name  = "CORS_ORIGINS"
      value = "*"
    }
    environment {
      name  = "LOG_LEVEL"
      value = "info"
    }
    environment {
      name  = "STARTING_CASH"
      value = "100000"
    }
    environment {
      name  = "MAX_WS_CLIENTS"
      value = "1000"
    }
  }

  scaling_target {
    min_task_count            = 1
    max_task_count            = 5
    auto_scaling_metric       = "AVERAGE_CPU"
    auto_scaling_target_value = 70
  }

  tags = { Name = "mockstarket-api" }
}

# ---------- ECS Express Mode: Web Frontend ----------

resource "aws_ecs_express_gateway_service" "web" {
  service_name           = "mockstarket-web"
  execution_role_arn     = aws_iam_role.ecs_execution.arn
  infrastructure_role_arn = aws_iam_role.ecs_infrastructure.arn

  cpu    = "256"
  memory = "512"

  health_check_path = "/"

  primary_container {
    image          = "${aws_ecr_repository.web.repository_url}:latest"
    container_port = 3000

    aws_logs_configuration {
      log_group         = aws_cloudwatch_log_group.web.name
      log_stream_prefix = "ecs"
    }
  }

  scaling_target {
    min_task_count            = 1
    max_task_count            = 3
    auto_scaling_metric       = "AVERAGE_CPU"
    auto_scaling_target_value = 70
  }

  tags = { Name = "mockstarket-web" }
}

# ---------- Outputs ----------

output "api_url" {
  value       = aws_ecs_express_gateway_service.backend.service_url
  description = "Backend API URL (HTTPS)"
}

output "web_url" {
  value       = aws_ecs_express_gateway_service.web.service_url
  description = "Web frontend URL (HTTPS)"
}

output "ws_url" {
  value       = "wss://${replace(aws_ecs_express_gateway_service.backend.service_url, "https://", "")}/ws"
  description = "WebSocket URL for real-time prices"
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
