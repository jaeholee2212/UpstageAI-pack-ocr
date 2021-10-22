FROM python

WORKDIR /src
RUN apt-get update && apt-get install -y python3-pip iputils-ping iproute2
RUN pip install fastapi &&  \
  pip install "uvicorn[standard]"

CMD uvicorn main:app --reload --host 0.0.0.0