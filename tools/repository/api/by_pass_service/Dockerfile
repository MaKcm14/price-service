FROM python:alpine3.21

WORKDIR /by-pass-service

COPY . .

RUN pip install --upgrade pip

RUN pip install flask requests curl_cffi simplejson

EXPOSE 9090

ENTRYPOINT [ "python", "server.py" ]
