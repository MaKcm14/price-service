FROM python:alpine3.21

LABEL maintainer="maksimacx50@gmail.com"

WORKDIR /by-pass-service

COPY . .

RUN pip install --upgrade pip

RUN pip install flask requests curl_cffi simplejson

ENTRYPOINT [ "python", "server.py" ]
