docker image build -f Dockerfile -t forumproject .
docker container run -p 8080:8080 --detach --name forum forumproject
docker ps -a
echo To stop container, run "Docker stop <container_name>".