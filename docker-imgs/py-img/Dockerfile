FROM python:3.9-slim


RUN useradd -m nonrootuser

COPY ./Evaluator.py /app/Evaluator.py
USER nonrootuser

WORKDIR /app

CMD ["python3", "/app/Evaluator.py" , "python"]
