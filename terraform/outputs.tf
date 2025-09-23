output "instance_id" {
  description = "The ID of the EC2 instance."
  value       = aws_instance.go_api_instance.id
}

output "security_group_id" {
  description = "The ID of the Security Group."
  value       = aws_security_group.app_sg.id
}