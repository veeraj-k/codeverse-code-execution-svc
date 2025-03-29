import json
import traceback
import subprocess
import io
import sys
from contextlib import redirect_stdout
import os

with open('/app/code/test_cases.json', 'r') as f:
    test_cases = json.load(f)

if len(sys.argv) > 1:
    language = sys.argv[1]
else:
    raise Exception("Language not specified. Please provide the language as a command-line argument.")

results = []

for i, test_case in enumerate(test_cases):
    test_result = {
        "status": None,
        "input": test_case['input'],
        "actual": test_case['output']['output'],
        "output": None,
        "stdout": None,
        "error": None
    }

    try:
        inputs = test_case['input']
        expected = test_case['output']['output']
        result = None

        if language == "python":

            result = subprocess.run(
                ["python", "/app/code/Runner.py", json.dumps(test_case)],
                capture_output=True, text=True
            )
            print(result)
           

        elif language == "java":
            print("Enterred the Java block")
            sb = subprocess.run(["javac","-cp","/app/lib/json-20250107.jar:.", "/app/code/Runner.java"], check=True)
            print(sb)
            result = subprocess.run(
                ["java", "-Xms64m", "-Xmx128m","-cp","/app/lib/json-20250107.jar:.", "Runner" , json.dumps(test_case)]  ,
                capture_output=True, text=True
            )
            print(result)


        elif language == "c":
            subprocess.run(["gcc", "/app/code/Runner.c", "/app/code/Solution.c", "-o", "runner", "-lcjson"], check=True)
            result = subprocess.run(
                ["./runner" , json.dumps(test_case)],
                capture_output=True, text=True
            )
            print(result)
        else:
            raise Exception(f"Unsupported language: {language}")


        raw_output = result.stdout.strip()
        lines = raw_output.split("\n")
        output = lines[-1] if lines else ""  
        captured_stdout = "\n".join(lines[:-1])  
        if result.stderr:
            raise Exception(f"Error occurred during execution: {result.stderr.strip()}")

        try:
            test_result["output"] = json.loads(output)
        except json.JSONDecodeError:
            test_result["output"] = output
        test_result["stdout"] = captured_stdout
        test_result["status"] = "passed" if output == str(expected) else "failed"

    except Exception as e:
        test_result["status"] = "error"
        test_result["error"] = str(e)

    results.append(test_result)

# Write results to output.json
with open('/app/code/output.json', 'w') as f:
    json.dump(results, f, indent=2)
