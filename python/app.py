import base64
import os
import vertexai
from vertexai.generative_models import GenerativeModel
import vertexai.preview.generative_models as generative_models
from flask import Flask, request, jsonify

app = Flask(__name__)

gcp_credentials = None

# Use custom AWS Security Credentials Supplier
if os.getenv("AWS_EXECUTION_ENV") == "AWS_ECS_FARGATE":
  from gcp_aws_credentials import credentials
  gcp_credentials = credentials

vertexai.init(project=os.getenv("GCP_PROJECT_ID"), location="us-east1", credentials=gcp_credentials)

model = GenerativeModel(
    "gemini-1.5-flash-001",
)

@app.route('/generate', methods=['POST'])
def generate_endpoint():
  """Generates text with the provided prompt."""
  data = request.get_json()
  prompt = data.get('prompt')

  if not prompt:
    return jsonify({'error': 'Missing "prompt" in request body'}), 400

  response = model.generate_content(
      [prompt],
      stream=False,
  )

  return jsonify({'text': response.text})

if __name__ == "__main__":
  port = int(os.environ.get("PORT", 8080))
  print(f"Starting app on port {port}")
  print("Example curl")
  print(f"""curl localhost:{port}/generate -X POST -H 'Content-Type: application/json' -d '{{"prompt": "Write a short story about a cat."}}""")
  app.run(host='0.0.0.0', port=port)
