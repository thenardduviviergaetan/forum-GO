docker image build -f Dockerfile -t forum .

docker container run -p 8080:8080 --detach --name forum forum

powershell.exe /c start http://localhost:8080

# firefox --new-window --full-sreen http://localhost:8080