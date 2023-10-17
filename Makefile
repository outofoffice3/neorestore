include .env

list-vars: 
	@echo "access key : $(accesKey)"
	@echo "secret access key : $(secretAccessKey)"
	@echo "region : $(region)"
	@echo "binary name : $(binaryName)"
	@echo "register prefix lambda source file : $(registerPrefixSourceFile)"
	@echo "register prefix lambda role arn : $(registerPrefixRoleArn)"
	@echo "register prefix lambda env vars : $(registerPrefixEnvVars)"
	@echo "register prefix lambda name : $(registerPrefixLambdaName)"