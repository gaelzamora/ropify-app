from fastapi import FastAPI, UploadFile, File, HTTPException
from fastapi.responses import Response, JSONResponse
from fastapi.middleware.cors import CORSMiddleware
import logging
import time
import os
from rembg import remove, new_session

# Configurar logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Crear aplicación FastAPI
app = FastAPI(
    title="Background Removal Service",
    description="Servicio para eliminar fondos de imágenes",
    version="1.0.0"
)

# Configurar CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Crear una sesión global para el modelo
try:
    logger.info("Inicializando modelo rembg...")
    session = new_session("u2net")
    logger.info("Modelo rembg inicializado correctamente")
except Exception as e:
    logger.error(f"Error inicializando modelo: {e}")
    session = None

@app.get("/health")
def health():
    """Verificar que el servicio está funcionando y el modelo cargado"""
    return {
        "status": "ok",
        "model_loaded": session is not None,
        "version": "1.0.0"
    }

@app.post("/remove-bg")
async def remove_bg(file: UploadFile = File(...)):
    """Elimina el fondo de una imagen"""
    if session is None:
        logger.error("Modelo no inicializado")
        raise HTTPException(status_code=503, detail="Servicio no disponible: modelo no inicializado")
    
    logger.info(f"Procesando imagen: {file.filename}, content_type: {file.content_type}")
    
    try:
        # Leer datos de imagen
        image_data = await file.read()
        logger.info(f"Imagen leída: {len(image_data)} bytes")
        
        # Procesar imagen
        output = remove(image_data, session=session)
        logger.info("Procesamiento completado")
        
        # Devolver imagen con fondo eliminado
        return Response(content=output, media_type="image/png")
    
    except Exception as e:
        logger.error(f"Error procesando imagen: {str(e)}")
        return JSONResponse(
            status_code=500,
            content={"error": f"Error al procesar la imagen: {str(e)}"}
        )