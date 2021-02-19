FROM python:latest

COPY ./requirements.txt requirements.txt
RUN pip install -U -r requirements.txt

COPY ./setup.py setup.py


RUN python setup.py develop

COPY . /app
WORKDIR /app
ENTRYPOINT ["python", "main.py"]
