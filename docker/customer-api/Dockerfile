FROM mcr.microsoft.com/windows/servercore:ltsc2022

SHELL ["powershell"]

RUN mkdir c:\\customer-api
COPY customerAPI.exe c:/customer-api

WORKDIR c:/customer-api
ENTRYPOINT []

CMD c:/customer-api/customerAPI.exe

