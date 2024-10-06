variable "supabase_url" {
  description = "URL for Supabase"
  type        = string
  sensitive   = true
}

variable "supabase_key" {
  description = "API key for Supabase"
  type        = string
  sensitive   = true
}

variable "openai_api_key" {
  description = "API key for OpenAI"
  type        = string
  sensitive   = true
}

variable "ip_info_api_key" {
  description = "API key for IP Info"
  type        = string
  sensitive   = true
}
