resource "aws_api_gateway_rest_api" "rag_diary_gateway" {
  name = "rag_diary_gateway"
  description = "API Gateway for Rag Diary"
}

# Resource for API Gateway /api endpoint
resource "aws_api_gateway_resource" "api" {
  rest_api_id = aws_api_gateway_rest_api.rag_diary_gateway.id
  parent_id = aws_api_gateway_rest_api.rag_diary_gateway.root_resource_id
  path_part = "api"
}

# Resource for API Gateway /api/v1 endpoint
resource "aws_api_gateway_resource" "v1" {
  rest_api_id = aws_api_gateway_rest_api.rag_diary_gateway.id
  parent_id = aws_api_gateway_resource.api.id
  path_part = "v1"
}

# Resource for API Gateway /api/v1/diary endpoint
resource "aws_api_gateway_resource" "diary" {
  rest_api_id = aws_api_gateway_rest_api.rag_diary_gateway.id
  parent_id = aws_api_gateway_resource.v1.id
  path_part = "diary"
}

# Resource for API Gateway /api/v1/diary/semantic-search endpoint
resource "aws_api_gateway_resource" "semantic_search" {
  rest_api_id = aws_api_gateway_rest_api.rag_diary_gateway.id
  parent_id = aws_api_gateway_resource.diary.id
  path_part = "semantic-search"
}

# Resource for API Gateway /api/v1/diary/rag-response endpoint
resource "aws_api_gateway_resource" "rag_response" {
  rest_api_id = aws_api_gateway_rest_api.rag_diary_gateway.id
  parent_id = aws_api_gateway_resource.diary.id
  path_part = "rag-response"
}

# Method for POST /api/v1/diary endpoint CREATE Diary entry
resource "aws_api_gateway_method" "post_diary" {
  rest_api_id = aws_api_gateway_rest_api.rag_diary_gateway.id
  resource_id = aws_api_gateway_resource.diary.id
  http_method = "POST"
  authorization = "NONE"
}

# Method for POST /api/v1/diary/semantic-search endpoint SEARCH Diary entry
resource "aws_api_gateway_method" "post_semantic_search" {
  rest_api_id = aws_api_gateway_rest_api.rag_diary_gateway.id
  resource_id = aws_api_gateway_resource.semantic_search.id
  http_method = "POST"
  authorization = "NONE"
}

# Method for POST /api/v1/diary/rag-response endpoint CREATE Rag Response
resource "aws_api_gateway_method" "post_rag_response" {
  rest_api_id = aws_api_gateway_rest_api.rag_diary_gateway.id
  resource_id = aws_api_gateway_resource.rag_response.id
  http_method = "POST"
  authorization = "NONE"
}

# Integration for POST /api/v1/diary endpoint
resource "aws_api_gateway_integration" "post_diary_integration" {
  rest_api_id = aws_api_gateway_rest_api.rag_diary_gateway.id
  resource_id = aws_api_gateway_resource.diary.id
  http_method = aws_api_gateway_method.post_diary.http_method

  integration_http_method = "POST"
  type = "AWS_PROXY"
  uri = aws_lambda_function.rag_diary.invoke_arn
}

# Integration for POST /api/v1/diary/semantic-search endpoint
resource "aws_api_gateway_integration" "post_semantic_search_integration" {
  rest_api_id = aws_api_gateway_rest_api.rag_diary_gateway.id
  resource_id = aws_api_gateway_resource.semantic_search.id
  http_method = aws_api_gateway_method.post_semantic_search.http_method

  integration_http_method = "POST"
  type = "AWS_PROXY"
  uri = aws_lambda_function.rag_diary.invoke_arn
}

# Integration for POST /api/v1/diary/rag-response endpoint
resource "aws_api_gateway_integration" "post_rag_response_integration" {
  rest_api_id = aws_api_gateway_rest_api.rag_diary_gateway.id
  resource_id = aws_api_gateway_resource.rag_response.id
  http_method = aws_api_gateway_method.post_rag_response.http_method

  integration_http_method = "POST"
  type = "AWS_PROXY"
  uri = aws_lambda_function.rag_diary.invoke_arn
}

# Method Response for POST /api/v1/diary/rag-response endpoint
resource "aws_api_gateway_method_response" "post_rag_response_method_response" {
  rest_api_id = aws_api_gateway_rest_api.rag_diary_gateway.id
  resource_id = aws_api_gateway_resource.rag_response.id
  http_method = aws_api_gateway_method.post_rag_response.http_method
  status_code = "200"

  response_parameters = {
    "method.response.header.Access-Control-Allow-Origin" = true
    "method.response.header.Access-Control-Allow-Methods" = true
    "method.response.header.Access-Control-Allow-Headers" = true
  }
}

# Integration Response for POST /api/v1/diary/rag-response endpoint
resource "aws_api_gateway_integration_response" "post_rag_response_integration_response" {
  rest_api_id = aws_api_gateway_rest_api.rag_diary_gateway.id
  resource_id = aws_api_gateway_resource.rag_response.id
  http_method = aws_api_gateway_method.post_rag_response.http_method
  status_code = aws_api_gateway_method_response.post_rag_response_method_response.status_code

  response_parameters = {
    "method.response.header.Access-Control-Allow-Origin" = "'https://dieg0code.site'"
    "method.response.header.Access-Control-Allow-Methods" = "'POST'"
    "method.response.header.Access-Control-Allow-Headers" = "'Content-Type,Authorization'"
  }
}

# Invoke permission for API Gateway to invoke Lambda
resource "aws_lambda_permission" "api_gateway_invoke_lambda" {
  statement_id = "AllowAPIGatewayInvoke"
  action = "lambda:InvokeFunction"
  function_name = aws_lambda_function.rag_diary.function_name
  principal = "apigateway.amazonaws.com"
  source_arn = "${aws_api_gateway_rest_api.rag_diary_gateway.execution_arn}/*/*/*"
}

# Gateway deployment
resource "aws_api_gateway_deployment" "api_gateway_deployment" {
  depends_on = [ 
    aws_api_gateway_integration.post_diary_integration,
    aws_api_gateway_integration.post_semantic_search_integration,
    aws_api_gateway_integration.post_rag_response_integration
   ]

   rest_api_id = aws_api_gateway_rest_api.rag_diary_gateway.id
   description = "API Gateway Rag Diary Deployment"

   triggers = {
    redeployment = sha1(jsonencode({
      post_diary_integration = aws_api_gateway_integration.post_diary_integration.id,
      post_semantic_search_integration = aws_api_gateway_integration.post_semantic_search_integration.id,
      post_rag_response_integration = aws_api_gateway_integration.post_rag_response_integration.id
    }))
   }

    lifecycle {
     create_before_destroy = true
    }
}

# Gateway stage
resource "aws_api_gateway_stage" "api_gateway_stage" {
  deployment_id = aws_api_gateway_deployment.api_gateway_deployment.id
  rest_api_id = aws_api_gateway_rest_api.rag_diary_gateway.id
  stage_name = "dev"

  depends_on = [ 
    aws_api_gateway_deployment.api_gateway_deployment
   ]
}

# Ivocation URL
output "api_gateway_url" {
  value = "https://${aws_api_gateway_rest_api.rag_diary_gateway.id}.execute-api.sa-east-1.amazonaws.com/${aws_api_gateway_stage.api_gateway_stage.stage_name}/api/v1/diary"
}