# syntax=docker/dockerfile:1.4
FROM python:3.10-slim
WORKDIR /app
RUN apt-get update \
    && apt-get install -y curl \
    && rm -rf /var/lib/apt/lists/*
COPY requirements.txt /app/
RUN --mount=type=cache,target=/root/.cache/pip \
    pip3 install -r requirements.txt
EXPOSE 8501
HEALTHCHECK CMD curl --fail http://localhost:8501/_stcore/health
ENTRYPOINT ["streamlit", "run", "app.py", "--server.port=8501", "--server.address=0.0.0.0"]
