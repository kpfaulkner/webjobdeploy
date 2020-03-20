Used to upload binaries + associated files as an Azure Webjob (to an existing app service).

Mainly developing this so I can streamline my Go webjob deployment process. Since we (Go developers) don't have Azure Functions
yet, we still need to make use of webjobs. Better than nothing :)

IMPORTANT NOTE:

There is an option (-store) to store the details (app service, username, PASSWORD etc) locally in your home directory (specifically ~/.webjobdeploy/config.json). This is to allow easy 
repeating of deployments. The credentials are NOT encrypted/hashed etc. Currently the assumption is that if your local machine is compromised you've got bigger issues that some deployment creds. 
Naive opinion I agree. Will probably add some sort of encryption, but unsure of immediate approach.

Example of usage:

Will zip and upload from c:\temp\webjob:
.\cmd.exe -username "myusername" -password "mypassword" -appServiceName "myappserviceplanname" -webjobName "mywebjob1" -webjobExeName "dummywebjob.exe" -uploadpath "c:\temp\webjob"

Assuming you've already zipped it:
.\cmd.exe -username "myusername" -password "mypassword" -appServiceName "myappserviceplanname" -webjobName "mywebjob1" -webjobExeName "dummywebjob.exe" -zipFileName "c:\temp\mywebjob.zip"

Username/Password used here are deployment credentials and NOT the ones used for the Azure Portal. For creating credentials please see https://docs.microsoft.com/en-us/azure/app-service/deploy-configure-credentials#userscope
