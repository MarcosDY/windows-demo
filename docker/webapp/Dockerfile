FROM mcr.microsoft.com/windows/servercore:ltsc2022

SHELL ["powershell"]

RUN mkdir c:\\webapp
COPY webapp.exe c:/webapp/

COPY startup-script.ps1 c:/startup-script.ps1
COPY fix-hosts.ps1 c:/fix-hosts.ps1

WORKDIR c:/webapp
ENTRYPOINT c:/startup-script.ps1

