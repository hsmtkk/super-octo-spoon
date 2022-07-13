build-lambda-zip -output main.zip main
aws lambda update-function-code --function-name fanout2 --zip-file fileb://main.zip
