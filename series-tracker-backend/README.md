# Series Tracker Backend API

[![Go Report Card](https://goreportcard.com/badge/github.com/Jonialen/Lab-6-Backend-only)](https://goreportcard.com/report/github.com/Jonialen/Lab-6-Backend-only) [![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

API RESTful desarrollada en Go para gestionar y realizar un seguimiento del progreso de visualización de series de televisión. Utiliza el router [Chi](https://github.com/go-chi/chi), el ORM [GORM](https://gorm.io/) para la interacción con la base de datos MySQL, y [Swaggo](https://github.com/swaggo/swag) para la generación de documentación OpenAPI (Swagger).

## ✨ Características Principales

* **Gestión CRUD de Series:** Crear, Leer (todas y por ID), Actualizar y Eliminar series.
* **Seguimiento de Progreso:**
    * Actualizar el estado de visualización (`Plan to Watch`, `Watching`, `Completed`, `Dropped`).
    * Registrar/incrementar el último episodio visto.
* **Ranking:** Sistema simple de votación (upvote/downvote) para las series.
* **API RESTful:** Diseño siguiendo principios REST.
* **Documentación Interactiva:** Endpoints documentados con Swagger UI.
* **Configuración Flexible:** Uso de variables de entorno para la configuración de la base de datos.
* **Middleware:** Incluye logging, recuperación de panics, CORS, etc., gracias a Chi middleware.
* **Cierre Grácil (Graceful Shutdown):** Implementado para permitir que las solicitudes en curso finalicen antes de apagar el servidor.
* **Contenerización:** Preparado para ejecutarse en un contenedor Docker (requiere base de datos externa).

## 🛠️ Tecnologías Utilizadas

* **Lenguaje:** [Go](https://golang.org/) (v1.18+)
* **Router HTTP:** [Chi (v5)](https://github.com/go-chi/chi)
* **ORM:** [GORM](https://gorm.io/)
* **Driver Base de Datos:** [MySQL Driver para GORM](https://gorm.io/docs/connecting_to_the_database.html#MySQL)
* **Documentación API:** [Swaggo (Swag & http-swagger)](https://github.com/swaggo/swag)
* **Base de Datos:** MySQL (o compatible)
* **Contenerización:** [Docker](https://www.docker.com/)

## 📋 Prerrequisitos

* **Go:** Versión 1.18 o superior (para desarrollo local y/o regenerar Swagger). [Instrucciones](https://golang.org/doc/install)
* **Docker:** Para construir y ejecutar la imagen del contenedor. [Instrucciones](https://docs.docker.com/engine/install/)
* **MySQL:** Una instancia de base de datos MySQL en ejecución (puede ser local, en otro contenedor Docker, o un servicio gestionado). La base de datos especificada debe existir.
* **Swag CLI:** (Opcional, solo si necesitas regenerar la documentación Swagger)
    ```bash
    go install [github.com/swaggo/swag/cmd/swag@latest](https://www.google.com/search?q=https://github.com/swaggo/swag/cmd/swag%40latest)
    ```

## 🚀 Instalación y Configuración (Local)

1.  **Clonar el repositorio:**
    ```bash
    git clone [https://github.com/Jonialen/Lab-6-Backend-only.git](https://github.com/Jonialen/Lab-6-Backend-only.git)
    cd Lab-6-Backend-only
    ```

2.  **Instalar dependencias (para desarrollo local):**
    ```bash
    go mod tidy
    ```

3.  **Configurar Variables de Entorno (para desarrollo local):**
    Crea un archivo `.env` en la raíz del proyecto o exporta las variables directamente en tu terminal:
    ```dotenv
    # .env (Ejemplo para desarrollo local)
    DB_HOST=localhost      # O 127.0.0.1 si MySQL corre localmente
    DB_PORT=3306
    DB_USER=app_user       # Usuario de tu base de datos
    DB_PASSWORD=app_password # Contraseña de tu base de datos
    DB_NAME=anime_db       # Nombre de tu base de datos (debe existir)
    PORT=8080              # Puerto para la API
    ```
    *Asegúrate de que la base de datos (`DB_NAME`) exista en tu instancia MySQL.* GORM (`AutoMigrate`) creará la tabla `series` si no existe.

## ▶️ Ejecutar la Aplicación (Localmente)

1.  Asegúrate de que tu instancia MySQL esté corriendo y accesible con las credenciales configuradas.
2.  Desde la raíz del proyecto, ejecuta el servidor:
    ```bash
    go run main.go
    ```
    El servidor debería iniciarse y mostrar logs indicando que está escuchando en el puerto configurado (por defecto `8080`). La API estará disponible en `http://localhost:8080`.

## 🐳 Ejecutar con Docker

Este proyecto incluye un `Dockerfile` para construir una imagen de contenedor para la aplicación Go. **Nota:** Esta configuración asume que la base de datos MySQL se ejecuta externamente al contenedor de la aplicación.

1.  **Asegúrate de que Docker esté instalado y corriendo.**

2.  **Asegúrate de que una instancia de MySQL esté accesible** desde donde correrás el contenedor Docker. Podría ser:
    * Una base de datos local (podrías necesitar usar `host.docker.internal` como `DB_HOST` en Linux/Mac o la IP de tu máquina en Windows).
    * Otra base de datos en un contenedor Docker (considera usar `docker-compose` para esto).
    * Un servicio de base de datos en la nube o en red.

3.  **Construir la imagen Docker:**
    Desde la raíz del proyecto (donde está el `Dockerfile`):
    ```bash
    docker build -t jonialen/series-tracker-backend:latest .
    ```
    (Puedes cambiar `jonialen/series-tracker-backend:latest` por el nombre de imagen y tag que prefieras).

4.  **Ejecutar el contenedor:**

5.  **Verificar:**
    La API debería estar accesible en `http://<tu_ip_o_localhost>:8080`. Puedes ver los logs con `docker logs series-tracker-app` (si no usaste `--rm` y el contenedor sigue corriendo en segundo plano con `-d`).

**Nota sobre Docker Compose:**
Para simplificar el despliegue local (incluyendo la base de datos), considera crear un archivo `docker-compose.yml`. Esto te permite definir y gestionar la aplicación y la base de datos como servicios interconectados con un solo comando (`docker-compose up`).

## 📚 Documentación de la API (Swagger)

Este proyecto utiliza `swaggo` para generar automáticamente la documentación OpenAPI (Swagger) a partir de los comentarios en el código fuente (`godoc`). Los archivos generados (`docs/`) están incluidos en el repositorio y son utilizados por la aplicación.

* **Acceder a la UI de Swagger:** Una vez que la aplicación (local o en Docker) esté corriendo, puedes acceder a la documentación interactiva en tu navegador visitando:
    [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) (Ajusta el puerto si usaste uno diferente o si accedes a Docker desde una IP distinta).
    La UI de Swagger te permite ver todos los endpoints, sus parámetros, cuerpos de solicitud/respuesta y probarlos directamente.

* **Regenerar Documentación:**
    Si realizas cambios en los comentarios de documentación (`// @Summary`, `// @Param`, etc.) en los archivos Go (`handlers`, `models`, `main.go`), necesitas regenerar los archivos de Swagger para que los cambios se reflejen en la UI. Ejecuta el siguiente comando desde la raíz del proyecto:
    ```bash
    swag init
    ```
    Asegúrate de hacer commit de los archivos actualizados en el directorio `docs/` si los cambios son permanentes. Si estás usando Docker, reconstruye la imagen después de regenerar los documentos.

## API Endpoints Principales

La API sigue un diseño RESTful con el prefijo base `/api`.

* `GET    /api/series`: Obtiene una lista de todas las series.
* `POST   /api/series`: Crea una nueva serie.
* `GET    /api/series/{id}`: Obtiene los detalles de una serie específica por su ID.
* `PUT    /api/series/{id}`: Actualiza completamente una serie existente por su ID.
* `DELETE /api/series/{id}`: Elimina una serie por su ID.
* `PATCH  /api/series/{id}/status`: Actualiza parcialmente el estado (`status`) de una serie.
* `PATCH  /api/series/{id}/episode`: Incrementa el contador de episodios vistos (`last_episode_watched`) de una serie.
* `PATCH  /api/series/{id}/upvote`: Incrementa el ranking (`ranking`) de una serie.
* `PATCH  /api/series/{id}/downvote`: Decrementa el ranking (`ranking`) de una serie.
* `GET    /health`: Endpoint simple para verificar si la API está en funcionamiento. Devuelve `{"status": "ok"}`.

*Para detalles completos sobre los parámetros de ruta, query params, cuerpos de solicitud JSON y códigos de respuesta, por favor consulta la documentación interactiva de Swagger.*

