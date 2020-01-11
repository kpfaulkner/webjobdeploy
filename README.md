Used to upload binaries + associated files as an Azure Webjob (to an existing app service).

Mainly developing this so I can streamline my Go webjob deployment process. Since we (Go developers) don't have Azure Functions
yet, we still need to make use of webjobs. Better than nothing :)

Example of usage:

Will zip and upload from c:\temp\webjob:
.\cmd.exe -username "myusername" -password "mypassword" -appServiceName "myappserviceplanname" -webjobName "mywebjob1" -webjobExeName "dummywebjob.exe" -uploadpath "c:\temp\webjob"

Assuming you've already zipped it:
.\cmd.exe -username "myusername" -password "mypassword" -appServiceName "myappserviceplanname" -webjobName "mywebjob1" -webjobExeName "dummywebjob.exe" -zipFileName "c:\temp\mywebjob.zip"
