import base64
import os
import vertexai
from vertexai.generative_models import GenerativeModel, Part, SafetySetting, FinishReason
import vertexai.preview.generative_models as generative_models
from gcp_aws_credentials import credentials


def generate():
  """Show how the credentials are used"""
  vertexai.init(project=os.getenv("GCP_PROJECT_ID"), location="us-east1", credentials=credentials)
  model = GenerativeModel(
    "gemini-1.5-flash-001",
  )
  response = model.generate_content(
      ["""tell me a funny running related joke"""],
      stream=False,
  )

  print(response.text)


generate()