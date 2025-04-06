# Aplicación de Seguimiento de Series (Lab 6)

## Descripción

Este proyecto es una aplicación web full-stack diseñada para ayudar a los usuarios a realizar un seguimiento de las series de TV que están viendo. Permite a los usuarios agregar nuevas series, actualizar su progreso de visualización (episodios vistos), cambiar el estado (por ejemplo, Viendo, Completada, Plan para Ver), clasificar series y ver su lista de seguimiento en una tabla filtrable.

La aplicación está contenedorizada utilizando Docker y orquestada con Docker Compose para una configuración y despliegue fáciles.

## Arquitectura

La aplicación sigue una arquitectura cliente-servidor estándar con una base de datos separada:

1. **Frontend (`series-tracker`)**: Una aplicación de una sola página (SPA) construida con JavaScript vanilla, HTML y CSS. Proporciona la interfaz de usuario para interactuar con los datos de las series.
2. **Backend (`series-tracker-backend`)**: Una API RESTful construida con Go (Golang). Maneja la lógica de negocio, interactúa con la base de datos y sirve datos al frontend. Incluye documentación Swagger para los endpoints de la API.
3. **Base de datos (`series-tracker-database`)**: Una base de datos MySQL (MariaDB) utilizada para almacenar toda la información de las series. Se proporciona un script `init.sql` para configurar la(s) tabla(s) necesaria(s) al inicializar.

## Stack Tecnológico

* **Frontend**: HTML, CSS, JavaScript (Vanilla JS)
* **Backend**: Go (Golang)
* **Base de datos**: MySQL (MariaDB)
* **Contenedorización**: Docker, Docker Compose

## Prerrequisitos

Antes de comenzar, asegúrate de tener instalado lo siguiente:

* [Docker](https://docs.docker.com/get-docker/)
* [Docker Compose](https://docs.docker.com/compose/install/) (Generalmente incluido con Docker Desktop)

## Empezando

1. **Clonar el Repositorio**:
    ```bash
    git clone https://github.com/Jonialen/Lab-6-Backend-only.git
    cd Lab-6-Backend-only
    ```

2. **Construir y Ejecutar con Docker Compose**:
    Navega al directorio raíz del proyecto (`Lab-6-Backend-only/`) donde se encuentra el archivo `docker-compose.yml` y ejecuta:
    ```bash
    docker-compose up --build -d
    ```
    * `--build`: Fuerza a Docker Compose a construir las imágenes (útil para la primera ejecución o después de cambios en el código).
    * `-d`: Ejecuta los contenedores en modo desconectado (en segundo plano).

    Este comando hará lo siguiente:
    * Construirá las imágenes Docker para los servicios de frontend y backend (si no están ya construidas).
    * Descargará la imagen oficial de PostgreSQL.
    * Creará y arrancará contenedores para los servicios de frontend, backend y base de datos.
    * Configurará la red necesaria para que los contenedores se comuniquen.
    * Inicializará la base de datos utilizando el script `init.sql`.

3. **Detener la Aplicación**:
    Para detener los contenedores en ejecución, ejecuta el siguiente comando en el mismo directorio:
    ```bash
    docker-compose down
    ```
    Para detener y eliminar volúmenes (como los datos de la base de datos), usa:
    ```bash
    docker-compose down -v
    ```

## Uso

Una vez que los contenedores estén en funcionamiento:

* **Aplicación Frontend**: Accede a la interfaz web en tu navegador en:
    `http://localhost:{puerto eleccion}`
    *(Este puerto está definido en el archivo `docker-compose.yml`)*

* **API Backend**: Los endpoints de la API se sirven en:
    `http://localhost:8080`
    *(Este puerto está definido en el archivo `docker-compose.yml`)*

* **Documentación de la API (Swagger)**: Puedes explorar los endpoints de la API utilizando la interfaz de usuario de Swagger disponible en:
    `http://localhost:5000/swagger/index.html`

## Estructura del Proyecto

Lab-6-Backend-only/
├── docker-compose.yml        # Configuración de Docker Compose para todos los servicios
├── series-tracker/           # Aplicación frontend (HTML, CSS, JS)
│   ├── components/           # Componentes de UI reutilizables (JS)
│   ├── pages/                # Lógica específica de páginas (JS)
│   ├── utils/                # Funciones utilitarias (llamadas a la API, renderizado, etc.)
│   ├── index.html            # Archivo HTML principal
│   ├── main.js               # Punto de entrada principal para JS del frontend
│   ├── styles.css            # Estilos CSS
│   └── README.md             # Detalles específicos del frontend
├── series-tracker-backend/   # API Backend (Go)
│   ├── docs/                 # Archivos de documentación Swagger
│   ├── handlers/             # Manejadores de solicitudes HTTP
│   ├── models/               # Definiciones de estructuras de datos (structs)
│   ├── repository/           # Lógica de interacción con la base de datos
│   ├── main.go               # Punto de entrada principal para la API backend
│   └── README.md             # Detalles específicos del backend
├── series-tracker-database/  # Configuración de la base de datos
│   ├── init.sql              # Script de inicialización de la base de datos (crea tablas)
│   └── README.md             # Detalles específicos de la base de datos
└── README.md                 # Este archivo (Visión General del Proyecto)


## Información Adicional

Para obtener información más detallada sobre cada componente, consulta sus respectivos archivos README:

* [README del Frontend](./series-tracker/README.md)
* [README del Backend](./series-tracker-backend/README.md)
* [README de la Base de Datos](./series-tracker-database/README.md)

Para usar esto:

1. Crea un archivo llamado `README.md` directamente dentro de la carpeta `Lab-6-Backend-only` (al mismo nivel que `docker-compose.yml`, `series-tracker/`, `series-tracker-backend/`, etc.).
2. Copia y pega el contenido de arriba en ese archivo.
3. Revisa y ajusta cualquier detalle (como el marcador de posición de la URL del repositorio) si es necesario.

