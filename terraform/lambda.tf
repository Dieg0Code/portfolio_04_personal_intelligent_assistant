resource "aws_lambda_function" "rag_diary" {
  filename = "rag_diary.zip"
  function_name = "rag_diary"
  role = aws_iam_role.rag_diary_lambda_role.arn
  handler = "bootstrap"
  runtime = "provided.al2023"
  memory_size = 128
  timeout = 160

  source_code_hash = filebase64sha256("rag_diary.zip")
}