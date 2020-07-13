FROM python:3.8

#COPY ./.docker/app/dev-entrypoint.sh /dev-entrypoint.sh
#RUN chmod +x /dev-entrypoint.sh

#COPY ./src/go.mod ./
#COPY ./src/go.sum ./
#RUN go mod download

COPY ./src /
RUN pip install aiohttp

CMD [ "python", "/main.py" ]