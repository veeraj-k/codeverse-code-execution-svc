FROM alpine:latest

WORKDIR /app/code
COPY ./Evaluator.py /app/Evaluator.py

RUN apk add --no-cache python3

RUN apk add --no-cache gcc musl-dev make build-base cjson-dev

# Add a non-root user for security
RUN adduser -D -u 1000 nonrootuser
USER nonrootuser


CMD ["python3", "/app/Evaluator.py" , "c"]

