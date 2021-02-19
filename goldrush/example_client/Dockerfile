FROM python:latest

COPY ./setup.py setup.py


RUN python setup.py develop

COPY . /app
WORKDIR /app
ENTRYPOINT ["python", "main.py"]
