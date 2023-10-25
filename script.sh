docker image build -f Dockerfile -t forum .

docker container run -p 8080:8080 --detach --name forum forum

powershell.exe /c start https://localhost:443

# firefox --new-window --full-sreen http://localhost:8080