#!/bin/bash
# setup.sh - Crea el entorno e instala dependencias

python3 -m venv yolo-env
source yolo-env/bin/activate
pip install --upgrade pip
pip install -r requirements.txt

echo "Entorno listo. Activalo con: source yolo-env/bin/activate"
