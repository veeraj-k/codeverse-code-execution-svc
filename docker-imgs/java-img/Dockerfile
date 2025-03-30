FROM openjdk:17-alpine

WORKDIR /app/code

COPY ./Evaluator.py /app/Evaluator.py


RUN apk update && apk add --no-cache python3 wget
RUN wget https://repo1.maven.org/maven2/org/json/json/20250107/json-20250107.jar -P /app/lib

RUN adduser -D -u 1000 nonrootuser
USER nonrootuser


CMD ["python3", "/app/Evaluator.py" , "java"]