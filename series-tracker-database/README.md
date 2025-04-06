# Proyecto Docker + Base de Datos

Este proyecto utiliza Docker para levantar un contenedor con una base de datos inicializada automáticamente mediante un script SQL (`init.sql`).

## Requisitos

- [Docker](https://www.docker.com/) instalado en tu sistema.

## Estructura del Proyecto

. ├── Dockerfile # Define la imagen del contenedor. ├── init.sql # Script SQL para inicializar la base de datos. ├── README.md # Este archivo.


## ¿Qué hace este proyecto?

1. Construye una imagen de Docker personalizada.
2. Levanta un contenedor con una base de datos (por ejemplo, PostgreSQL o MySQL).
3. Ejecuta automáticamente el script `init.sql` al iniciar el contenedor para crear las tablas y datos iniciales.

## Cómo usarlo

### 1. Construir la imagen

Primero, debes construir la imagen de Docker usando el `Dockerfile` que se encuentra en el proyecto. Para hacerlo, abre tu terminal y navega a la carpeta del proyecto. Luego ejecuta el siguiente comando:

```bash
docker build -t nombre-de-tu-imagen .

