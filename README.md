# API de Conversión de Monedas

API RESTful para conversión de monedas con soporte para USD, EUR y VES (Bolívar Soberano de Venezuela).

## Tabla de Contenidos

- [Características](#características)
- [Requisitos](#requisitos)
- [Instalación](#instalación)
- [Configuración](#configuración)
- [Uso](#uso)
  - [Estructura del Proyecto](#estructura-del-proyecto)
  - [Endpoints](#endpoints)
    - [Obtener Tasa de Cambio](#obtener-tasa-de-cambio)
    - [Obtener Todas las Tasas](#obtener-todas-las-tasas)
    - [Convertir Moneda](#convertir-moneda)
    - [Health Check](#health-check)
- [Ejemplos de Uso](#ejemplos-de-uso)
- [Formatos de Moneda](#formatos-de-moneda)
- [Manejo de Errores](#manejo-de-errores)
- [Caché](#caché)
- [Pruebas](#pruebas)
- [Despliegue](#despliegue)
- [Contribución](#contribución)
- [Licencia](#licencia)

## Características

- Conversión en tiempo real entre USD, EUR y VES
- Tasa de cambio actualizada automáticamente
- Caché para mejorar el rendimiento
- API RESTful con documentación detallada
- Soporte para múltiples formatos de entrada de moneda
- Sistema de logging integrado

## Requisitos

- Go 1.16 o superior
- Python 3.6 o superior (para el script de obtención de tasas)
- Módulos de Python: `requests`

## Instalación

1. Clonar el repositorio:
   ```bash
   git clone https://github.com/tu-usuario/conversor-monedas.git
   cd conversor-monedas
   ```

2. Instalar dependencias de Python:
   ```bash
   pip install requests
   ```

3. Construir y ejecutar la aplicación:
   ```bash
   go build -o conversor
   ./conversor
   ```

## Configuración

La aplicación se puede configurar mediante variables de entorno:

- `PORT`: Puerto para el servidor HTTP (por defecto: `:8080`)
- `CACHE_TTL`: Tiempo de vida del caché en minutos (por defecto: `5`)

## Uso

### Estructura del Proyecto

```
.
├── cmd/                  # Punto de entrada de la aplicación
├── internal/
│   ├── core/             # Entidades y casos de uso
│   ├── infrastructure/   # Implementaciones concretas (repositorios)
│   └── interfaces/       # Controladores HTTP y rutas
├── get_dolar.py          # Script Python para obtener tasas
└── README.md             # Este archivo
```

## Endpoints

### Obtener Tasa de Cambio

Obtiene la tasa de cambio actual para una moneda específica.

- **Método**: `GET`
- **URL**: `/api/v1/rates`
- **Parámetros de consulta**:
  - `currency` (opcional): Código de moneda (`usd` o `eur`). Por defecto: `usd`
- **Ejemplo de respuesta exitosa (200 OK)**:
  ```json
  {
    "currency": "USD",
    "price": 138.13,
    "last_updated": "2023-08-20T12:00:00Z",
    "symbol": "Bs.",
    "percent_change": 0.5
  }
  ```
- **Códigos de error**:
  - `400`: Moneda no soportada
  - `500`: Error al obtener la tasa de cambio

### Obtener Todas las Tasas

Obtiene las tasas de cambio para todas las monedas soportadas.

- **Método**: `GET`
- **URL**: `/api/v1/rates/all`
- **Ejemplo de respuesta exitosa (200 OK)**:
  ```json
  {
    "USD": {
      "currency": "USD",
      "price": 138.13,
      "last_updated": "2023-08-20T12:00:00Z",
      "symbol": "Bs.",
      "percent_change": 0.5
    },
    "EUR": {
      "currency": "EUR",
      "price": 152.45,
      "last_updated": "2023-08-20T12:00:00Z",
      "symbol": "Bs.",
      "percent_change": 0.3
    }
  }
  ```

### Convertir Moneda

Realiza la conversión entre dos monedas.

- **Método**: `POST`
- **URL**: `/api/v1/convert`
- **Cuerpo de la solicitud (JSON)**:
  ```json
  {
    "amount": 100,
    "from": "USD",
    "to": "VES"
  }
  ```
  - `amount` (number, requerido): Cantidad a convertir
  - `from` (string, requerido): Moneda de origen (USD, EUR, VES)
  - `to` (string, requerido): Moneda de destino (USD, EUR, VES)

- **Ejemplo de respuesta exitosa (200 OK)**:
  ```json
  {
    "from_amount": 100,
    "to_amount": 13813,
    "rate": 138.13,
    "from_symbol": "USD",
    "to_symbol": "VES"
  }
  ```

- **Códigos de error**:
  - `400`: Parámetros inválidos o monedas no soportadas
  - `500`: Error al realizar la conversión

### Health Check

Verifica que el servicio esté en funcionamiento.

- **Método**: `GET`
- **URL**: `/api/v1/health`
- **Respuesta exitosa (200 OK)**:
  ```
  OK
  ```

## Formatos de Moneda

La API acepta múltiples formatos para cada moneda:

- **USD**: "USD", "$", "US$", "US"
- **EUR**: "EUR", "€", "EURO", "EUROS"
- **VES**: "VES", "BS", "BS.", "Bs", "Bs.", "bs", "bs.", "BOLIVAR", "BOLIVARES"

## Manejo de Errores

La API devuelve respuestas de error en el siguiente formato:

```json
{
  "error": "Mensaje de error descriptivo",
  "code": "CÓDIGO_DEL_ERROR"
}
```

## Caché

Las tasas de cambio se almacenan en caché para mejorar el rendimiento. El tiempo de vida del caché es de 5 minutos por defecto, pero puede configurarse mediante la variable de entorno `CACHE_TTL`.

## Pruebas

Para ejecutar las pruebas unitarias:

```bash
go test ./...
```

## Despliegue

La aplicación puede desplegarse en cualquier plataforma que soporte aplicaciones Go. Se recomienda usar un sistema de gestión de procesos como PM2 o systemd.

## Contribución

1. Haz un fork del repositorio
2. Crea una rama para tu característica (`git checkout -b feature/nueva-caracteristica`)
3. Haz commit de tus cambios (`git commit -am 'Añadir nueva característica'`)
4. Haz push a la rama (`git push origin feature/nueva-caracteristica`)
5. Crea un nuevo Pull Request

## Licencia

Este proyecto está licenciado bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para más detalles.
