# filepath: background-removal-service/main.py
from fastapi import FastAPI, File, UploadFile
from fastapi.responses import Response
from rembg import remove

app = FastAPI()

@app.post("/remove-background")
async def remove_background(file: UploadFile = File(...)):
    image_bytes = await file.read()
    print(f"Received file: {file.filename}, size: {len(image_bytes)} bytes")

    output = remove(image_bytes)
    print(f"Output size: {len(output)} bytes")

    return Response(content=output, media_type="image/png")